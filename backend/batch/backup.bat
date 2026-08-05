if exist saved-goals\backup rmdir /s /q saved-goals\backup &&^
go build &&^
hinst-website.exe --mode=backup
