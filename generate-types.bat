@echo off
rem Generates TypeScript types from the Golang db_objects package using tygo.
rem Requires Go. tygo is installed automatically if missing.
setlocal
cd /d "%~dp0backend"
where tygo >nul 2>nul || go install github.com/gzuidhof/tygo@latest
tygo generate
endlocal
