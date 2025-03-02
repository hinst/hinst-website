package main

import "net/http"

func main() {
	const listenAddress = ":8080"
	const path = "hinst-website"
	(&webApp{
		path: path,
	}).start()
	assertError(http.ListenAndServe(listenAddress, nil))
}
