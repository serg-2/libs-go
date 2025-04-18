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
) (*LLMClient, bool) {
	l := &LLMClient{}

	// Set Model
	if !slices.Contains(availableModelsOllama, modelToSet) {
		log.Printf("Model %s unsupported.", modelToSet)
		return nil, false
	}

	l.model = modelToSet

	// options part
	l.options = optionsToSet

	// system request message
	l.systemRequestMessages = systemRequestMessages

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

func getRequestOllama(l *LLMClient, question string, previosMessages []SystemMessages) *api.ChatRequest {
	streamEnabled := false
	
	requestMessages := getApiMessagesOllama(l.systemRequestMessages)
	
	for _, prevMessage := range previosMessages {
		requestMessages = append(requestMessages, api.Message{
			Role:    prevMessage.Role,
			Content: prevMessage.Content,
		})
	}

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

	requestMessages = append(requestMessages, api.Message{
		Role:    "user",
		Content: question,
	})
	
	return &api.ChatRequest{
		Model:    l.model,
		Messages: requestMessages,
		Stream:   &streamEnabled,
		Options:  l.options,
	}
}

func getResonseFunctionOllama(l *LLMClient, id string) func(resp api.ChatResponse) error {
	return func(resp api.ChatResponse) error {
		tmpVal := l.requests.Get(id).(request)

		tmpVal.finished = true
		tmpVal.result = resp.Message.Content
		tmpVal.duration = time.Now().Sub(tmpVal.startTime)

		l.requests.Add(id, tmpVal)
		return nil
	}
}
