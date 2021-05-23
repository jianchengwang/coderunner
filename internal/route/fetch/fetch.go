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
		err = os.MkdirAll(workDir, 0755)
		if err != nil {
			panic(err)
		}
		RunGitCommand(workDir,"git", "init")
		RunGitCommand(workDir,"git", "remote", "add", "origin", gitUrl)
		RunGitCommand(workDir,"git", "config", "core.sparsecheckout", "true")
		RunGitCommand(workDir,"echo", dir, ">>", ".git/info/sparse-checkout")
		RunGitCommand(workDir,"git", "pull", "origin", branch)
		RunGitCommand(workDir,"git", "checkout")
	}
	targetDir := path.Join(workDir, dir)
	fmt.Println(targetDir)

	fileMap := map[string]interface{}{}
	WorkDir(targetDir, fileMap, fileExt)

	return r.MakeSuccessJSON(gin.H{
		"fileMap": fileMap,
	})
}

func WorkDir(targetDir string, fileMap  map[string]interface{}, fileExt string) {
	files, _ := ioutil.ReadDir(targetDir)
	for _, f := range files {
		fmt.Println(f.Name())
		if f.IsDir() {
			WorkDir(path.Join(targetDir, f.Name()), fileMap, fileExt)
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


func RunGitCommand(workDir string, name string, arg ...string) (string, error) {

	cmd := exec.Command(name, arg...)
	cmd.Dir = workDir
	//cmd.Stderr = os.Stderr
	msg, err := cmd.CombinedOutput()
	cmd.Run()

	// 报错时 exit status 1
	fmt.Println(string((msg)))
	return string(msg), err
}