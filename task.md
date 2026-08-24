Save schema.yaml to disk after every build.

# Documentation

To save your OpenAPI schema as a schema.yaml file on each build, you can use a custom Go binary or template function that executes during go generate or a Makefile build step. The [Huma REST framework](https://huma.rocks/) processes its schema dynamically at runtime, meaning you must programmatically fetch the schema from your API instance and write it to disk.
Here is the most reliable approach using go generate.
## 1. Create a Schema Generator File
Create a small utility file (e.g., cmd/openapi/main.go) that instantiates your Huma API, marshals the OpenAPI spec to YAML, and writes it to your desired file location.

package main
import (
	"context"
	"fmt"
	"os"

	"://github.com"
	"://github.com/adapters/humachi"
	"://github.com"
)
func main() {
	// 1. Create a dummy router and your Huma API instance
	// Ensure this mirrors your production configuration setup
	router := chi.NewRouter()
	api := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))

	// 2. (Optional) Register your routes if you need route-specific schemas included
	// setupRoutes(api)

	// 3. Marshal the internal OpenAPI structure to YAML
	yamlBytes, err := api.OpenAPI().MarshalYAML()
	if err != nil {
		fmt.Printf("Failed to marshal OpenAPI to YAML: %v\n", err)
		os.Exit(1)
	}

	// 4. Save the schema file to disk
	filename := "schema.yaml"
	err = os.WriteFile(filename, yamlBytes, 0644)
	if err != nil {
		fmt.Printf("Failed to write %s: %v\n", filename, err)
		os.Exit(1)
	}

	fmt.Printf("Successfully saved schema to %s\n", filename)
}

## 2. Connect via go generate
Add a go:generate directive at the top of your main application entry point (e.g., main.go). This ensures that running a single command updates your schema.

package main
//go:generate go run ./cmd/openapi/main.go
func main() {
    // Your standard server execution logic goes here
}

## 3. Automate on Every Build
To enforce that the schema.yaml updates on every build, you have two options depending on your workflow:
## Option A: Using a Makefile (Recommended)
If you use a Makefile for your build pipeline, chain the generation step directly into your build target.

.PHONY: build
build:
	go generate ./...
	go build -o bin/server main.go

## Option B: Go Build Constraint
Alternatively, you can just run go generate ./... in your CI/CD pipeline right before executing your standard go build command.
------------------------------
If you would like to proceed further, let me know:

*
* Which HTTP router adapter (Chi, Gin, Echo, Fiber) you are using.
* If you want to bundle this schema.yaml back into your final Go binary using go:embed.
* How your project routes are structured, so we can ensure they load properly into the generator tool.
*


