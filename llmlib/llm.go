package llmlib

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ollama/ollama/api"
	cl "github.com/serg-2/libs-go/commonlib"
)

type llmClient struct {
	client               *api.Client
	model                string
	options              map[string]any
	systemRequestMessage string
	requests             cl.ContainerId
}

type request struct {
	finished        bool
	result          string
	startTime       time.Time
	duration        time.Duration
	finishedChannel chan struct{}
}

var availableModels []string = []string{"gemma2:2B", "gemma2:9B", "gemma2:27B", "llava:13b", "llava:34b"}

func InitClient(
	urlString string,
	modelToSet string,
	optionsToSet map[string]any,
	systemRequestMessageToSet string,
) (*llmClient, bool) {
	l := &llmClient{}

	// Set Model
	if !slices.Contains(availableModels, modelToSet) {
		log.Printf("Model %s unsupported.", modelToSet)
		return nil, false
	}

	l.model = modelToSet

	// options part
	l.options = optionsToSet

	// system request message
	l.systemRequestMessage = systemRequestMessageToSet

	// Client part
	serverUrl, err := url.Parse(urlString)
	if err != nil {
		log.Println("Can't parse url!")
		log.Println(err)
		return nil, false
	}
	// Client init
	l.client = api.NewClient(serverUrl, http.DefaultClient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*1000))
	defer cancel()

	// Check Client
	if l.client.Heartbeat(ctx) != nil {
		log.Println("Server is not ok: " + urlString)
		return nil, false
	}

	// Set requests
	l.requests = *cl.NewContainerId()

	// All ok
	return l, true
}

func (l *llmClient) AddRequest(question string) string {
	id := uuid.New().String()

	// Context Part
	ctx := context.Background()

	// Request Part
	streamEnabled := false
	req := &api.ChatRequest{
		Model:    l.model,
		Messages: getMessages(l.systemRequestMessage, question),
		Stream:   &streamEnabled,
		Options:  l.options,
	}

	// Response part
	respFunc := func(resp api.ChatResponse) error {
		tmpVal := l.requests.Get(id).(request)

		tmpVal.finished = true
		tmpVal.result = resp.Message.Content
		tmpVal.duration = time.Now().Sub(tmpVal.startTime)

		l.requests.Add(id, tmpVal)
		return nil
	}

	waitCh := make(chan struct{})
	go func() {
		err := l.client.Chat(ctx, req, respFunc)
		if err != nil {
			log.Println("Error in Chat handling")
			log.Println(err)
			tmpReq := l.requests.Get(id)
			if tmpReq == nil {
				log.Fatalln("UNSUPPORTED!")
			}
			tmpVal := tmpReq.(request)

			tmpVal.finished = true
			tmpVal.result = "Error in Chat handling: " + err.Error()
			tmpVal.duration = time.Now().Sub(tmpVal.startTime)

			l.requests.Add(id, tmpVal)
		}

		close(waitCh)
	}()

	var newReq request = request{
		false,
		"answer is not ready",
		time.Now(),
		0,
		waitCh,
	}
	l.requests.Add(id, newReq)

	return id
}

func (l *llmClient) CheckRequest(id string) bool {
	tmpVal := l.requests.Get(id)
	if tmpVal == nil {
		return false
	}
	tmpReq := tmpVal.(request)
	return tmpReq.finished
}

func (l *llmClient) GetAnswer(id string) string {
	tmpReq := l.requests.Get(id)
	if tmpReq == nil {
		return "No such id."
	}
	tmpVal := tmpReq.(request)

	return strings.TrimSuffix(tmpVal.result, "\n")
}

func (l *llmClient) GetFinishChannel(id string) *chan struct{} {
	tmpReq := l.requests.Get(id)
	if tmpReq == nil {
		return nil
	}
	tmpVal := tmpReq.(request)

	return &tmpVal.finishedChannel
}

func (l *llmClient) GetCompletedFor(id string) time.Duration {
	tmpReq := l.requests.Get(id)
	if tmpReq == nil {
		return 0
	}
	tmpVal := tmpReq.(request)
	return tmpVal.duration
}

func (l *llmClient) GetDurationStatus(id string) time.Duration {
	tmpReq := l.requests.Get(id)
	if tmpReq == nil {
		return 0
	}
	tmpVal := tmpReq.(request)
	return time.Now().Sub(tmpVal.startTime)
}

// local function to get messages array using different roles
func getMessages(systemRequestMessage string, question string) []api.Message {
	messages := []api.Message{
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

		{
			Role:    "system",
			Content: systemRequestMessage,
		},
		{
			Role:    "user",
			Content: question,
		},
	}
	return messages
}
