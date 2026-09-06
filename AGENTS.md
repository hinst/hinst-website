When searching files, avoid launching `find /` in root filesystem, because root filesystem has millions of files.
When listing files recursively, avoid `node_modules` trap, otherwise thousands of files from `node_modules` will spam the output.
