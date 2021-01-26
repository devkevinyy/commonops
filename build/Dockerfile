FROM alpine:3.9
MAINTAINER ChujieYang "yangchujie1@163.com"

RUN echo "http://mirrors.aliyun.com/alpine/v3.8/main" >> /etc/apk/repositories
RUN echo "http://mirrors.aliyun.com/alpine/v3.8/community" >> /etc/apk/repositories

RUN apk add --no-cache ca-certificates

RUN set -eux; \
	apk add --no-cache --virtual .build-deps \
		bash \
		openssl \
		nginx;

RUN apk add mariadb mariadb-client

WORKDIR /go/src/github.com/chujieyang/ops/

COPY ./ref/* ./
COPY ops.go ./
ADD ./web.tar ./
RUN mkdir ./conf && chmod +x ./mysql.sh && chmod +x ./entrypoint.sh

EXPOSE 80
EXPOSE 9999

ENTRYPOINT ["./entrypoint.sh"]