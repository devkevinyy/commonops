FROM alpine:3.9
MAINTAINER ChujieYang "yangchujie1@163.com"

RUN apk add tzdata
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo "Asia/Shanghai" > /etc/timezone
RUN apk --no-cache add openjdk11 --repository=http://dl-cdn.alpinelinux.org/alpine/edge/community

RUN echo "http://mirrors.aliyun.com/alpine/v3.9/main" >> /etc/apk/repositories
RUN echo "http://mirrors.aliyun.com/alpine/v3.9/community" >> /etc/apk/repositories

RUN apk add --no-cache ca-certificates

ENV JAVA_HOME=/usr/lib/jvm/java-11-openjdk
ENV PATH="$JAVA_HOME/bin:${PATH}"

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
RUN mkdir ./jars && mkdir ./conf && chmod +x ./mysql.sh && chmod +x ./entrypoint.sh
COPY ./ref/sqlparser.jar ./jars/

EXPOSE 80
EXPOSE 9999

ENTRYPOINT ["./entrypoint.sh"]