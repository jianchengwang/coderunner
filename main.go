package main

import (
	"fmt"
	"github.com/jianchengwang/coderunner/internal/route"
	"github.com/jianchengwang/coderunner/internal/task"
	"os"
	"os/signal"
	"syscall"
	"time"
	log "unknwon.dev/clog/v2"
)

func main()  {

	//创建监听退出chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Println("退出", s)
				ExitFunc(true)
			default:
				fmt.Println("other", s)
			}
		}
	}()

	_ = log.NewConsole()
	defer log.Stop()

	// Check environment config, make sure the application is safe enough.
	appURL := os.Getenv("APP_URL")
	appPassword := os.Getenv("APP_PASSWORD")
	appContainerPath := os.Getenv("APP_CONTAINER_PATH")

	if appURL == "" {
		log.Fatal("Empty APP_URL")
	}
	if appPassword == "" || len(appPassword) < 8 {
		log.Fatal("APP_PASSWORD is not strong enough")
	}
	if appContainerPath == "" {
		log.Fatal("APP_CONTAINER_PATH is empty")
	}

	r := route.New()
	err := r.Run()
	if err != nil {
		log.Fatal("Failed to start HTTP server: %v", err)
	}

	// 执行定时器定时关闭容器
	go func() {
		for t := range time.Tick(time.Minute * 5) {
			fmt.Println("Tick at", t)
			ExitFunc(false)
		}
	}()

	//fmt.Println("hello world")
	//
	//selectLang := "python"
	//code := "print(\"hello world\")"
	//
	//startAt := time.Now().UnixNano()
	//t, err := task.NewTask(selectLang, []byte(code))
	//if err != nil {
	//	panic(err)
	//}
	//output, err := t.Run()
	//if err != nil {
	//	panic(err)
	//}
	//endAt := time.Now().UnixNano()
	//
	//
	//fmt.Println("startAt", startAt)
	//fmt.Println("endAt", endAt)
	//for _,op := range output {
	//	fmt.Println(op.Body)
	//}

	//err := os.MkdirAll("/coderunner/volume", 0755)
	//if err != nil {
	//	log.Fatal("Failed to create path /elaina/volume: %v", err)
	//}
	//log.Trace("Create /elaina/volume succeed!")
	//
	//ctx := context.Background()
	//cli, err := client.NewClientWithOpts()
	//if err != nil {
	//	fmt.Println("Unable to create docker client")
	//	panic(err)
	//}
	//cli.NegotiateAPIVersion(ctx)
	//fmt.Println(cli.ImageList(context.Background(), types.ImageListOptions{}))
	//
	//os.Setenv("APP_CONTAINER_PATH", "E:\\tmp")
	//var hostVolumePath = path.Join(os.Getenv("APP_CONTAINER_PATH"), "volume")
	//uid := uuid.NewV4().String()
	//sourceVolumePath := path.Join(hostVolumePath, uid)
	//err = os.MkdirAll(sourceVolumePath, 0755)
	//if err != nil {
	//	panic(err)
	//}
	//// Make runner folder.
	////coderunnerVolumePath := path.Join("/coderunner/volume", uid)
	////err = os.MkdirAll(coderunnerVolumePath, 0755)
	////
	////// Make the `runner` folder and create code file, `code.<ext>`.
	////runnerPath := path.Join(coderunnerVolumePath, "runner")
	////err = os.MkdirAll(runnerPath, 0755)
	////if err != nil {
	////	panic(err)
	////}
	//code := "print(\"hellko world\")"
	//fileName := "code.py"
	//filePath := path.Join(sourceVolumePath, fileName)
	//err = ioutil.WriteFile(filePath, []byte(code), 0755)
	//if err != nil {
	//	panic(err)
	//}
	//
	//var networkMode container.NetworkMode
	//createContainerResp, err := cli.ContainerCreate(ctx,
	//	&container.Config{
	//		Image: "python:3.9.1-alpine",
	//		Tty:   true,
	//		WorkingDir: "/runtime",
	//	},
	//	&container.HostConfig{
	//		NetworkMode: networkMode,
	//		Mounts: []mount.Mount{
	//			{
	//				Type:   mount.TypeBind,
	//				Source: sourceVolumePath,
	//				Target: "/runtime",
	//			},
	//		},
	//		Resources: container.Resources{
	//			NanoCPUs: 2 * 1000000000,    // 0.0001 * CPU of cpu
	//			Memory:   100 * 1024 * 1024, // Minimum memory limit allowed is 6MB.
	//		},
	//	}, nil, nil, uid)
	//if err != nil {
	//	panic(err)
	//}
	//containerID := createContainerResp.ID
	//
	//// Clean containers and folder after executed.
	//defer func() {
	//	if err := cli.ContainerStop(ctx, containerID, nil); err != nil {
	//		log.Error("Failed to stop container: %v", err)
	//	}
	//
	//	if err := cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
	//		RemoveVolumes: true,
	//		Force:         true,
	//	}); err != nil {
	//		log.Error("Failed to remove container: %v", err)
	//	}
	//
	//	err := os.RemoveAll(sourceVolumePath)
	//	if err != nil {
	//		log.Error("Failed to remove volume folder: %v", err)
	//	}
	//}()
	//
	//if err := cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
	//	panic(err)
	//}
	//idResponse, err :=cli.ContainerExecCreate(ctx,containerID,types.ExecConfig{
	//	Cmd:[]string{"python3", "code.py"},
	//	Tty:true,
	//	AttachStderr:true,
	//	AttachStdout:true,
	//	AttachStdin:true,
	//	Detach:true,
	//})
	//
	//// 附加到上面创建的/bin/bash进程中
	//hr, err := cli.ContainerExecAttach(ctx, idResponse.ID, types.ExecStartCheck{Detach: false, Tty: true})
	//if err != nil {
	//	panic(err)
	//}
	//// 关闭I/O
	//defer hr.Close()
	//// 输入
	////hr.Conn.Write([]byte("ls\r"))
	//// 输出
	//scanner := bufio.NewScanner(hr.Conn)
	//for scanner.Scan() {
	//	fmt.Println(scanner.Text())
	//}

	//if err := cli.ContainerExecStart(ctx,idResponse.ID,types.ExecStartCheck{
	//
	//}); err != nil{
	//	log.Error("error when exec start ", err)
	//}
	//
	//reader, err := cli.ContainerLogs(ctx,containerID,types.ContainerLogsOptions{
	//	ShowStdout:true,
	//	ShowStderr:true,
	//})
	//
	//if err != nil{
	//	log.Error("error when containerLogs",err)
	//}
	//
	//buf := new(bytes.Buffer)
	//buf.ReadFrom(reader)
	//newStr := buf.String()
	//
	//fmt.Printf(newStr)
	//
	//go io.Copy(os.Stdout,reader)

	//<- make(chan struct{})
}

func ExitFunc(exited bool)  {
	fmt.Println("开始退出...")
	fmt.Println("执行清理...")
	hh, _ := time.ParseDuration("1h")
	for _,v := range task.SandboxMap {
		expireTime := v.LastOptTime.Add(hh)
		if len(v.T.ContainerID)>0 && (exited || expireTime.Before(time.Now())) {
			v.T.Clean()
		}
	}
	fmt.Println("结束退出...")
	os.Exit(0)
}
