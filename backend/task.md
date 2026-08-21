Check commit 79a2d2994b7313c36c33ad9974f269362a113076 to see what we have been doing.

Please proceed with the next step: migrate API endpoint "/api/goal/image"

Example: how to return image from Huma Rest framework.

```go
package main

import (
	"context"
	"net/http"

	"://github.com"
	"://github.com/adapters/humago"
)

// 1. Define the output struct with a []byte body field.
type ImageOutput struct {
	Body []byte
}

func main() {
	// Initialize standard library router (Go 1.22+)
	mux := http.NewServeMux()
	api := humago.New(mux, huma.DefaultConfig("Image API", "1.0.0"))

	// 2. Register the operation with custom Content-Type metadata
	huma.Register(api, huma.Operation{
		OperationID: "get-image",
		Method:      http.MethodGet,
		Path:        "/image",
		Summary:     "Returns a sample binary PNG image",
		Responses: map[string]*huma.Response{
			"200": {
				Description: "Image returned successfully",
				Content: map[string]*huma.MediaType{
					"image/png": {}, // Informs OpenAPI/Swagger of the correct format
				},
			},
		},
	}, func(ctx context.Context, input *struct{}) (*ImageOutput, error) {

		// 3. Replace this with your actual image reading or generation logic
		// (e.g., using os.ReadFile("avatar.png") or image.Encode())
		dummyPngBytes := []byte{
			0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // Minimal PNG header
			0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		}

		// 4. Return the raw bytes wrapped in your output struct
		return &ImageOutput{Body: dummyPngBytes}, nil
	})

	http.ListenAndServe(":8888", mux)
}

```
