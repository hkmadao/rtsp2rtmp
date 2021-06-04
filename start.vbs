Dim WinScriptHost
Set WinScriptHost = CreateObject("WScript.Shell")
WinScriptHost.Run Chr(34) & "rtsp2rtmp.exe" & Chr(34), 0
Set WinScriptHost = Nothing