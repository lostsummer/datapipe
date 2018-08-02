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
    <httpserver enable="true">
        <importer enable="true" name="PageClick" servertype="redis" serverurl="192.168.8.175:6379" toqueue="EMoney.Tongji:PageClickStringQueue"/>
        <importer enable="true" name="PageView" servertype="redis" serverurl="192.168.8.175:6379" toqueue="EMoney.Tongji:PageViewStringQueue"/>
        <importer enable="true" name="WebData" servertype="redis" serverurl="192.168.8.175:6379" toqueue="EMoney.Tongji:Data:WebDataStringQueue"/>
        <importer enable="true" name="AppData" servertype="redis" serverurl="192.168.8.175:6379" toqueue="EMoney.Tongji:Data:AppDataStringQueue"/>
        <importer enable="true" name="PayLog" servertype="redis" serverurl="192.168.8.175:6379" toqueue="EMoney.Tongji:Data:PayLogStringQueue"/>
        <importer enable="true" name="UserLog" servertype="redis" serverurl="192.168.8.175:6379" toqueue="EMoney.Tongji:Data:UserLogStringQueue"/>
        <importer enable="true" name="Soft" servertype="redis" serverurl="192.168.8.175:6379" toqueue="EMoney.Tongji:SoftLogStringQueue"/>
        <importer enable="true" name="SoftActionLog" servertype="redis" serverurl="192.168.8.175:6379" toqueue="EMoney.Tongji:SoftActionLogJsonQueue"/>
        <importer enable="true" name="FrontEndLog" servertype="redis" serverurl="192.168.8.175:6379" toqueue="EMoney.Tongji:FrontEndLogJsonQueue"/>
        <importer enable="true" name="LiveDruation" servertype="redis" serverurl="192.168.8.175:6379" toqueue="EMoney.Tongji:LiveDurationStringQueue"/>
    </httpserver>
```

httpserver enable 属性是指定是否开启httpserver的开关
httpserver 下属 importer 对应每个 http route handler，配置了对应的写入redis地址和对列key, 也有单独的加载开关

httpserver自身运行的配置文件在单独的 dotweb.conf 中

```xml
<?xml version="1.0" encoding="UTF-8"?>
<config>
<app logpath="logs/" enabledlog="true" runmode="development"/>
<server isrun="true" port="8080" />
<appset>
    <set key="set1" value="1" />
</appset>
<middlewares>
</middlewares>
<routers>
</routers>
<groups>
</groups>
</config>　

```
## 通用计数服务

2018-08 根据需求新增，详见：[wiki](http://git.emoney.cn/techplat/datapipe/wikis/common-counter)

配置, httpserver 下:

```xml
   <httpserver enable="true">
        <accumulator enable="true" name="PVCounter" servertype="redis" serverurl="172.28.1.118:6379" tocounter="EMoney.DataPipe:Counter"/>
        <accumulator enable="true" name="UVCounter" servertype="redis" serverurl="172.28.1.118:6379" tocounter="EMoney.DataPipe:Counter" toset="EMoney.DataPipe:UserSet"/>
    </httpserver>
```

"tcounter" 为redis中计数key前缀
"toset" 为统计去重用户数（用于UV统计）使用的set的key前缀
