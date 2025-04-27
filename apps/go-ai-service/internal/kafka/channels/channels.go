package kafka

import "github.com/Vovarama1992/go-ai-service/pkg/types"

var AutoReplyInput = make(chan types.EnhancedAutoreplyRequest, 100)
var AutoReplyOutput = make(chan types.AiAutoreplyResponse, 100)

var AdviceInput = make(chan types.EnhancedAdviceRequest, 100)
var AdviceOutput = make(chan types.AiAdviceResponse, 100)
