#!/bin/sh

echo "开始生成 ops 集成镜像 ..."

rm -f ops.go

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./ref/ops.go ./ops/

rm -f web.tar
tar -cvf web.tar ./web/build/*

docker build -t yangchujie/ops_system:v1 .

echo "镜像生成完成！"
