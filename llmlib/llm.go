package llmlib

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/cohesion-org/deepseek-go"
	"github.com/google/uuid"
)

var availableRoles []string = []string{
	"system",
	"user",
	"assistant",
	"tool",
}

func (l *LLMClient) AddRequest(
	senderId int64,
	sourceChatId int64,
	question string,
	previosMessages []SystemMessages,
	tools []deepseek.Tool,
	passedFunction PassedFunction,
) string {
	// Validate system messages
	if !validateSystemMessages(previosMessages) {
		log.Println("Can't validate request.")
		return ""
	}

	//log.Printf("Additional context questions: %s\n", js.JsonAsString(previosMessages))

	// Validate question
	if !validateQuestion(question) {
		log.Println("Can't validate question.")
		return ""
	}

	id := uuid.New().String()
	ctx := context.Background()
	waitCh := make(chan struct{})

	if l.clientOllama != nil {
		// OLLAMA
		// NO TOOLS For Now
		go waitForAnswerOllama(
			getRequestOllama(l, question, previosMessages),
			getResonseFunctionOllama(l, id),
			ctx,
			waitCh,
			l,
			id,
		)
	} else if l.clientDS != nil {
		go waitForAnswerDS(
			getRequestDS(l, question, previosMessages, tools),
			ctx,
			waitCh,
			l,
			id,
			passedFunction,
			sourceChatId,
			senderId,
		)
	} else {
		log.Println("Can't find clients.")
		return ""
	}

	var newReq request = request{
		false,
		time.Now(),
		0,
		waitCh,
		0,
		"answer is not ready",
		[]SystemToolCalls{},
		[]SystemMessages{},
	}
	l.requests.Add(id, newReq)

	return id
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

func (l *LLMClient) GetHistory(id string) []SystemMessages {
	tmpReq := l.requests.Get(id)
	if tmpReq == nil {
		return nil
	}
	tmpVal := tmpReq.(request)
	return tmpVal.history
}

func (l *LLMClient) GetRetries(id string) int {
	tmpReq := l.requests.Get(id)
	if tmpReq == nil {
		return 0
	}
	tmpVal := tmpReq.(request)
	return tmpVal.numberRetries
}

func (l *LLMClient) GetCallsFromAnswer(id string) []SystemToolCalls {
	tmpReq := l.requests.Get(id)
	if tmpReq == nil {
		return []SystemToolCalls{}
	}
	tmpVal := tmpReq.(request)
	return tmpVal.resultCalls
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
