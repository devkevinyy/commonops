.DEFAULT_GOAL := display

DomainName = commonops.com:9999
AllInOneImageVersion = 1.0.4
ImageName = registry.cn-hangzhou.aliyuncs.com/kevin_yang/ops_system

WebImageVersion = 1.0.0
WebImageName = registry.cn-hangzhou.aliyuncs.com/xxxxxx/ops_web
SvrImageVersion = 1.0.0
SvrImageName = registry.cn-hangzhou.aliyuncs.com/xxxxxx/ops_svr

.PHONY = display
display:
	@echo "请选择你想要构建的阶段: all、clean、buildWeb、buildSvr、buildImage"

.PHONY = all
all: clean buildWeb buildSvr buildImage clean

.PHONY = clean
clean:
	@rm -rf ./web/build/*
	@rm -rf ./build/ops.go
	@rm -rf ./build/web.tar
	echo "编译后的文件已清理完成"

.PHONY = buildWeb
buildWeb:
	@echo "开始构建前端项目"
	@sed -i 's/DOMAIN_NAME/${DomainName}/' ./web/src/config.js
	@cd ./web && yarn build
	@tar -cvf ./build/web.tar ./web/build/*
	@sed -i 's/${DomainName}/DOMAIN_NAME/' ./web/src/config.js	
	@echo "前端项目构建完成"

.PHONY = buildSvr
buildSvr:
	@echo "开始编译后端项目"
	@cd ./ops && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../build/ops.go .
	@echo "后端项目编译完成"

.PHONY = buildAllInOneImage
buildAllInOneImage: ./build/Dockerfile
	@echo "开始根据Dockerfile生成镜像: ${ImageName}:${AllInOneImageVersion}"
	@cd build && docker build -t ${ImageName}:${AllInOneImageVersion} .
	@echo "镜像生成成功"

.PHONY = buildWebImage
buildWebImage: ./build/web.tar ./build/frontendDockerfile
	@echo "开始生产前端项目镜像: ${WebImageName}:${WebImageVersion}"
	@cd build && docker build -f ./frontendDockerfile -t ${WebImageName}:${WebImageVersion} .
	@echo "前端镜像生成成功"

.PHONY = buildSvrImage
buildSvrImage: ./build/ops.go ./build/backendDockerfile
	@echo "开始生产后端项目镜像: ${SvrImageName}:${SvrImageVersion}"
	@cd build && docker build -f ./backendDockerfile -t ${SvrImageName}:${SvrImageVersion} .
	@echo "前端镜像生成成功"
