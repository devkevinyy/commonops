### 安装指南


#### 1、docker 快速访问

* 执行命令

```
   * 生成镜像
   项目目录下执行：./build.sh
   
   * 启动容器
   docker run -d -p 8880:80 -p 19999:9999 registry.cn-hangzhou.aliyuncs.com/kevin_yang/ops_system:0.3
   
   * web 访问
   http://localhost:8880
   
   * 账号信息
   账号：admin@ops.com 密码：123456
```

#### 2、普通安装

* 环境准备

```
    * Node & yarn 环境安装
    参考文档： https://nodejs.org/zh-cn/ 
    
    * Go 环境安装
    参考文档： https://www.runoob.com/go/go-environment.html
    
    * Mysql 安装
    参考文档： https://www.runoob.com/mysql/mysql-install.html

```

* 获取代码

```
    # 获取代码
    https://github.com/chujieyang/commonops.git
```

* 依赖包安装

```text
    # 进入前端项目目录 web。  （备注：node 安装 yarn）
    yarn install 
    
    # 后端的依赖包已经放置在后端代码 vendor 目录中
```

* 数据库初始化

```text
    # 使用 [project_path]/commonops/ops/init_sql 下的 common_ops.sql 初始化数据库。
    # mysql client cmd 
    
    CREATE DATABASE common_ops CHARSET=UTF8;
    USE common_ops;
    SOURCE common_ops.sql;
```

* 配置 & 启动

```text
    # 配置数据库信息
    [project_path]/commonops/ops/conf/mysql.go
    
    # 配置后端服务端口信息
    [project_path]/commonops/ops/conf/app.go
    
    # 配置前端 api 访问地址
    [project_path]/commonops/web/src/config.js
    
    # 启动后端服务
    go run [project_path]/commonops/ops/main.go
    
    # 启动前端站点
    cd [project_path]/commonops/web/ & yarn start
    
    # 使用浏览器访问：http://localhost:3000
    
```
