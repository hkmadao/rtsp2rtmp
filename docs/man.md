#### rtsp2rtmp

![](./images/rtsp2rtmpad.png)

##### 项目说明：

1. 用户配置摄像头信息，包括（摄像头编号：code、摄像头rtsp地址：rtspURL、推送的rtmp地址：rtmpURL等）
2. 系统连接到摄像头，获取音视频数据
3. 系统解析摄像头数据，保存为flv文件
4. 系统推送搭到配置好的rtmp服务器
5. 用户请求观看视频，系统返回视频数据给用户播放

##### 解析说明：

1. 音视频编解码使用的是[开源项目](https://github.com/deepch/vdk.git)的功能
2. 服务器连接到摄像头，服务器获取到音视频数据，解析为av.packet，分发给FileFlvManager和HttpFlvManager处理
3. FileFlvManager将数据封装为flv文件的数据格式，写入文件
4. RtmpFlvManager将数据封装为rtmp流的数据格式，发送到rtmp服务器
5. 用户通过http方式和服务器连接请求视频数据，HttpFlvManager将av.packet封装为httpflv格式数据返回

##### 配置说明：

```
server:
    user:
        name: admin #网页登录用户名
        password: admin #网页登录密码
    rtmp:
    httpflv:
        port: 8080
        static:
            path: ./resources/static #页面所在文件夹
    fileflv:
        path: ./resources/output/live #录像所在文件夹
    log:
        path: ./resources/output/log #日志所在文件夹  
        level: 6 #1-7 7输出的信息最多 
    database:
        driver-type: 4 #数据库类型
        driver: postgres #数据库驱动
        url: user=postgres password=123456 dbname=rtsp2rtmp host=localhost port=5432 sslmode=disable TimeZone=UTC #数据库url
        show-sql: false     #是否打印sql                
```

##### 开发说明：

程序分为服务器和页面，服务端采用golang开发，前端采用react+materia-ui，完成后编译页面文件放入服务器的resources/static文件夹,或者修改配置文件页面所在文件夹的路径

###### 服务器开发说明：

1. 安装golang
2. 获取[服务器源码](https://github.com/hkmadao/rtsp2rtmp.git)
3. 安装postgresql数据库，根据配置文件"resources/conf/conf-prod.yml"创建数据库
4. 根据"docs/init/rtsp2rtmp-postgresql.sql"文件创建表    
5. 进入项目目录
6. go build开发

###### 页面开发说明：

1. 安装node
2. 下载[页面源码](https://github.com/hkmadao/rtsp2rtmp-web.git)
3. 进入项目目录
4. npm install
5. npm run start