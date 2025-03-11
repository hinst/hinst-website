package main

type libreTranslateRequest struct {
	Query  string `json:"q"`
	Source string `json:"source"`
	Target string `json:"target"`
	Format string `json:"format"`
}

type libreTranslateResponse struct {
	TranslatedText string `json:"translatedText"`
}
