#!/bin/bash
#./build.sh 0.0.1
#vscode每次保存会将linux换行符替换为window换行符，如果此文件不能执行，请自行替换换行符
ver=$1
if [ -n "${ver}" ]; then
    echo package version "${ver}"
else
    echo no version param
    exit 1
fi
#打多个平台的包
platforms="linux_amd64 linux_arm windows_amd64"
rm -rf ./resources/output/sqlite3/
for platform in $platforms; do

    GOOS_VAR=$(echo "${platform}" | gawk 'BEGIN{FS="_"} {print $1}')
    GOARCH_VAR=$(echo "${platform}" | gawk 'BEGIN{FS="_"} {print $2}')
    CGO_ENABLED_VAR=1
    echo "${GOOS_VAR}"_"${GOARCH_VAR}"
    if [[ "${GOOS_VAR}" == "windows" ]]; then
        GOOS=${GOOS_VAR} GOARCH=${GOARCH_VAR} CGO_ENABLED=${CGO_ENABLED_VAR} CC=x86_64-w64-mingw32-gcc go build -o ./resources/output/sqlite3/rtsp2rtmp_"${ver}"_"${GOOS_VAR}"_"${GOARCH_VAR}"/rtsp2rtmp.exe main.go
    elif [[ "${platform}" == "linux_arm" ]]; then
        GOOS=${GOOS_VAR} GOARCH=${GOARCH_VAR} CGO_ENABLED=${CGO_ENABLED_VAR} CC=arm-linux-gnueabihf-gcc go build -o ./resources/output/sqlite3/rtsp2rtmp_"${ver}"_"${GOOS_VAR}"_"${GOARCH_VAR}"/rtsp2rtmp main.go
    else
        GOOS=${GOOS_VAR} GOARCH=${GOARCH_VAR} CGO_ENABLED=${CGO_ENABLED_VAR} go build -o ./resources/output/sqlite3/rtsp2rtmp_"${ver}"_"${GOOS_VAR}"_"${GOARCH_VAR}"/rtsp2rtmp main.go
    fi

    mkdir -p ./resources/output/sqlite3/rtsp2rtmp_"${ver}"_"${GOOS_VAR}"_"${GOARCH_VAR}"/resources/output/live
    mkdir -p ./resources/output/sqlite3/rtsp2rtmp_"${ver}"_"${GOOS_VAR}"_"${GOARCH_VAR}"/resources/output/log
    mkdir -p ./resources/output/sqlite3/rtsp2rtmp_"${ver}"_"${GOOS_VAR}"_"${GOARCH_VAR}"/resources/static
    mkdir -p ./resources/output/sqlite3/rtsp2rtmp_"${ver}"_"${GOOS_VAR}"_"${GOARCH_VAR}"/resources/conf
    mkdir -p ./resources/output/sqlite3/rtsp2rtmp_"${ver}"_"${GOOS_VAR}"_"${GOARCH_VAR}"/resources/db

    cp -r ./resources/static/* ./resources/output/sqlite3/rtsp2rtmp_"${ver}"_"${GOOS_VAR}"_"${GOARCH_VAR}"/resources/static
    cp -r ./resources/conf/* ./resources/output/sqlite3/rtsp2rtmp_"${ver}"_"${GOOS_VAR}"_"${GOARCH_VAR}"/resources/conf
    cp -r ./resources/db/* ./resources/output/sqlite3/rtsp2rtmp_"${ver}"_"${GOOS_VAR}"_"${GOARCH_VAR}"/resources/db

    cd ./resources/output/sqlite3/ || exit
    rm -rf rtsp2rtmp_"${ver}"_"${GOOS_VAR}"_"${GOARCH_VAR}".tar.gz
    tar -zcvf ./rtsp2rtmp_"${ver}"_"${GOOS_VAR}"_"${GOARCH_VAR}".tar.gz rtsp2rtmp_"${ver}"_"${GOOS_VAR}"_"${GOARCH_VAR}"/

    # rm -rf ./rtsp2rtmp_"${ver}"_"${GOOS_VAR}"_"${GOARCH_VAR}"/
    cd ../../../
done
