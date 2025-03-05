package main

import (
	"log"
	"net/http"
)

func main() {
	webMain()
}

func webMain() {
	const webPath = "/hinst-website"
	const netAddress = ":8080"
	log.Printf("Starting: netAddress=%v, webPath=%v", netAddress, webPath)
	var webApp = &webApp{webPath: webPath}
	webApp.start()
	assertError(http.ListenAndServe(netAddress, nil))
}
