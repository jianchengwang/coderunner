package task

import (
	"bufio"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"os"
	"path"
	"strings"
	log "unknwon.dev/clog/v2"
)

var hostVolumePath = path.Join(os.Getenv("APP_CONTAINER_PATH"), "volume")
const RuntimePath = "/runtime"
type Task struct {
	ctx context.Context

	UUID   string
	RUNNER *runner

	cli         *client.Client
	ContainerID string

	SourceVolumePath string // Folder in Host: /home/<your_user>/coderunner/volume/<UUID>/
	fileName         string
}

type Output struct {
	Error bool   `json:"error"`
	Body  string `json:"body"`
}

func NewTask(language string, code []byte) (*Task, error) {
	uid := uuid.NewV4().String()

	// Set the programming language RUNNER.
	var runner *runner
	for _, r := range LangRunners {
		if r.Name == language {
			runner = &r
			break
		}
	}
	if runner == nil {
		return nil, errors.Errorf("unexpected language: %v", language)
	}

	// Create a new docker client.
	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return nil, err
	}
	cli.NegotiateAPIVersion(ctx)

	sourceVolumePath := path.Join(hostVolumePath, uid)
	err = os.MkdirAll(sourceVolumePath, 0755)
	if err != nil {
		return nil, err
	}

	fileName := "code" + runner.Ext
	filePath := path.Join(sourceVolumePath, fileName)
	err = ioutil.WriteFile(filePath, code, 0755)
	if err != nil {
		return nil, err
	}

	return &Task{
		ctx:              ctx,
		UUID:             uid,
		RUNNER:           runner,
		cli:              cli,
		SourceVolumePath: sourceVolumePath,
		fileName:         fileName,
	}, nil
}

// Run runs a task.
func (t *Task) Run() ([]*Output, error) {
	output := make([]*Output, 0, 2)
	var networkMode container.NetworkMode
	createContainerResp, err := t.cli.ContainerCreate(t.ctx,
		&container.Config{
			Image: t.RUNNER.Image,
			Tty:   true,
			WorkingDir: RuntimePath,
		},
		&container.HostConfig{
			NetworkMode: networkMode,
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: t.SourceVolumePath,
					Target: RuntimePath,
				},
			},
			Resources: container.Resources{
				NanoCPUs: t.RUNNER.MaxCPUs * 1000000000,    // 0.0001 * CPU of cpu
				Memory:   t.RUNNER.MaxMemory * 1024 * 1024, // Minimum memory limit allowed is 6MB.
			},
		}, nil, nil, t.UUID)
	if err != nil {
		return nil, err
	}
	t.ContainerID = createContainerResp.ID

	// Clean containers and folder after executed.
	//defer t.Clean()

	if err := t.cli.ContainerStart(t.ctx, t.ContainerID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	setupOutput, err := t.setupEnvironment()
	if err != nil {
		return output, err
	}
	output = append(output, setupOutput)
	if setupOutput.Error {
		return output, nil
	}

	// Execute code.
	runOutput, err := t.Exec("", "")
	if err != nil {
		return output, err
	}
	output = append(output, runOutput)

	return output, nil
}

func (t *Task) setupEnvironment() (*Output, error) {
	if len(t.RUNNER.BuildCmd) != 0 {
		return t.Exec("", "")
	}
	return &Output{}, nil
}

func (t *Task) CreateFile(code string, fileName string) (*Output, error) {
	// Create a new docker client.
	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return nil, err
	}
	cli.NegotiateAPIVersion(ctx)

	sourceVolumePath := path.Join(hostVolumePath, t.UUID)
	err = os.MkdirAll(sourceVolumePath, 0755)
	if err != nil {
		return nil, err
	}

	if len(fileName) == 0 || t.RUNNER.Name != "java" || path.Ext(fileName) != t.RUNNER.Ext  {
		fileName = "code" + t.RUNNER.Ext
	} else {
		if strings.HasPrefix(code, "package") {
			code = code[strings.Index(code, "\n"):]
		}
	}
	filePath := path.Join(sourceVolumePath, fileName)
	err = ioutil.WriteFile(filePath, []byte(code), 0755)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *Task) Exec(cmd string, fileName string) (*Output, error) {
	if len(cmd) == 0 {
		if len(t.RUNNER.BuildCmd) > 0 {
			// java file check
			if len(fileName) > 0 && t.RUNNER.Name == "java" && path.Ext(fileName) == t.RUNNER.Ext {
				cmd = strings.ReplaceAll(t.RUNNER.BuildCmd, t.RUNNER.DefaultFileName, fileName) + " && "+ strings.ReplaceAll(t.RUNNER.RunCmd, "code", strings.ReplaceAll(fileName, t.RUNNER.Ext, ""))
			} else {
				cmd = t.RUNNER.BuildCmd + " && " + t.RUNNER.RunCmd
			}
		} else {
			cmd = t.RUNNER.RunCmd
		}
	}
	idResponse, err := t.cli.ContainerExecCreate(t.ctx, t.ContainerID, types.ExecConfig{
		Env: t.RUNNER.Env,
		Cmd:[]string{"/bin/sh", "-c", cmd},
		Tty:true,
		AttachStderr:true,
		AttachStdout:true,
		AttachStdin:true,
		Detach:true,
	})
	if err != nil {
		return &Output{
			Error: true,
			Body:  err.Error(),
		}, nil
	}

	// 附加到上面创建的/bin/bash进程中
	hr, err := t.cli.ContainerExecAttach(t.ctx, idResponse.ID, types.ExecStartCheck{Detach: false, Tty: true})
	if err != nil {
		return &Output{
			Error: true,
			Body:  err.Error(),
		}, nil
	}
	// 关闭I/O
	defer hr.Close()
	// 输入
	//hr.Conn.Write([]byte("ls\r"))
	// 输出
	scanner := bufio.NewScanner(hr.Conn)
	respBody := ""
	for scanner.Scan() {
		respBody += scanner.Text() + "\n"
	}

	return &Output{
		Error: false,
		Body:  respBody,
	}, nil
}

func (t *Task) Clean() {
	if err := t.cli.ContainerStop(t.ctx, t.ContainerID, nil); err != nil {
		log.Error("Failed to stop container: %v", err)
	}

	if err := t.cli.ContainerRemove(t.ctx, t.ContainerID, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}); err != nil {
		log.Error("Failed to remove container: %v", err)
	}

	err := os.RemoveAll(t.SourceVolumePath)
	if err != nil {
		log.Error("Failed to remove volume folder: %v", err)
	}
}
