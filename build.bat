@echo off
chcp 65001
set /p ver=请输入版本：  
echo 版本：%ver% 打包开始
@REM window_amd64
@REM SET GOOS=windows
@REM SET GOARCH=amd64
SET CGO_ENABLED=1
go build -o rtsp2rtmp_%ver%_window_amd64.exe main.go

@REM window_amd64
@REM rmdir /S /Q .\output\rtsp2rtmp_%ver%_window_amd64

@REM md .\output\rtsp2rtmp_%ver%_window_amd64\output\live
@REM md .\output\rtsp2rtmp_%ver%_window_amd64\output\log
@REM md .\output\rtsp2rtmp_%ver%_window_amd64\conf

@REM xcopy /S /Y /E .\static .\output\rtsp2rtmp_%ver%_window_amd64\static\
@REM xcopy /S /Y /E .\db .\output\rtsp2rtmp_%ver%_window_amd64\db\
@REM xcopy .\conf\conf.yml .\output\rtsp2rtmp_%ver%_window_amd64\conf
@REM xcopy .\rtsp2rtmp_%ver%_window_amd64.exe .\output\rtsp2rtmp_%ver%_window_amd64\rtsp2rtmp.exe
@REM xcopy .\start.vbs .\output\rtsp2rtmp_%ver%_window_amd64\start.vbs

@REM linux_adm64
@REM rmdir /S /Q .\output\rtsp2rtmp_%ver%_linux_amd64

@REM md .\output\rtsp2rtmp_%ver%_linux_amd64\output\live
@REM md .\output\rtsp2rtmp_%ver%_linux_amd64\output\log
@REM md .\output\rtsp2rtmp_%ver%_linux_amd64\conf

@REM xcopy /S /Y /E .\static .\output\rtsp2rtmp_%ver%_linux_amd64\static\
@REM xcopy /S /Y /E .\db .\output\rtsp2rtmp_%ver%_linux_amd64\db\
@REM xcopy .\conf\conf.yml .\output\rtsp2rtmp_%ver%_linux_amd64\conf
@REM xcopy .\rtsp2rtmp_%ver%_linux_amd64 .\output\rtsp2rtmp_%ver%_linux_amd64\rtsp2rtmp

@REM linux_armv6
@REM rmdir /S /Q .\output\rtsp2rtmp_%ver%_linux_armv6

@REM md .\output\rtsp2rtmp_%ver%_linux_armv6\output\live
@REM md .\output\rtsp2rtmp_%ver%_linux_armv6\output\log
@REM md .\output\rtsp2rtmp_%ver%_linux_armv6\conf

@REM xcopy /S /Y /E .\static .\output\rtsp2rtmp_%ver%_linux_armv6\static\
@REM xcopy /S /Y /E .\db .\output\rtsp2rtmp_%ver%_linux_armv6\db\
@REM xcopy .\conf\conf.yml .\output\rtsp2rtmp_%ver%_linux_armv6\conf
@REM xcopy .\rtsp2rtmp_%ver%_linux_armv6 .\output\rtsp2rtmp_%ver%_linux_armv6\rtsp2rtmp

pause