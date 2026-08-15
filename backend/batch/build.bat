@echo off
setlocal
cd /d "%~dp0.."

rem Ensure tygo is installed (needed by the //go:generate directive in main.go)
where tygo >nul 2>nul || go install github.com/gzuidhof/tygo@latest
for /f "tokens=*" %%g in ('go env GOPATH') do set "GOPATH=%%g"
set "PATH=%PATH%;%GOPATH%\bin"

rem Regenerate TypeScript types from Go source, then build
go generate ./...
if errorlevel 1 exit /b 1
go build -v .
exit /b %errorlevel%
