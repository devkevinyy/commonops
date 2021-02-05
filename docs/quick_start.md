### 安装指南

#### 1、docker 方式

> 1.1、启动容器：
>
> docker run -d -p 8880:80 -p 9999:9999 registry.cn-hangzhou.aliyuncs.com/kevin_yang/ops_system:0.6
>
> 1.2、浏览器访问：
> http://[you ip]:8880
>
> 1.3、登陆信息：
> 账号：admin@ops.com 密码：123456

#### 2、宿主机方式

> 1、构建环境准备
>
> 请预先准备并安装 Node、yarn、golang-1.15.2 等组件
>
> 2、运行环境准备
>
> 请预先准备并安装 mysql、nginx 等组件
>
> 3、获取代码
>
> git clone https://github.com/chujieyang/commonops.git
>
> 4、后端源码构建
>
> > 4.1、配置数据库信息
> >
> >     commonops/ops/conf/mysql.go
> >
> > 4.2、配置后端服务端口信息
> >
> >     commonops/ops/conf/app.go
> >
> > 4.3、构建
> >
> >     make clean buildSvr
> >
> >     构建制品输出到 commonops/build/ops.go (二进制文件)
>
> 5、前端源码构建
>
> > 5.1、配置前端访问后端 api 地址
> >
> >     commonops/web/src/config.js
> >
> > 5.2、构建
> >
> >     make clean buildWeb
> >
> >     构建制品输出到 commonops/build/web.tar (前端文件压缩包)
>
> 6、数据库初始化
>
> 使用 commonops/build/ref/ 下的 common_ops.sql 初始化数据库。
>
> ```
> CREATE DATABASE common_ops CHARSET=UTF8;
> USE common_ops;
> SOURCE common_ops.sql;
> ```
>
> 7、Nginx 配置
>
> 使用 common_ops/build/ref/nginx.conf 进行前后端服务的配置
>
> 8、服务启动
>
> 8.1、后端服务启动:
>
> nohup ./ops.go 2>&1 > ops.log &
>
> 8.2、nginx 启动:
>
> /usr/sbin/nginx -c /go/src/github.com/chujieyang/ops/nginx.conf
