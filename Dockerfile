FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone

# install git - apt-get replace with apk
RUN apk update && \
    apk upgrade && \
    apk add --no-cache bash git openssh

RUN mkdir /etc/coderunner
WORKDIR /etc/coderunner

ADD coderunner /etc/coderunner

RUN chmod 655 /etc/coderunner/coderunner

ENTRYPOINT ["/etc/coderunner/coderunner"]
EXPOSE 8080