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

## datapipe 线上部署情况

目前线上 datapipe 在线上部署多台服务器， httpserver 和 taskservice 按不同进程实例不同配置部署，主要因为：

1. 升级时功能互相不影响
2. 减少redis连接池争用的情况

### httpserver (上报消息息接入，队列生成者)

部署服务器：

1. 172.28.1.45
2. 172.28.1.46
3. 172.31.37.8

部署目录：

`/emoney/datapipe`

配置文件：

```
/emoney/datapipe/conf/production/app.conf
/emoney/datapipe/conf/production/dotweb.conf
```

- app.conf

  ```xml
  <?xml version="1.0" encoding="UTF-8"?>
  <config>
      <httpserver enable="true">
          <importer enable="true" name="PageClick" servertype="redis" serverurl="172.28.1.118:6379" toqueue="EMoney.DataPipe:PageClick"/>
          <importer enable="true" name="PageView" servertype="redis" serverurl="172.28.1.118:6379" toqueue="EMoney.DataPipe:PageView"/>
          <importer enable="true" name="ADView" servertype="redis" serverurl="172.28.1.118:6379" toqueue="EMoney.Tongji:Data:ADViewStringQueue"/>
          <importer enable="true" name="ADClick" servertype="redis" serverurl="172.28.1.118:6379" toqueue="EMoney.Tongji:Data:ADClickStringQueue"/>
          <importer enable="true" name="WebData" servertype="redis" serverurl="172.28.1.118:6379" toqueue="EMoney.Tongji:Data:WebDataStringQueue"/>
          <importer enable="true" name="AppData" servertype="redis" serverurl="172.28.1.118:6379" toqueue="EMoney.Tongji:Data:AppDataStringQueue"/>
          <importer enable="true" name="PayLog" servertype="redis" serverurl="172.28.1.118:6379" toqueue="EMoney.Tongji:Data:PayLogStringQueue"/>
          <importer enable="true" name="UserLog" servertype="redis" serverurl="172.28.1.118:6379" toqueue="EMoney.Tongji:Data:UserLogStringQueue"/>
          <importer enable="true" name="Soft" servertype="redis" serverurl="172.28.1.118:6379" toqueue="EMoney.Tongji:SoftLogStringQueue"/>
          <importer enable="true" name="SoftActionLog" servertype="redis" serverurl="172.28.1.118:6379" toqueue="EMoney.Tongji:SoftActionLogJsonQueue"/>
          <importer enable="true" name="ActLog" servertype="redis" serverurl="172.28.1.118:6379" toqueue="EMoney.Tongji:SoftActionLogJsonQueue"/>
          <importer enable="true" name="FrontEndLog" servertype="redis" serverurl="172.28.1.118:6379" toqueue="EMoney.Tongji:FrontEndLogJsonQueue"/>
          <importer enable="true" name="LiveDuration" servertype="redis" serverurl="172.28.1.118:6379" toqueue="EMoney.Tongji:LiveDurationStringQueue"/>
          <accumulator enable="true" name="PVCounter" servertype="redis" serverurl="172.28.1.118:6379" tocounter="EMoney.DataPipe:Counter"/>
          <accumulator enable="true" name="UVCounter" servertype="redis" serverurl="172.28.1.118:6379" tocounter="EMoney.DataPipe:Counter" toset="EMoney.DataPipe:UserSet"/>
      </httpserver>
  <redis keycommonpre="EMoney.TongJiService"></redis>
  <mongodb serverurl="mongodb://webloguser:weblogem123@172.28.1.39:27017/emoney_tongji" dbname="emoney_tongji"></mongodb>
  <kafka serverurl="172.31.37.35:9092,172.31.37.36:9092,172.31.37.37:9092"></kafka>
  <log filepath="/emoney/datapipe/logs"></log>
  <tasks>
  </tasks>
  </config>
  ```

  

启动脚本：

