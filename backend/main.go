package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	var modePointer = flag.String("mode", "web", "mode: web, translate")
	flag.Parse()
	switch *modePointer {
	case "web":
		webMain()
	case "translate":
		translateMain()
	default:
		log.Fatalf("Unknown mode: %v", *modePointer)
	}
}

func webMain() {
	const netAddress = ":8080"
	var webApp = &webApp{}
	webApp.start()
	log.Printf("Starting: netAddress=%v, webPath=%v", netAddress, webApp.webPath)
	assertError(http.ListenAndServe(netAddress, nil))
}

func translateMain() {
	new(translator).run()
}
