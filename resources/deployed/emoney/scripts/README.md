# datapipe 程序控制脚本

## 部署

export dpctl 脚本所在目录到PATH中或者在没有和其他系统工具冲突的情况下放在 /usr/bin 目录下

## 使用方法

dpctl (stat|stop|status) (dp1|dp2|tk1|tk2)

简称指代以下程序在服务器上的唯一运行实例:

  - dp1: /emoney/dp_api2/datapipe/datapipe
  - dp2: /emoney/dp_api2/datapipe_task/datapipe_task
  - tk1: /emoney/dp_aliapi/datapipe2/datapipe2
  - tk2: /emoney/dp_aliapi/datapipe_task2/datapipe_task2
