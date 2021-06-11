@REM Windows
SET GOOS=windows
SET GOARCH=amd64
SET CGO_ENABLED=1
go build -o rtsp2rtmp.exe main.go

rmdir /S /Q .\output\rtsp2rtmp

md .\output\rtsp2rtmp\output\live
md .\output\rtsp2rtmp\output\log
md .\output\rtsp2rtmp\conf

xcopy /S /Y /E .\static .\output\rtsp2rtmp\static\
xcopy /S /Y /E .\db .\output\rtsp2rtmp\db\
xcopy .\conf\conf.yml .\output\rtsp2rtmp\conf
xcopy .\rtsp2rtmp .\output\rtsp2rtmp
xcopy .\rtsp2rtmp.exe .\output\rtsp2rtmp
xcopy .\start.vbs .\output\rtsp2rtmp

pause