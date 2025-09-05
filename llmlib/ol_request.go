package llmlib

import (
	"context"
	"log"
	"time"

	"github.com/ollama/ollama/api"
)

func waitForAnswerOllama(
	chatRequest *api.ChatRequest,
	responseFunction func(resp api.ChatResponse) error,
	ctx context.Context,
	waitCh chan struct{},
	l *LLMClient,
	id string,
) {
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
}
