package server

import (
	"log"
	"os"
)

// generateSchema registers all API routes and writes the OpenAPI schema
// to the given file. No database connection is required, because routes
// only use the database at request time.
func (me *program) generateSchema(filename string) {
	var webApp = new(webApp)
	webApp.init(me.database)
	var yaml, err = webApp.webApi.OpenAPI().YAML()
	if err != nil {
		log.Fatalf("Failed to marshal OpenAPI to YAML: %v", err)
	}
	if err = os.WriteFile(filename, yaml, 0644); err != nil {
		log.Fatalf("Failed to write %v: %v", filename, err)
	}
	log.Printf("Successfully saved OpenAPI schema to %v", filename)
}
