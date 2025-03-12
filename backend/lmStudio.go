package main

const lm_studio_role_system = "system"
const lm_studio_role_user = "user"
const prompt_russian_to_english = "You are a professional Russian-to-English translator specializing in diary blog posts. Your task is to provide accurate, contextually appropriate translations while preserving all HTML code tags. Do not provide any explanations or commentary - just the direct translation."

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
