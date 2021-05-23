FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone

RUN mkdir /etc/coderunner
WORKDIR /etc/coderunner

ADD main /etc/coderunner

RUN chmod 655 /etc/coderunner/main

ENTRYPOINT ["/etc/coderunner/main"]
EXPOSE 8080