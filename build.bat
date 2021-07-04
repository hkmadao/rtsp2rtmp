@echo off
chcp 65001
set /p ver=请输入版本：  
echo 版本：%ver% 打包开始
rmdir /S /Q .\resources\output\releases
@REM windows_amd64
echo 打包windows_amd64平台
SET GOOS=windows
SET GOARCH=amd64
SET CGO_ENABLED=0
go build -o .\resources\output\releases\rtsp2rtmp_%ver%_%GOOS%_%GOARCH%\rtmp2flv.exe main.go
echo =============%GOOS%_%GOARCH%
md .\resources\output\releases\rtsp2rtmp_%ver%_%GOOS%_%GOARCH%\output\live
md .\resources\output\releases\rtsp2rtmp_%ver%_%GOOS%_%GOARCH%\output\log
md .\resources\output\releases\rtsp2rtmp_%ver%_%GOOS%_%GOARCH%\conf

xcopy /S /Y /E .\resources\static .\resources\output\releases\rtsp2rtmp_%ver%_%GOOS%_%GOARCH%\static\
xcopy /S /Y /E .\resources\db .\resources\output\releases\rtsp2rtmp_%ver%_%GOOS%_%GOARCH%\db\
xcopy /S /Y /E .\resources\conf .\resources\output\releases\rtsp2rtmp_%ver%_%GOOS%_%GOARCH%\conf
cd .\resources\output\releases\
7z a -ttar -so rtsp2rtmp_%ver%_%GOOS%_%GOARCH%_demo.tar rtsp2rtmp_%ver%_%GOOS%_%GOARCH%/ | 7z a -si rtsp2rtmp_%ver%_%GOOS%_%GOARCH%_demo.tar.gz
cd ..\..\..\

pause