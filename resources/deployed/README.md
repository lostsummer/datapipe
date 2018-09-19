# datapipe 线上部署配置

--update 2018-09-19

datapipe 部署结构变迁参考[wiki](http://git.emoney.cn/techplat/datapipe/wikis/datapipe-deploy-merge)

当前两套 datapipe httpserver+taskserver 都均衡部署于三台万国主机：

1. 172.31.37.8
2. 172.31.37.9
3. 172.31.37.45 

每台的部署程序、配置相同，映射于目标主机的 /emoney 目录

http://api2.tongji.emoney.cn/ 对应 /emoney/dp_api2

其中

httpserver： /emoney/dp_api2/datapipe

taskserver: /emoney/dp_api2/datapipe_task

http://aliapi.tongji.emoney.cn/ 对应 /emoney/dp_aliapi

其中 

httpserver： /emoney/dp_aliapi/datapipe2

taskserver: /emoney/dp_aliapi/datapipe_task2
