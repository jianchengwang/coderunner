package task

type runner struct {
	Name     string
	Ext      string
	Image    string
	BuildCmd string
	RunCmd   string
	DefaultFileName string
	Env 	 []string
	MaxCPUs  int64
	MaxMemory int64
	Example  string
}

var LangRunners = []runner{
	{
		Name:     "go",
		Ext:      ".go",
		Image:    "golang:1.15-alpine",
		BuildCmd: "rm -rf go.mod && go mod init code-runner && go build -v .",
		RunCmd:   "./code-runner",
		Env: []string{"GOPROXY=https://goproxy.io,direct"},
		DefaultFileName: "code.go",
		MaxCPUs: 2,
		MaxMemory: 100,
		Example: "package main\n\n" +
			"import \"fmt\"\n\n" +
			"func main() {\n\n" +
			"  fmt.Println(\"hello world.\")\n\n" +
			"}",
	},
	{
		Name:     "python",
		Ext:      ".py",
		Image:    "python:3.9.1-alpine",
		BuildCmd: "",
		RunCmd:   "python3 code.py",
		Env: []string{},
		DefaultFileName: "code.py",
		MaxCPUs: 2,
		MaxMemory: 100,
		Example: "print(\"hello world.\")",
	},
	{
		Name:     "java",
		Ext:      ".java",
		Image:    "openjdk:8u232-jdk",
		BuildCmd: "javac code.java",
		RunCmd:   "java code",
		Env: []string{},
		DefaultFileName: "code.java",
		MaxCPUs: 2,
		MaxMemory: 100,
		Example: "class code {\n  public static void main(String[] args) {\n    System.out.println(\"Hello, World!\"); \n  }\n}",
	},
	{
		Name:     "javascript",
		Ext:      ".js",
		Image:    "node:lts-alpine",
		BuildCmd: "npm config set registry https://registry.npm.taobao.org",
		RunCmd:   "node code.js",
		Env: []string{},
		DefaultFileName: "code.js",
		MaxCPUs: 2,
		MaxMemory: 50,
		Example: "console.log(\"hello world.\");",
	},
	{
		Name:     "c",
		Ext:      ".c",
		Image:    "gcc:latest",
		BuildCmd: "gcc -v code.c -o code",
		RunCmd:   "./code",
		Env: []string{},
		DefaultFileName: "code.c",
		MaxCPUs: 2,
		MaxMemory: 50,
		Example: "#include <stdio.h>\nint main()\n{\n  printf(\"Hello, World!\");\n  return 0;\n}",
	},
}
