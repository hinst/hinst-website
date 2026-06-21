package server

const AI_ROLE_SYSTEM = "system"
const AI_ROLE_USER = "user"

type openAiRequest struct {
	Model    string          `json:"model"`
	Messages []openAiMessage `json:"messages"`
	Stream   bool            `json:"stream"`
}

type openAiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAiResponse struct {
	Choices []lmStudioChoice `json:"choices"`
}

type lmStudioChoice struct {
	Message openAiMessage
}
