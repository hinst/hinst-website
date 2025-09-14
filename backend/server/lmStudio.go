package server

const lm_studio_role_system = "system"
const lm_studio_role_user = "user"
const lm_studio_multilingual_model_id = "aya-expanse-8B"
const lm_studio_default_url = "http://localhost:1235/v1/chat/completions"
const prompt_Russian_to_something = "You are a professional Russian-to-{something} translator specializing in diary blog posts. Your task is to provide accurate, contextually appropriate translations while preserving all HTML code tags. Do not provide any explanations or commentary - just the direct translation."
const prompt_generate_title = "You are a professional blog post editor. Your task is to create an engaging title for the following text. The title should be one sentence. Ignore all HTML code tags. Do not provide any explanations or commentary - just the title in the same language as the original text."

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
