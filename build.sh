#!/bin/bash
#./build.sh 0.0.1
ver=$1
if [ -n "${ver}" ]; then
    echo package version "${ver}"
else
    echo no version param
    exit 1
fi
#打多个平台的包
platforms="linux_amd64"
rm -rf ./resources/output/demo/
for platform in $platforms; do

    export GOOS=$(echo "$platform" | gawk 'BEGIN{FS="_"} {print $1}')
    export GOARCH=$(echo "$platform" | gawk 'BEGIN{FS="_"} {print $2}')
    export CGO_ENABLED=1
    echo "${GOOS}"_"${GOARCH}"
    if [[ "${GOOS}" == "windows" ]]; then
        go build -o ./resources/output/demo/rtsp2rtmp_"${ver}"_"${GOOS}"_"${GOARCH}"/rtsp2rtmp.exe main.go
    else
        go build -o ./resources/output/demo/rtsp2rtmp_"${ver}"_"${GOOS}"_"${GOARCH}"/rtsp2rtmp main.go
    fi

    mkdir -p ./resources/output/demo/rtsp2rtmp_"${ver}"_"${GOOS}"_"${GOARCH}"/resources/output/live
    mkdir -p ./resources/output/demo/rtsp2rtmp_"${ver}"_"${GOOS}"_"${GOARCH}"/resources/output/log
    mkdir -p ./resources/output/demo/rtsp2rtmp_"${ver}"_"${GOOS}"_"${GOARCH}"/resources/conf

    cp -r ./resources/static ./resources/output/demo/rtsp2rtmp_"${ver}"_"${GOOS}"_"${GOARCH}"/resources/static/
    cp -r ./resources/conf ./resources/output/demo/rtsp2rtmp_"${ver}"_"${GOOS}"_"${GOARCH}"/resources/conf

    cd ./resources/output/demo/ || exit
    rm -rf rtsp2rtmp_"${ver}"_"${GOOS}"_"${GOARCH}".tar.gz
    tar -zcvf ./rtsp2rtmp_"${ver}"_"${GOOS}"_"${GOARCH}".tar.gz rtsp2rtmp_"${ver}"_"${GOOS}"_"${GOARCH}"/

    # rm -rf ./rtsp2rtmp_"${ver}"_"${GOOS}"_"${GOARCH}"/
    cd ../../../
done
