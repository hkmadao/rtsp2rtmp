module github.com/hkmadao/rtsp2rtmp

go 1.16

require (
	github.com/beego/beego/v2 v2.0.1
	// github.com/deepch/vdk v0.0.0-20241120073805-439b6309323c //gitUrl v0.0.0-timestamp-commitId
	github.com/deepch/vdk v0.0.27
	github.com/gin-gonic/gin v1.7.2
	github.com/google/uuid v1.3.0
	github.com/lib/pq v1.10.9
)

replace github.com/deepch/vdk => github.com/hkmadao/vdk v0.0.0-20241127071358-df60b9bc5ae8
