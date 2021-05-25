>Docker-based remote code runner.

## 效果图
![效果图1](https://github.com/jianchengwang/coderunner/raw/main//doc/assets/show1.png)

## 特性
- [x] terminal
- [x] fetch gitrep
- [x] upload files
- [x] support go,python,java,javascript,c...
- [ ] support markdown
- [ ] support jsbin

## 部署
参照`deploy`目录

**go build & docker build**

```shell
set GOARCH=amd64
set GOOS=linux
go build -o coderunner main.go

docker build -t coderunner:v0.0.1 .
docker tag coderunner:v0.0.1 jianchengwang/coderunner
docker login
docker push jianchengwang/coderunner
```

这里根据自己的需求，打包成基于哪种架构的二进制文件，然后生成docker镜像即可，

**docker-compose**

```yaml
version: '3'
services:
  coderunner:
    image: jianchengwang/coderunner:latest
    ports:
      - "8080:8080"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - /root/coderunner:/root/coderunner
    environment:
      APP_URL: http://localhost:8080
      APP_PASSWORD: 12345678
      APP_CONTAINER_PATH: /root/coderunner
```

这里要注意将`APP_CONTAINER_PATH`要跟docker目录进行映射，`APP_URL`就是配置允许跨域的域名地址了，

**nginx proxy**

你如果使用nginx进行代理转发的话，要配置下跨域相关，否则可能导致跨域问题，

```nginx
 proxy_set_header    Host            $host;
 proxy_set_header    X-Real-IP       $remote_addr;
 proxy_set_header    X-Forwarded-For $proxy_add_x_forwarded_for;
 location / {
    proxy_pass http://172.17.0.6:8902;
    add_header Access-Control-Allow-Origin *;
 }
```

## 参考
[Elaina](https://github.com/wuhan005/Elaina) Docker-based remote code runner.

## License
MIT，

Jut do everything you want. You happy is ok.

