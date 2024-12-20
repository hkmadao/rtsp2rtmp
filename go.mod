module github.com/hkmadao/rtsp2rtmp

go 1.16

require (
	github.com/beego/beego/v2 v2.0.1
	// github.com/deepch/vdk v0.0.0-20241120073805-439b6309323c //gitUrl v0.0.0-timestamp-commitId
	github.com/deepch/vdk v0.0.27
	github.com/gin-gonic/gin v1.7.2
	github.com/go-cmd/cmd v1.4.3
	github.com/google/uuid v1.3.0
	github.com/lib/pq v1.10.9
	github.com/matoous/go-nanoid/v2 v2.1.0 // indirect
	github.com/u2takey/go-utils v0.3.1
)

replace github.com/deepch/vdk => github.com/hkmadao/vdk v0.0.0-20241127071358-df60b9bc5ae8
