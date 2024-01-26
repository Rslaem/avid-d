@echo off
setlocal enabledelayedexpansion
set /a n=%1-1
for /l %%i in (0,1,%n%) do (
    start cmd /k go run main.go -n=%1 -f=%2 -id=%%i -path=%~dp0
)