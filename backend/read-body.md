See file: `server\webAppGoals.go`
Function `setGoalPostText` currently uses `humago.Unwrap()` to read input text.
However we should read input text differently: using this example

```go
package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com"
	"github.com/adapters/humago"
)

// Define the input structure requesting the raw text body
type TextRequest struct {
	// RawBody captures the unparsed request payload
	RawBody []byte `contentType:"text/plain"`
}

// Define a simple response structure
type TextResponse struct {
	Body struct {
		Message string `json:"message"`
	}
}

func main() {
	// Create a standard Go HTTP router (Go 1.22+)
	mux := http.NewServeMux()

	// Initialize the Huma API config
	config := huma.DefaultConfig("Text Reader API", "1.0.0")
	api := humago.New(mux, config)

	// Register the text endpoint
	huma.Register(api, huma.Operation{
		OperationID: "post-plain-text",
		Method:      http.MethodPost,
		Path:        "/submit-text",
		Summary:     "Accepts a plain text payload",
	}, func(ctx context.Context, input *TextRequest) (*TextResponse, error) {
		// Convert the raw byte slice into a readable Go string
		receivedText := string(input.RawBody)

		// Print to server logs for verification
		fmt.Printf("Received plain text: %s\n", receivedText)

		// Respond back to the client
		res := &TextResponse{}
		res.Body.Message = fmt.Sprintf("Successfully processed %d characters.", len(receivedText))
		return res, nil
	})

	// Start the server
	fmt.Println("Server running on :8080...")
	http.ListenAndServe(":8080", mux)
}
```
