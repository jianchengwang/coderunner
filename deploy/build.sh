# 请在主目录执行
set GOARCH=amd64
set GOOS=linux
go build -o coderunner main.go

docker build -t coderunner:v0.0.1 .
docker tag coderunner:v0.0.1 jianchengwang/coderunner
docker login
docker push jianchengwang/coderunner

