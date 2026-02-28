package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	var file, fileError = os.Open("./go.mod")
	if fileError != nil {
		panic(fileError)
	}
	defer file.Close()

	var scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		var line = scanner.Text()
		if strings.HasPrefix(line, "replace ") {
			continue
		}
		fmt.Println(line)
	}
}
