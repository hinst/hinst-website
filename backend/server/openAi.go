package server

const AI_REASONING_EFFORT_MEDIUM = "medium"

type openAiRequest struct {
	Model           string          `json:"model"`
	Messages        []openAiMessage `json:"messages"`
	Stream          bool            `json:"stream"`
	ReasoningEffort string          `json:"reasoning_effort"`
}

const AI_ROLE_SYSTEM = "system"
const AI_ROLE_USER = "user"

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

const OLLAMA_DEFAULT_PORT = 11434
