# DataPipe 



## 功能结构图

![](http://ooyi4zkat.bkt.clouddn.com/arch.png)

## TaskService 介绍

//todo: 简言之，从 Redis队列中取出数据，整理后推向下游数据库或队列



## HttpServer 介绍

"功能结构图"中， HttpServer为新增部分，功能上在TaskService上游，负责接收web前端向http server发送的各种日志，用户行为统计数据，整理后写入队列。其实和TaskService是独立的两部分，在设计上个人更倾向于将它们分开。但从部署上一个程序更方便。

### 功能沿袭

PHP项目EMoney.Tongji： https://192.168.0.36/svn/WebPlatform/BasicPlatform/EMoney.TongJi/trunk/EMoney.TongJi.PHP

线上地址：http://api2.tongji.emoney.cn

该项目使用Yii框架，核心业务逻辑都在 EMoney.TongJi.PHP\EMoney.TongJi.FrontAPI\protected\controllers\ 目录下的*Controller.php中， 功能逻辑对应 Go http 的路由回调handler.

###  代码

- 新增 TechPlat/datapipe/httpserver/*
- 修改 TechPlat/config/*

http server 使用 dotweb 框架

路由绑定 TechPlat/datapipe/httpserver/router.go

业务逻辑 TechPlat/datapipe/httpserver/handlers/*

### 配置

参见 TechPlat/datapipe/resources/develop/app.conf

新增了 httpserver、importers

```xml
    <httpserver httpport="80"/>
    <importers>
        <importer id="PageClick" toserver="192.168.8.175:6379" toqueue="EMoney.Tongji:PageClickStringQueue"/>
        <importer id="PageView" toserver="192.168.8.175:6379" toqueue="EMoney.Tongji:PageViewStringQueue"/>
        <importer id="WebData" toserver="192.168.8.175:6379" toqueue="EMoney.Tongji:Data:WebDataStringQueue"/>
        <importer id="AppData" toserver="192.168.8.175:6379" toqueue="EMoney.Tongji:Data:AppDataStringQueue"/>
        <importer id="PayLog" toserver="192.168.8.175:6379" toqueue="EMoney.Tongji:Data:PayLogStringQueue"/>
        <importer id="UserLog" toserver="192.168.8.175:6379" toqueue="EMoney.Tongji:Data:UserLogStringQueue"/>
        <importer id="Soft" toserver="192.168.8.175:6379" toqueue="EMoney.Tongji:SoftLogStringQueue"/>
        <importer id="SoftActionLog" toserver="192.168.8.175:6379" toqueue="EMoney.Tongji:SoftActionLogJsonQueue"/>
        <importer id="FrontEndLog" toserver="192.168.8.175:6379" toqueue="EMoney.Tongji:FrontEndLogJsonQueue"/>
        <importer id="LiveDruation" toserver="192.168.8.175:6379" toqueue="EMoney.Tongji:LiveDurationStringQueue"/>
    </importers>

```

httpserver目前只配置的server端口

importers下属importer对应每个http route handler，配置了对应的写入redis地址和对列key

### 不完善

- __服务优雅重启__

  目前`kill -HUP <datapipe pid>`只能重新加载TaskServic， ~~httpserver 的重启没有在信号处理函数中实现~~(已实现了一个不优雅重启)，因为做到优雅重启有许多框架以外基础工作要做，可参考：

  https://segmentfault.com/a/1190000004445975

- __可配置项提取__

  httpserver log 开关等


