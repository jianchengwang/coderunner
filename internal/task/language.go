package task

type runner struct {
	Name     string
	Ext      string
	Image    string
	BuildCmd string
	RunCmd   string
	MaxCPUs  int64
	MaxMemory int64
	Example  string
}

var LangRunners = []runner{
	{
		Name:     "python",
		Ext:      ".py",
		Image:    "python:3.9.1-alpine",
		BuildCmd: "",
		RunCmd:   "python3 code.py",
		MaxCPUs: 2,
		MaxMemory: 100,
		Example: "print(\"hello world.\")",
	},
	{
		Name:     "go",
		Ext:      ".go",
		Image:    "golang:1.15-alpine",
		BuildCmd: "rm -rf go.mod && go mod init code-runner && go build -v .",
		RunCmd:   "./code-runner",
		MaxCPUs: 2,
		MaxMemory: 100,
		Example: "package main\n\n" +
				 "import \"fmt\"\n\n" +
			     "func main() {\n\n" +
			     "  fmt.Println(\"hello world.\")\n\n" +
			     "}",
	},
	{
		Name:     "javascript",
		Ext:      ".js",
		Image:    "node:lts-alpine",
		BuildCmd: "",
		RunCmd:   "node code.js",
		MaxCPUs: 2,
		MaxMemory: 50,
		Example: "console.log(\"hello world.\");",
	},
}
