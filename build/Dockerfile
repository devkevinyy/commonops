FROM registry.cn-hangzhou.aliyuncs.com/kevin_yang/ops_java_db_base:1.0
MAINTAINER ChujieYang "yangchujie1@163.com"

WORKDIR /go/src/github.com/chujieyang/ops/

COPY ./ref/* ./
COPY ops.go ./
ADD ./web.tar ./
RUN mkdir ./jars && mkdir ./conf && chmod +x ./mysql.sh && chmod +x ./entrypoint.sh && touch ops.log
COPY ./ref/config.ini ./conf/
COPY ./ref/sqlparser.jar ./jars/

EXPOSE 80
EXPOSE 9999

ENTRYPOINT ["./entrypoint.sh"]