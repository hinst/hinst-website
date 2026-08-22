See file: `server\webAppGoals.go`
Function `setGoalPostText` currently uses `humago.Unwrap()` to read input text.
However we should read input text differently: using this example

```go
huma.Register(api, huma.Operation{
	OperationID: "post-plain-text",
	Method:      http.MethodPost,
	Path:        "/text",
	Summary:     "Example to post plain text input",
}, func(ctx context.Context, input *struct {
	RawBody []byte `contentType:"text/plain"`
}) (*struct{}, error) {
	fmt.Println("Got input:", input.RawBody)
	return nil, nil
}
```
