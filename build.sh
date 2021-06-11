#/bin/bash
# Linux
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=1
go build -o rtsp2rtmp main.go

# Windows
# export GOOS=windows
# export GOARCH=amd64
# export CGO_ENABLED=1
# go build -o rtsp2rtmp.exe main.go

#rm -rf ./output/releases

#mkdir -p ./output/releases/output/live
#mkdir -p ./output/releases/output/log
#mkdir -p ./output/releases/conf

# cp -r ./static ./output/releases/static/
# cp -r ./db ./output/releases/db/
# cp -r ./conf/conf.yml ./output/releases/conf
# cp -r ./rtsp2rtmp ./output/releases
# cp -r ./rtsp2rtmp.exe ./output/releases
# cp -r ./start.vbs ./output/releases
