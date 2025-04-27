package types

type AiAutoreplyRequest struct {
	ChatId       int    `json:"chatId"`
	TargetUserId int    `json:"targetUserId"`
	Text         string `json:"text"`
	SenderId     int    `json:"senderId"`
	CustomPrompt string `json:"customPrompt,omitempty"`
}

type AiAdviceRequest struct {
	ChatId       int    `json:"chatId"`
	TargetUserId int    `json:"targetUserId"`
	SenderId     int    `json:"senderId"`
	SourceText   string `json:"sourceText"`
	CustomPrompt string `json:"customPrompt,omitempty"`
}

type AiAdviceResponse struct {
	ChatId       int    `json:"chatId"`
	TargetUserId int    `json:"targetUserId"`
	Advice       string `json:"advice"`
	SourceText   string `json:"sourceText"`
}

type AiAutoreplyResponse struct {
	ChatId       int    `json:"chatId"`
	TargetUserId int    `json:"targetUserId"`
	Text         string `json:"text"`
	SenderId     int    `json:"senderId"`
}

type EnhancedAdviceRequest struct {
	Request  AiAdviceRequest
	ThreadId string
	UserName string
}

type EnhancedAutoreplyRequest struct {
	Request  AiAutoreplyRequest
	ThreadId string
	UserName string
}
