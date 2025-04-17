package llmlib

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/ollama/ollama/api"
	cl "github.com/serg-2/libs-go/commonlib"
)

var availableModelsOllama []string = []string{
	"gemma3:12B", "gemma3:27B",
	"llama3.3:70B",
	"deepseek-r1:671B",
	"llava:13b", "llava:34b",
}

func InitOllamaClient(
	urlString string,
	modelToSet string,
	optionsToSet map[string]any,
	systemRequestMessages []SystemMessages,
) (*llmClient, bool) {
	l := &llmClient{}

	// Set Model
	if !slices.Contains(availableModelsOllama, modelToSet) {
		log.Printf("Model %s unsupported.", modelToSet)
		return nil, false
	}

	l.modelOllama = modelToSet

	// options part
	l.options = optionsToSet

	// system request message
	l.systemRequestMessageOllama = getApiMessagesOllama(systemRequestMessages)

	// Client part
	serverUrl, err := url.Parse(urlString)
	if err != nil {
		log.Println("Can't parse url!")
		log.Println(err)
		return nil, false
	}
	// Client init
	l.clientOllama = api.NewClient(serverUrl, http.DefaultClient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*1000))
	defer cancel()

	// Check Client
	if l.clientOllama.Heartbeat(ctx) != nil {
		log.Println("Server is not ok: " + urlString)
		return nil, false
	}

	// Set requests
	l.requests = *cl.NewContainerId()

	// All ok
	return l, true
}

func getApiMessagesOllama(systemRequestMessages []SystemMessages) []api.Message {
	var result []api.Message
	for _, message := range systemRequestMessages {
		result = append(result, api.Message{
			Role:    message.Role,
			Content: message.Content,
		})
	}
	return result
}

func getRequestOllama(l *llmClient, question string) *api.ChatRequest {
	streamEnabled := false
	return &api.ChatRequest{
		Model: l.modelOllama,
		Messages: getMessagesOllama(
			l.systemRequestMessageOllama,
			question,
		),
		Stream:  &streamEnabled,
		Options: l.options,
	}
}

func getResonseFunctionOllama(l *llmClient, id string) func(resp api.ChatResponse) error {
	return func(resp api.ChatResponse) error {
		tmpVal := l.requests.Get(id).(request)

		tmpVal.finished = true
		tmpVal.result = resp.Message.Content
		tmpVal.duration = time.Now().Sub(tmpVal.startTime)

		l.requests.Add(id, tmpVal)
		return nil
	}
}

// local function to get messages array using different roles
func getMessagesOllama(environmentMessages []api.Message, question string) []api.Message {
	// api.Message{
	// 	Role:    "system",
	// 	Content: "Provide very brief, concise responses",
	// },
	// api.Message{
	// 	Role:    "user",
	// 	Content: "Name some unusual animals",
	// },
	// api.Message{
	// 	Role:    "assistant",
	// 	Content: "Monotreme, platypus, echidna",
	// },
	// api.Message{
	// 	Role:    "user",
	// 	Content: "which of these is the most dangerous?",
	// },

	messages := append(environmentMessages, api.Message{
		Role:    "user",
		Content: question,
	},
	)
	return messages
}
