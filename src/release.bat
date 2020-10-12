rsrc -manifest app.manifest -ico=app.ico -o rsrc.syso
go build -ldflags="-H windowsgui" -o "..\\MWFileChecker.exe"
