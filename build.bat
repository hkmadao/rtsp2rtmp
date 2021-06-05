:: Linux
SET CGO_ENABLED=1
SET GOOS=linux
SET GOARCH=amd64
go build -o rtsp2rtmp main.go

:: Windows
SET CGO_ENABLED=1
SET GOOS=windows
SET GOARCH=amd64
go build -o rtsp2rtmp.exe main.go

pause