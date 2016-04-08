#亲加直播(GotyeLive)——APP视频技术整体解决方案

##简介

##直播APP技术框架
<div align="center">
<img src="https://github.com/QPlus/GotyeLive/blob/master/pic/freamwork.jpg" width="800" alt="直播APP技术框架图" align="center"/>
</div>

目前亲加提供:

[后台服务器(golang)](https://github.com/QPlus/qplus-live-server) [IOS App(object-C)](https://github.com/QPlus/qplus-live-ios)

[Android  App(java)]亲加希望大家通过学习之后,有开源精神的程序员能够参与到这个开源项目当中，与亲加共同开发android版本。


##亲加直播交流QQ群
亲加直播—交流群_01 : [544476772](https://github.com/QPlus/GotyeLive/blob/master/pic/qpluslive-group01.png)

<img src="https://github.com/QPlus/GotyeLive/blob/master/pic/qpluslive-group01.png" width="200" alt="全民直播App视频技术"/>


##亲加直播qpluslive服务器程序

###服务器简介
该项目是亲加直播客户端[QPlusLive For IOS](https://github.com/QPlus/GotyeLive_IOS)的直播业务服务器。

该项目是使用Golang编写的直播业务服务器，可以直接运行，为了方便大家测试使用，可以使用已编译版本，[点击下载]
(https://github.com/QPlus/GotyeLive/blob/master/doc/gotyelive_server.tar.gz)


该项目完整安装全民直播APP而设计，目前是1.0版本，后期会不断更新，敬请期待。

### 使用方式
该项目需要使用Mysql,所以首先系统得安装Mysql。

创建数据库和表的SQL脚本为[gotye_open_live.sql](https://github.com/QPlus/GotyeLive/blob/master/doc/gotye_open_live.sql), 
下载解压后的.tar.gz包里面既有。

压缩包中提供了编译好的支持`Linux`的可执行文件。

压缩包中的[config.ini](https://github.com/QPlus/GotyeLive/blob/master/doc/config.ini)是服务器的配置文件，其中的内容请安格式修改, 具体说明如下:
```
[apiserver]
#服务器监听端口
http_port = 8080

[mysql]
#数据库地址
address = 192.168.1.10
#数据库名
dbname  = gotye_open_live
#数据库账号
account = app
#数据库密码
password = 123456
```
压缩包中的[run](https://github.com/QPlus/GotyeLive/blob/master/doc/run)是服务器启动脚本，使用说明如下:
run -s or start : 启动服务器程序
run -k or stop  : 关闭服务器程序
run -i or info  : 查看服务器信息
run -h or help  : 获得脚本使用说明

##API说明

[应用层协议]请查看[API_DOC](api_doc.md)


