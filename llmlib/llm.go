package llmlib

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/go-deepseek/deepseek"
	dsr "github.com/go-deepseek/deepseek/request"
	"github.com/google/uuid"
	"github.com/ollama/ollama/api"
	cl "github.com/serg-2/libs-go/commonlib"
)

type llmClient struct {
	clientOllama               *api.Client
	clientDS                   deepseek.Client
	modelOllama                string
	modelDS                    string
	options                    map[string]any
	systemRequestMessageOllama []api.Message
	systemRequestMessageDS     []*dsr.Message
	requests                   cl.ContainerId
}

type request struct {
	finished        bool
	result          string
	startTime       time.Time
	duration        time.Duration
	finishedChannel chan struct{}
}

type SystemMessages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (l *llmClient) AddRequest(question string) string {
	id := uuid.New().String()

	// Context Part
	ctx := context.Background()

	waitCh := make(chan struct{})

	if l.clientOllama != nil {
		// OLLAMA
		go func(chatRequest *api.ChatRequest, responseFunction func(resp api.ChatResponse) error) {
			err := l.clientOllama.Chat(ctx, chatRequest, responseFunction)
			if err != nil {
				log.Println("Error in Chat handling")
				log.Println(err)
				tmpVal := l.requests.Get(id).(request)
				tmpVal.finished = true
				tmpVal.result = "Error in Chat handling: " + err.Error()
				tmpVal.duration = time.Now().Sub(tmpVal.startTime)
				l.requests.Add(id, tmpVal)
			}
			close(waitCh)
		}(
			getRequestOllama(l, question),
			getResonseFunctionOllama(l, id),
		)
	} else if l.clientDS != nil {
		go func(chatRequest *dsr.ChatCompletionsRequest) {
			chatResp, err := l.clientDS.CallChatCompletionsChat(ctx, chatRequest)
			tmpVal := l.requests.Get(id).(request)
			if err != nil {
				log.Println("Error in Chat handling DS")
				log.Println(err)
				tmpVal.result = "Error in Chat handling: " + err.Error()
			} else {
				tmpVal.result = chatResp.Choices[0].Message.Content
			}
			tmpVal.duration = time.Now().Sub(tmpVal.startTime)
			tmpVal.finished = true
			l.requests.Add(id, tmpVal)
			close(waitCh)
		}(
			getRequestDS(l, question),
		)
	} else {
		panic('2')
	}

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
