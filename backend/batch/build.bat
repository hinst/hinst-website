cd /d "%~dp0.."

go build -v .
exit /b %errorlevel%
