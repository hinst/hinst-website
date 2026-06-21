package server

const lm_studio_role_system = "system"
const lm_studio_role_user = "user"
const lm_studio_multilingual_model_id = "aya-expanse-8B"
const lm_studio_default_url = "http://localhost:1235/v1/chat/completions"

type lmStudioRequest struct {
	Model    string            `json:"model"`
	Messages []lmStudioMessage `json:"messages"`
	Stream   bool              `json:"stream"`
}

type lmStudioMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type lmStudioResponse struct {
	Choices []lmStudioChoice `json:"choices"`
}

type lmStudioChoice struct {
	Message lmStudioMessage
}
