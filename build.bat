@echo off

go build -v -ldflags "-H windowsgui" main.go

main.exe

pause
