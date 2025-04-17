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
) (*llmClient, bool) {
	l := &llmClient{}

	// Set Model
	if !slices.Contains(availableModelsDS, modelToSet) {
		log.Printf("Model %s unsupported.", modelToSet)
		return nil, false
	}

	if modelToSet == "chat" {
		l.modelDS = deepseek.DEEPSEEK_CHAT_MODEL
	} else {
		panic("1")
	}

	// system request message
	l.systemRequestMessageDS = getApiMessagesDS(systemRequestMessages)

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

func getRequestDS(l *llmClient, question string) *dsr.ChatCompletionsRequest {
	streamEnabled := false
	return &dsr.ChatCompletionsRequest{
		Model:  l.modelDS,
		Stream: streamEnabled,
		Messages: getMessagesDS(
			l.systemRequestMessageDS,
			question,
		),
	}
}

// local function to get messages array using different roles
func getMessagesDS(environmentMessages []*dsr.Message, question string) []*dsr.Message {
	messages := append(environmentMessages, &dsr.Message{
		Role:    "user",
		Content: question,
	},
	)
	return messages
}
