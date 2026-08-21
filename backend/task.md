Check commit 79a2d2994b7313c36c33ad9974f269362a113076 to see what we have been doing.

Please proceed with the next step: migrate API endpoint "/api/goal/image"

Example: how to return image from Huma Rest framework.

```go
package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"
)

// Options for the CLI.
type Options struct {
	Port int `help:"Port to listen on" short:"p" default:"8888"`
}

// ImageOutput represents the image operation response.
type ImageOutput struct {
	ContentType string `header:"Content-Type"`
	Body        []byte
}

func main() {
	// Create a CLI app which takes a port option.
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		// Create a new router & API
		router := chi.NewMux()
		api := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))

		// Register GET /image
		huma.Register(api, huma.Operation{
			OperationID: "get-image",
			Summary:     "Get an image",
			Method:      http.MethodGet,
			Path:        "/image",
			Responses: map[string]*huma.Response{
				"200": {
					Description: "Image response",
					Content: map[string]*huma.MediaType{
						"image/jpeg": {},
					},
				},
			},
		}, func(ctx context.Context, input *struct{}) (*ImageOutput, error) {
			resp := &ImageOutput{}
			resp.ContentType = "image/png"
			resp.Body = []byte{ /* ... image bytes here ... */ }
			return resp, nil
		})

		// Tell the CLI how to start your server.
		hooks.OnStart(func() {
			fmt.Printf("Starting server on port %d...\n", options.Port)
			http.ListenAndServe(fmt.Sprintf(":%d", options.Port), router)
		})
	})

	// Run the CLI. When passed no commands, it starts the server.
	cli.Run()
}

```
