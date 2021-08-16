#!/bin/sh

/usr/sbin/nginx -c /go/src/github.com/chujieyang/ops/backend_nginx.conf
nohup ./ops.go 2>&1 > ops.log &
tail -f ops.log
