package llmlib

import (
	"log"
	"slices"

	"github.com/go-deepseek/deepseek"
	dsr "github.com/go-deepseek/deepseek/request"
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
		l.model = deepseek.DEEPSEEK_CHAT_MODEL
	} else {
		log.Println("Can't understand DS model to set.")
		return nil, false
	}

	// system request message
	l.systemRequestMessages = systemRequestMessages

	// Client init
	var err error
	l.clientDS, err = deepseek.NewClient(apiKey)
	if err != nil {
		log.Println("Can't get DS client")
		log.Println(err)
		return nil, false
	}

	// Set requests
	l.requests = *cl.NewContainerId()

	// All ok
	return l, true
}

func getApiMessagesDS(systemRequestMessages []SystemMessages) []*dsr.Message {
	var result []*dsr.Message
	for _, message := range systemRequestMessages {
		result = append(result, &dsr.Message{
			Role:    message.Role,
			Content: message.Content,
		})
	}
	return result
}

func getRequestDS(l *LLMClient, question string, previosMessages []SystemMessages) *dsr.ChatCompletionsRequest {
	streamEnabled := false

	requestMessages := getApiMessagesDS(l.systemRequestMessages)
	// Add Previous
	for _, prevMessage := range previosMessages {
		requestMessages = append(requestMessages, &dsr.Message{
			Role:    prevMessage.Role,
			Content: prevMessage.Content,
		})
	}
	// Add question
	requestMessages = append(requestMessages, &dsr.Message{
		Role:    "user",
		Content: question,
	})
	return &dsr.ChatCompletionsRequest{
		Model:    l.model,
		Stream:   streamEnabled,
		Messages: requestMessages,
	}
}
