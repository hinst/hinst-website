cd /d "%~dp0.."

rem Generate TypeScript interface definitions from backend
where tygo >nul 2>nul || go install github.com/gzuidhof/tygo@latest
go generate ./...
if errorlevel 1 exit /b 1

go build -v .
exit /b %errorlevel%
