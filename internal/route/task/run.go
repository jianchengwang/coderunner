package route

import (
	r "github.com/jianchengwang/coderunner/internal/response"
	"github.com/thanhpk/randstr"
	"io/ioutil"
	"net/http"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jianchengwang/coderunner/internal/task"
)

func NewTaskHandler(c *gin.Context)  {
	uid := randstr.String(10)
	selectLang := task.LangRunners[0].Name
	code := task.LangRunners[0].Example

	t, err := task.NewTask(selectLang, []byte(code))
	if err != nil {
		panic(err)
	}
	task.SandboxMap[uid] = task.Sandbox {
		T:           t,
		LastOptTime: time.Now(),
	}

	c.Redirect(http.StatusMovedPermanently, "/r/" + uid)
}

func EditorHandler(c *gin.Context) {
	uid := c.Param("uid")
	_, prs := task.SandboxMap[uid]
	if !prs {
		c.Redirect(http.StatusMovedPermanently, "/r/new")
	}
	c.HTML(200, "sandbox.html", nil)
}

func InitTaskHandler(c *gin.Context) (int, interface{}) {
	if len(task.LangRunners) == 0 {
		c.HTML(404, "sandbox_404.html", nil)
		return -1, nil
	}
	return r.MakeSuccessJSON(gin.H{
		"sandbox": map[string]string{"name": "未命名", "lang": task.LangRunners[0].Name,  "placeholder": task.LangRunners[0].Example},
		"langRunners": task.LangRunners,
	})
}

func RunTaskHandler(c *gin.Context) (int, interface{}) {
	uid := c.Param("uid")
	selectLang := c.PostForm("lang")
	code := c.PostForm("code")

	sandbox, prs := task.SandboxMap[uid]
	var t *task.Task
	if prs {
		t = sandbox.T
	}
	hh, _ := time.ParseDuration("1h")
	expireTime := sandbox.LastOptTime.Add(hh)
	if !prs || len(t.ContainerID)==0 || t.RUNNER.Name != selectLang || expireTime.Before(time.Now()) {
		if t != nil {
			t.Clean()
		}
		t, err := task.NewTask(selectLang, []byte(code))
		if err != nil {
			return r.MakeErrJSON(500, "Failed to create task: %v", err)
		}
		task.SandboxMap[uid] = task.Sandbox {
			T:           t,
			LastOptTime: time.Now(),
		}
		startAt := time.Now().UnixNano()
		output, err := t.Run()
		if err != nil {
			if err != nil {
				return r.MakeErrJSON(500, "Failed to run task: %v", err)
			}
		}
		endAt := time.Now().UnixNano()
		return r.MakeSuccessJSON(gin.H{
			"containerId": t.ContainerID,
			"containerLang": t.RUNNER.Name,
			"result":   output,
			"startAt": startAt,
			"endAt":   endAt,
		})
	} else {
		sandbox.LastOptTime = time.Now()
		startAt := time.Now().UnixNano()
		_, err := t.CreateFile([]byte(code))
		if err != nil {
			if err != nil {
				return r.MakeErrJSON(500, "Failed to create code file: %v", err)
			}
		}
		output, err := t.Exec("")
		if err != nil {
			if err != nil {
				return r.MakeErrJSON(500, "Failed to exec cmd: %v", err)
			}
		}
		endAt := time.Now().UnixNano()
		return r.MakeSuccessJSON(gin.H{
			"containerId": t.ContainerID,
			"containerLang": t.RUNNER.Name,
			"result":   output,
			"startAt": startAt,
			"endAt":   endAt,
		})
	}
}

func ExecTerminalCmdHandler(c *gin.Context) (int, interface{}) {
	uid := c.Param("uid")
	cmd := c.PostForm("cmd")
	if len(cmd) == 0 {
		return r.MakeErrJSON(401, "Cmd cant be empty")
	}
	sandbox, prs := task.SandboxMap[uid]
	if !prs {
		return r.MakeErrJSON(404, "Cant find uid task container")
	}
	t := sandbox.T
	startAt := time.Now().UnixNano()
	output, err := t.Exec(cmd)
	if err != nil {
		if err != nil {
			return r.MakeErrJSON(500, "Failed to exec cmd: %v", err)
		}
	}
	endAt := time.Now().UnixNano()
	return r.MakeSuccessJSON(gin.H{
		"containerId": t.ContainerID,
		"containerLang": t.RUNNER.Name,
		"result":   output,
		"startAt": startAt,
		"endAt":   endAt,
	})
}

func UploadVolumeFilesHandler(c *gin.Context) (int, interface{}) {
	uid := c.Param("uid")
	sandbox, prs := task.SandboxMap[uid]
	if !prs {
		return r.MakeErrJSON(404, "Cant find uid task container")
	}
	sourceVolume := sandbox.T.SourceVolumePath
	form, err := c.MultipartForm()
	if err !=nil{
		return r.MakeErrJSON(500, "Failed to upload volume files: %v", err)
	}
	files := form.File["file"]
	fileMap := map[string]string{}
	for _, file := range files {
		if err := c.SaveUploadedFile(file, path.Join(sourceVolume, file.Filename)); err !=nil{
			return r.MakeErrJSON(500, "Failed to save volume file: %v", err)
		}
		if path.Ext(file.Filename) == sandbox.T.RUNNER.Ext {
			data, err := ioutil.ReadFile(path.Join(sourceVolume,file.Filename))
			if err != nil {
				panic(err)
			}
			fileMap[file.Filename] = string(data)
		}
	}
	if len(fileMap) > 0 {
		return r.MakeSuccessJSON(gin.H{
			"fileMap": fileMap,
		})
	}
	return r.MakeSuccessJSON(nil)
}
