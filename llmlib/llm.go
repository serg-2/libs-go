package llmlib

import (
	"context"
	"log"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-deepseek/deepseek"
	dsr "github.com/go-deepseek/deepseek/request"
	"github.com/google/uuid"
	"github.com/ollama/ollama/api"
	cl "github.com/serg-2/libs-go/commonlib"
	js "github.com/serg-2/libs-go/jsonlib"
)

type LLMClient struct {
	clientOllama          *api.Client
	clientDS              deepseek.Client
	model                 string
	options               map[string]any
	systemRequestMessages []SystemMessages
	requests              cl.ContainerId
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

var availableRoles []string = []string{
	"system",
	"user",
	"assistant",
	"tool",
}

func (l *LLMClient) AddRequest(question string, previosMessages []SystemMessages) string {
	// Validate system messages
	if !validateSystemMessages(previosMessages) {
		log.Println("Can't validate request.")
		return ""
	}

	log.Printf("Additional context questions: %s\n", js.JsonAsString(previosMessages))

	// Validate question
	if !validateQuestion(question) {
		log.Println("Can't validate question.")
		return ""
	}

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
			getRequestOllama(l, question, previosMessages),
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
			getRequestDS(l, question, previosMessages),
		)
	} else {
		log.Println("Can't find clients.")
		return ""
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

func validateQuestion(question string) bool {
	if utf8.RuneCountInString(question) > 300 {
		log.Printf("Too big user request message!\n")
		return false
	}
	return true
}

func validateSystemMessages(previosMessages []SystemMessages) bool {
	if len(previosMessages) > 30 {
		log.Printf("Too big context!\n")
		return false
	}
	for _, message := range previosMessages {
		// Check role
		if !slices.Contains(availableRoles, message.Role) {
			log.Printf("Role %s unsupported.\n", message.Role)
			return false
		}
		// Check length?
	}
	return true
}

func (l *LLMClient) CheckRequest(id string) bool {
	tmpVal := l.requests.Get(id)
	if tmpVal == nil {
		return false
	}
	tmpReq := tmpVal.(request)
	return tmpReq.finished
}

func (l *LLMClient) GetAnswer(id string) string {
	tmpReq := l.requests.Get(id)
	if tmpReq == nil {
		return "No such id."
	}
	tmpVal := tmpReq.(request)
	return strings.TrimSuffix(tmpVal.result, "\n")
}

func (l *LLMClient) GetFinishChannel(id string) *chan struct{} {
	tmpReq := l.requests.Get(id)
	if tmpReq == nil {
		return nil
	}
	tmpVal := tmpReq.(request)
	return &tmpVal.finishedChannel
}

func (l *LLMClient) GetCompletedFor(id string) time.Duration {
	tmpReq := l.requests.Get(id)
	if tmpReq == nil {
		return 0
	}
	tmpVal := tmpReq.(request)
	return tmpVal.duration
}

func (l *LLMClient) GetDurationStatus(id string) time.Duration {
	tmpReq := l.requests.Get(id)
	if tmpReq == nil {
		return 0
	}
	tmpVal := tmpReq.(request)
	return time.Now().Sub(tmpVal.startTime)
}
