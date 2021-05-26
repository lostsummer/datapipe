# DataPipe 



## 功能结构图

![](http://ooyi4zkat.bkt.clouddn.com/arch.png)

## TaskService 介绍

//todo: 简言之，从 Redis队列中取出数据，整理后推向下游数据库或队列



## HttpServer 介绍

"功能结构图"中， HttpServer为新增部分，功能上在TaskService上游，负责接收web前端向http server发送的各种日志，用户行为统计数据，整理后写入队列。其实和TaskService是独立的两部分，在设计上个人更倾向于将它们分开。但从部署上一个程序更方便。


###  代码

- 新增 TechPlat/datapipe/httpserver/*
- 修改 TechPlat/config/*

http server 使用 dotweb 框架

路由绑定 TechPlat/datapipe/httpserver/router.go

业务逻辑 TechPlat/datapipe/httpserver/handlers/*


