### 运维 Ops 平台

目前该项目为个人开发维护的一个运维管理平台。目前功能还在持续迭代开发完善中。
该平台采用 React + Golang 技术体系构建。
你可以通过以下方式来尝试体验，相关的建议和意见可以通过 Issue 来进行反馈；如果你想参与到项目中，可通过邮件联系：yangchujie@sina.com

#### 功能介绍

-   阿里云资源的同步集成，可查看资源核心指标监控图
-   通过 Web SSH 的方式快速便捷管理服务器
-   通过 DMS 在保障安全的前提下方便用户进行数据库访问
-   Kubernetes 集群管理，可查看 Node、Pod 的核心指标监控图
-   通过 Web SSH 的方式安全便捷访问 Kubernetes 容器
-   工单模块帮助团队协作
-   基于 RBAC 的权限系统可细分控制团队成员角色的权限

#### 体验方式

为方便用户能快速体验该项目，目前提供 docker image 直接访问的方式：

```
启动方式

docker run -d -p 8880:80 -p 9999:9999 registry.cn-hangzhou.aliyuncs.com/kevin_yang/ops_system:0.6

Host 设置
设置访问侧 host 文件: commonops.com  --->  you server ip

web 访问

http://commonops.com:8880

初次访问账号

账号：admin@ops.com 密码：123456
```

#### 项目详细文档

[项目文档](https://chujieyang.github.io/commonops/)
