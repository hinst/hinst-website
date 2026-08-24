cd /d "%~dp0.."

go build -v .
if errorlevel 1 exit /b %errorlevel%

hinst-website.exe --mode=generateSchema
exit /b %errorlevel%
