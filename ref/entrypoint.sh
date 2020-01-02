#!/bin/sh

./mysql.sh
mkdir -p /run/nginx/
/usr/sbin/nginx -c /go/src/github.com/chujieyang/ops/nginx.conf
nohup ./ops.go 2>&1 > ops.log &
tail -f ops.log