`/emoney/datapipe/rundatapipe.sh`

- rundatapipe.sh

  ```bash
  #!/bin/bash
  
  export RunEnv=production && cd /emoney/datapipe && nohup ./datapipe > /emoney/datapipe/logs/stdout.log &
  ```

日志文件目录：

```
/emoney/datapipe/logs
/emoney/datapipe/innerlogs
```



### taskservice （消息分发，队列消费者）

部署服务器：

1. 172.28.1.45
2. 172.28.1.46

部署目录：

`/emoney/datapipe_task`

配置文件：

```
/emoney/datapipe_task/conf/production/app.conf
/emoney/datapipe_task/conf/production/dotweb.conf
```

- app.conf

  ```xml
  <?xml version="1.0" encoding="UTF-8"?>
  <config>
      <httpserver enable="false">
          <importer enable="true" name="PageClick" servertype="redis" serverurl="172.28.1.118:6379" toqueue="EMoney.DataPipe:PageClick"/>
          <importer enable="true" name="PageView" servertype="redis" serverurl="172.28.1.118:6379" toqueue="EMoney.DataPipe:PageView"/>
          <importer enable="true" name="ADView" servertype="redis" serverurl="172.28.1.118:6379" toqueue="EMoney.Tongji:Data:ADViewStringQueue"/>
          <importer enable="true" name="ADClick" servertype="redis" serverurl="172.28.1.118:6379" toqueue="EMoney.Tongji:Data:ADClickStringQueue"/>
          <importer enable="true" name="WebData" servertype="redis" serverurl="172.28.1.118:6379" toqueue="EMoney.Tongji:Data:WebDataStringQueue"/>
          <importer enable="true" name="AppData" servertype="redis" serverurl="172.28.1.118:6379" toqueue="EMoney.Tongji:Data:AppDataStringQueue"/>
          <importer enable="true" name="PayLog" servertype="redis" serverurl="172.28.1.118:6379" toqueue="EMoney.Tongji:Data:PayLogStringQueue"/>
          <importer enable="true" name="UserLog" servertype="redis" serverurl="172.28.1.118:6379" toqueue="EMoney.Tongji:Data:UserLogStringQueue"/>
          <importer enable="true" name="Soft" servertype="redis" serverurl="172.28.1.118:6379" toqueue="EMoney.Tongji:SoftLogStringQueue"/>
          <importer enable="true" name="SoftActionLog" servertype="redis" serverurl="172.28.1.118:6379" toqueue="EMoney.Tongji:SoftActionLogJsonQueue"/>
          <importer enable="true" name="ActLog" servertype="redis" serverurl="172.28.1.118:6379" toqueue="EMoney.Tongji:SoftActionLogJsonQueue"/>
          <importer enable="true" name="FrontEndLog" servertype="redis" serverurl="172.28.1.118:6379" toqueue="EMoney.Tongji:FrontEndLogJsonQueue"/>
          <importer enable="true" name="LiveDuration" servertype="redis" serverurl="172.28.1.118:6379" toqueue="EMoney.Tongji:LiveDurationStringQueue"/>
          <accumulator enable="true" name="Counter" servertype="redis" serverurl="172.28.1.118:6379" tocounter="EMoney.DataPipe:Counter"/>
      </httpserver>
  <redis keycommonpre="EMoney.TongJiService"></redis>
  <mongodb serverurl="mongodb://webloguser:weblogem123@172.28.1.39:27017/emoney_tongji" dbname="emoney_tongji"></mongodb>
  <kafka serverurl="172.31.37.35:9092,172.31.37.36:9092,172.31.37.37:9092"></kafka>
  <log filepath="/emoney/datapipe_task/logs"></log>
  <tasks>
      <task taskid="PageView" fromserver="172.28.1.118:6379" fromqueue="EMoney.DataPipe:PageView" targettype="mongodb" targetvalue="capped_datapipe_pageview"  triggerserver="172.28.1.118:6379"  triggerqueue="EMone
  y.DataPipe:PageView_kafka" counterserver="172.28.1.118:6379" counterkey="EMoney.DataPipe:PageViewCounter"/>
      <task taskid="PageView_2" fromserver="172.28.1.118:6379" fromqueue="EMoney.DataPipe:PageView" targettype="mongodb" targetvalue="capped_datapipe_pageview"  triggerserver="172.28.1.118:6379"  triggerqueue="EMo
  ney.DataPipe:PageView_kafka" counterserver="172.28.1.118:6379" counterkey="EMoney.DataPipe:PageViewCounter"/>
      <task taskid="PageView_Kafka" fromserver="172.28.1.118:6379" fromqueue="EMoney.DataPipe:PageView_kafka" targettype="kafka" targetvalue="datapipe-pageview"  triggerserver="172.28.1.118:6379"  triggerqueue="EM
  oney.Tongji:PageView_pipe"/>
      <task taskid="PageClick" fromserver="172.28.1.118:6379" fromqueue="EMoney.DataPipe:PageClick" targettype="mongodb" targetvalue="capped_datapipe_pageclick"  triggerserver="172.28.1.118:6379"  triggerqueue="EM
  oney.DataPipe:PageClick_kafka" counterserver="172.28.1.118:6379" counterkey="EMoney.DataPipe:PageClickCounter"/>
      <task taskid="PageClick_Kafka" fromserver="172.28.1.118:6379" fromqueue="EMoney.DataPipe:PageClick_kafka" targettype="kafka" targetvalue="datapipe-pageclick"  triggerserver="172.28.1.118:6379"  triggerqueue=
  "EMoney.Tongji:PageClick_pipe"/>
      <task taskid="SoftLog_ForApp_Http" fromserver="172.28.1.118:6379" fromqueue="EMoney.DataPipe:SoftLog_ForApp" targettype="http" targetvalue="http://stat.m.emoney.cn/statistics/Uninstall/PopUp" />
      <task taskid="Page_OnlineTime" fromserver="172.28.1.118:6379" fromqueue="EMoney.Tongji:Page_OnlineTime" targettype="mongodb" targetvalue="capped_datapipe_page_onlinetime"/>
      <task taskid="UserLog_MongoDB" fromserver="172.28.1.118:6379" fromqueue="EMoney.Tongji:Data:UserLogStringQueue" targettype="mongodb" targetvalue="capped_tongji_userlog"  triggerserver="172.28.1.118:6379"  tr
  iggerqueue="EMoney.Tongji:Data:UserLogStringQueue_triggerhttp"/>
      <task taskid="UserLog_Http" fromserver="172.28.1.118:6379" fromqueue="EMoney.Tongji:Data:UserLogStringQueue_triggerhttp" targettype="http" targetvalue="http://172.28.1.146/api/CMP/1.0/Customer.saveDistribute
  Customer?gate_appid=10076" />
      <task taskid="AppData_MongoDB" fromserver="172.28.1.118:6379" fromqueue="EMoney.Tongji:Data:AppDataStringQueue" targettype="mongodb" targetvalue="capped_tongji_appdata" triggerserver="172.28.1.118:6379"  tri
  ggerqueue="EMoney.Tongji:Data:AppDataStringQueue_triggerhttp"/>
      <task taskid="AppData_Http" fromserver="172.28.1.118:6379" fromqueue="EMoney.Tongji:Data:AppDataStringQueue_triggerhttp" targettype="http" targetvalue="http://172.28.1.146/api/quanshang/1.0/SendAppData?gate_
  appid=10076" />
      <task taskid="WebData_MongoDB" fromserver="172.28.1.118:6379" fromqueue="EMoney.Tongji:DataPipe:WebDataStringQueue" targettype="mongodb" targetvalue="capped_tongji_webdata" triggerserver="172.28.1.118:6379"
   triggerqueue="EMoney.Tongji:Data:WebDataStringQueue_triggerhttp"/>
      <task taskid="WebDataTrigger_Http" fromserver="172.28.1.118:6379" fromqueue="EMoney.Tongji:Data:WebDataStringQueue_triggerhttp" targettype="http" targetvalue="http://172.28.1.146/api/CMP/1.0/customer.userreg
  ister?gate_appid=10076" />
      <task taskid="ADView_MongoDB" fromserver="172.28.1.118:6379" fromqueue="EMoney.Tongji:Data:ADViewStringQueue" targettype="mongodb" targetvalue="capped_tongji_adview" triggerserver="172.28.1.118:6379" trigger
  queue="EMoney.Tongji:DataPipe:ADView_kafka"/>
      <task taskid="ADView_Kafka" fromserver="172.28.1.118:6379" fromqueue="EMoney.Tongji:DataPipe:ADView_kafka" targettype="kafka" targetvalue="datapipe-adview"/>
      <task taskid="ADClick_MongoDB" fromserver="172.28.1.118:6379" fromqueue="EMoney.Tongji:Data:ADClickStringQueue" targettype="mongodb" targetvalue="capped_tongji_adclick" triggerserver="172.28.1.118:6379" trig
  gerqueue="EMoney.Tongji:DataPipe:ADClick_kafka"/>
      <task taskid="ADClick_Kafka" fromserver="172.28.1.118:6379" fromqueue="EMoney.Tongji:DataPipe:ADClick_kafka" targettype="kafka" targetvalue="datapipe-adclick"/>
      <task taskid="FrontEndLog_MongoDB" fromserver="172.28.1.118:6379" fromqueue="EMoney.Tongji:FrontEndLogJsonQueue" targettype="mongodb" targetvalue="capped_tongji_frontendlog"/>
      <task taskid="SoftAction_MongoDB" fromserver="172.28.1.118:6379" fromqueue="EMoney.Tongji:SoftActionLogJsonQueue" targettype="mongodb" targetvalue="capped_tongji_softaction"/>
      <task taskid="SoftAction2_MongoDB" fromserver="172.28.1.118:6379" fromqueue="EMoney.Tongji:SoftActionLogJsonQueue" targettype="mongodb" targetvalue="capped_tongji_softaction"/>
      <task taskid="SoftAction3_MongoDB" fromserver="172.28.1.118:6379" fromqueue="EMoney.Tongji:SoftActionLogJsonQueue" targettype="mongodb" targetvalue="capped_tongji_softaction"/>
      <task taskid="SoftAction_free_MongoDB" fromserver="172.28.1.118:6379" fromqueue="EMoney.Tongji:SoftActionLogJsonQueue_FreeUser" targettype="mongodb" targetvalue="capped_tongji_softaction_free" triggerserver=
  "172.28.1.118:6379" triggerqueue="EMoney.Tongji:DataPipe:SoftAction_kafka"/>
      <task taskid="SoftAction2_free_MongoDB" fromserver="172.28.1.118:6379" fromqueue="EMoney.Tongji:SoftActionLogJsonQueue_FreeUser" targettype="mongodb" targetvalue="capped_tongji_softaction_free" triggerserver
  ="172.28.1.118:6379" triggerqueue="EMoney.Tongji:DataPipe:SoftAction_kafka"/>
      <task taskid="SoftAction_Kafka" fromserver="172.28.1.118:6379" fromqueue="EMoney.Tongji:DataPipe:SoftAction_kafka" targettype="kafka" targetvalue="datapipe-softaction"/>
  </tasks>
  </config>
  ```

启动脚本：

`/emoney/datapipe_task/rundatapipetask.sh`

- rundatapipetask.sh

  ```bash
  #!/bin/bash
  
  export RunEnv=production && cd /emoney/datapipe_task && nohup ./datapipe_task > /emoney/datapipe_task/logs/stdout.log &
  ```