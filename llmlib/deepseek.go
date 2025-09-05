package llmlib

import (
	"log"
	"slices"

	"github.com/cohesion-org/deepseek-go"
	cl "github.com/serg-2/libs-go/commonlib"
)

var availableModelsDS []string = []string{
	"chat",
}

func InitDSClient(
	modelToSet string,
	systemRequestMessages []SystemMessages,
	apiKey string,
) (*LLMClient, bool) {
	l := &LLMClient{}

	// Set Model
	if !slices.Contains(availableModelsDS, modelToSet) {
		log.Printf("Model %s unsupported.", modelToSet)
		return nil, false
	}

	if modelToSet == "chat" {
		l.model = deepseek.DeepSeekChat
	} else {
		log.Println("Can't understand DS model to set.")
		return nil, false
	}

	// system request message
	l.systemRequestMessages = systemRequestMessages

	// Client init
	l.clientDS = deepseek.NewClient(apiKey)

	// Set requests
	l.requests = *cl.NewContainerId()

	// All ok
	return l, true
}

func getApiMessagesDS(systemRequestMessages []SystemMessages) []deepseek.ChatCompletionMessage {
	var result []deepseek.ChatCompletionMessage
	for _, message := range systemRequestMessages {
		result = append(result, deepseek.ChatCompletionMessage{
			Role:       message.Role,
			Content:    message.Content,
			ToolCallID: message.ToolCallId,
		})
	}
	return result
}

func getRequestDS(
	l *LLMClient,
	question string,
	previosMessages []SystemMessages,
	tools []deepseek.Tool,
) *deepseek.ChatCompletionRequest {
	requestMessages := getApiMessagesDS(l.systemRequestMessages)
	// Add Previous
	for _, prevMessage := range previosMessages {
		requestMessages = append(requestMessages, deepseek.ChatCompletionMessage{
			Role:    prevMessage.Role,
			Content: prevMessage.Content,
		})
	}
	// Add question
	requestMessages = append(requestMessages, deepseek.ChatCompletionMessage{
		Role:    "user",
		Content: question,
	})
	return &deepseek.ChatCompletionRequest{
		Model:    l.model,
		Messages: requestMessages,
		Tools:    tools,
	}
}
