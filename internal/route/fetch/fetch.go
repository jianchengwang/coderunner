package fetch

import (
	"fmt"
	"github.com/gin-gonic/gin"
	r "github.com/jianchengwang/coderunner/internal/response"
	"github.com/jianchengwang/coderunner/internal/task"
	"github.com/jianchengwang/coderunner/internal/utils"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"
)

var hostGitPath = path.Join(os.Getenv("APP_CONTAINER_PATH"), "git")
func FetchGitRep(c *gin.Context) (int, interface{}) {
	gitRep := c.Query("gitRep")
	lang := c.Query("lang")

	fileExt := ""
	for _, r := range task.LangRunners {
		if r.Name == lang {
			fileExt = r.Ext
			break
		}
	}

	u, err := url.Parse(gitRep)
	if err != nil {
		panic(err)
	}
	host, _, _ := net.SplitHostPort(u.Host)
	gitPath := u.Path
	fmt.Println(host)
	fmt.Println(gitPath)

	user := strings.Split(gitPath, "/")[1]
	rep := strings.Split(gitPath, "/")[2]
	branch := strings.Split(gitPath, "/")[4]
	dir := gitPath[strings.Index(gitPath, branch) + len(branch) + 1:]
	gitUrl := "https://github.com/" + user + "/" + rep + ".git"
	fmt.Println(user, rep, branch, dir, gitUrl)

	var workDir = path.Join(hostGitPath, rep)
	if !utils.Exists(workDir) {
		fmt.Println(workDir, "not existed")
		err = os.MkdirAll(workDir, 0755)
		if err != nil {
			return r.MakeErrJSON(500, err.Error())
		}
		runGitCommand(workDir,"git", "init")
		runGitCommand(workDir,"git", "remote", "add", "origin", gitUrl)
		runGitCommand(workDir,"git", "config", "core.sparsecheckout", "true")
	} else {
		fmt.Println(workDir, "existed")
	}
	runGitCommand(workDir,"/bin/sh", "-c", "echo \"" + dir + "/*\">> .git/info/sparse-checkout")
	runGitCommand(workDir,"git", "pull", "origin", branch)
	runGitCommand(workDir,"git", "checkout")
	targetDir := path.Join(workDir, dir)
	fmt.Println(targetDir)
	fileMap := map[string]interface{}{}
	buildFileMap(targetDir, fileMap, fileExt)
	return r.MakeSuccessJSON(gin.H{
		"fileMap": fileMap,
	})
}

func buildFileMap(targetDir string, fileMap  map[string]interface{}, fileExt string) {
	files, _ := ioutil.ReadDir(targetDir)
	for _, f := range files {
		fmt.Println(f.Name())
		if f.IsDir() {
			buildFileMap(path.Join(targetDir, f.Name()), fileMap, fileExt)
		} else {
			if len(fileExt)>0 && path.Ext(f.Name()) == fileExt {
				data, err := ioutil.ReadFile(path.Join(targetDir, f.Name()))
				if err != nil {
					panic(err)
				}
				fileMap[f.Name()] = string(data)
			}
		}
	}
}


func runGitCommand(workDir string, name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	cmd.Dir = workDir
	//cmd.Stderr = os.Stderr
	msg, err := cmd.CombinedOutput()
	cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(msg))
	return string(msg), err
}