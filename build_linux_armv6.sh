#/bin/bash
ver=$1
if [ -n "${ver}" ] 
then
    echo package version ${ver}
else
    echo no version param
    exit 1
fi
# Linux
# export GOOS=linux
# export GOARCH=amd64
export CGO_ENABLED=1
go build -o rtsp2rtmp_${ver}_linux_armv6 main.go

# Windows
# export GOOS=windows
# export GOARCH=amd64
# export CGO_ENABLED=1
# go build -o rtsp2rtmp.exe main.go

#package linux_amd64
rm -rf ./output/rtsp2rtmp_${ver}_linux_amd64

mkdir -p ./output/rtsp2rtmp_${ver}_linux_amd64/output/live
mkdir -p ./output/rtsp2rtmp_${ver}_linux_amd64/output/log
mkdir -p ./output/rtsp2rtmp_${ver}_linux_amd64/conf

cp -r ./static ./output/rtsp2rtmp_${ver}_linux_amd64/static/
cp -r ./db ./output/rtsp2rtmp_${ver}_linux_amd64/db/
cp -r ./conf/conf.yml ./output/rtsp2rtmp_${ver}_linux_amd64/conf
cp -r ./rtsp2rtmp_${ver}_linux_amd64 ./output/rtsp2rtmp_${ver}_linux_amd64/rtsp2rtmp

#package linux_armv6
rm -rf ./output/rtsp2rtmp_${ver}_linux_armv6

mkdir -p ./output/rtsp2rtmp_${ver}_linux_armv6/output/live
mkdir -p ./output/rtsp2rtmp_${ver}_linux_armv6/output/log
mkdir -p ./output/rtsp2rtmp_${ver}_linux_armv6/conf

cp -r ./static ./output/rtsp2rtmp_${ver}_linux_armv6/static/
cp -r ./db ./output/rtsp2rtmp_${ver}_linux_armv6/db/
cp -r ./conf/conf.yml ./output/rtsp2rtmp_${ver}_linux_armv6/conf
cp -r ./rtsp2rtmp_${ver}_linux_armv6 ./output/rtsp2rtmp_${ver}_linux_armv6/rtsp2rtmp

#package window_amd64
rm -rf ./output/rtsp2rtmp_${ver}_window_amd64

mkdir -p ./output/rtsp2rtmp_${ver}_window_amd64/output/live
mkdir -p ./output/rtsp2rtmp_${ver}_window_amd64/output/log
mkdir -p ./output/rtsp2rtmp_${ver}_window_amd64/conf

cp -r ./static ./output/rtsp2rtmp_${ver}_window_amd64/static/
cp -r ./db ./output/rtsp2rtmp_${ver}_window_amd64/db/
cp -r ./conf/conf.yml ./output/rtsp2rtmp_${ver}_window_amd64/conf
cp -r ./rtsp2rtmp_${ver}_window_amd64.exe ./output/rtsp2rtmp_${ver}_window_amd64/rtsp2rtmp.exe
cp -r ./start.vbs ./output/rtsp2rtmp_${ver}_window_amd64/start.vbs

cd ./output/
tar -zcvf ./rtsp2rtmp_${ver}_linux_amd64.tar.gz ./rtsp2rtmp_${ver}_linux_amd64/
tar -zcvf ./rtsp2rtmp_${ver}_linux_armv6.tar.gz ./rtsp2rtmp_${ver}_linux_armv6/
tar -zcvf ./rtsp2rtmp_${ver}_window_amd64.tar.gz ./rtsp2rtmp_${ver}_window_amd64/

rm -rf ./rtsp2rtmp_${ver}_linux_amd64/
rm -rf ./rtsp2rtmp_${ver}_linux_armv6/
rm -rf ./rtsp2rtmp_${ver}_window_amd64/
