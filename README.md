#### rtsp2rtmp

![](./docs/images/rtsp2rtmpad.png)

##### 项目功能：

1. rtsp转httpflv播放
2. rtsp转rtmp推送
3. rtsp视频录像，录像文件为flv格式

##### 运行说明：

1. 下载[程序文件](https://github.com/hkmadao/rtsp2rtmp/releases)，解压   
2. 执行程序文件：window下执行rtsp2rtmp.exe，linux下执行rtsp2rtmp   
3. 浏览器访问程序服务地址：http://[server_ip]:8080/   
4. 在网页配置摄像头的rtsp地址、要推送到的rtmp服务器地址等信息  
5. 配置好后重启服务器  
6. 再次进入网页观看视频      

> 注意：
>
> ​	程序目前支持h264视频编码、aac音频编码，若不能正常播放，关掉摄像头推送的音频再尝试

##### 目录结构：

```
--rtsp2rtmp #linux执行文件
--rtsp2rtmp.exe #window执行文件
  --statics #程序的网页文件夹
  --conf #配置文件文件夹
    --conf.yml #配置文件
  --db #sqlite3 #数据库文件夹
    --rtsp2rtmp.db #sqlite3数据库文件（存放摄像头的url、推送的rtmp服务器地址等信息）
  --output #程序输出文件夹
    --live #保存摄像头录像的文件夹，录像格式为flv
    --log #程序输出的日志文件夹
```

##### 配置说明：

```
server:
    httpflv:
        port: 8080 #程序的http端口
    fileflv:
        save: true #是否保存录像文件
        path: ./output/live #录像文件夹
    log:
        path: ./output/log #日志文件夹
        
```

##### 开发说明：

程序分为服务器和页面，服务端采用golang开发，前端采用react+materia-ui，完成后编译页面文件放入服务器的statics文件夹

###### 服务器开发说明：

1. 安装golang，gc++编译器(sqlite3模块的需要用到，window下可选择安装MinGW)
2. 获取[服务器源码](https://github.com/hkmadao/rtsp2rtmp.git)
3. 进入项目目录
4. go build开发

###### 页面开发说明：

1. 安装node
2. 下载[页面源码](https://github.com/hkmadao/rtsp2rtmp-web.git)
3. 进入项目目录
4. npm install
5. npm run start