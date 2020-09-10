#!/bin/sh

echo "开始生成 ops 集成镜像 ..."

rm -f ops.go

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./ref/ops.go ./ops/

rm -f web.tar
cd web && yarn build && cd -
tar -cvf web.tar ./web/build/*

docker build -t registry.cn-hangzhou.aliyuncs.com/kevin_yang/ops_system:0.3 .

echo "镜像生成完成！"
