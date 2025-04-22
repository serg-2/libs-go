package llmlib

import (
	"context"
	"log"
	"time"

	dsr "github.com/go-deepseek/deepseek/request"
	js "github.com/serg-2/libs-go/jsonlib"
)

func waitForAnswerDS(
	chatRequest *dsr.ChatCompletionsRequest,
	ctx context.Context,
	waitCh chan struct{},
	l *LLMClient,
	id string,
	passedFunction PassedFunction,
) {
	log.Printf("FULL Request: %s\n", js.JsonAsString(chatRequest))
	firstResponse, err := l.clientDS.CallChatCompletionsChat(ctx, chatRequest)
	log.Printf("FULL FIRST Response: %s\n", js.JsonAsString(firstResponse))
	currentRequest := l.requests.Get(id).(request)
	if err != nil {
		log.Println("Error in Chat handling DS")
		log.Println(err)
		currentRequest.result = "Error in Chat handling: " + err.Error()
	} else {
		parseResult(&currentRequest, firstResponse.Choices[0], chatRequest.Messages)
		// Check need to make functions
		if currentRequest.resultCalls != nil {
			// Call functions
			var toolsAnswers []*dsr.Message

			toolsAnswers = append(
				toolsAnswers,
				SystemToDS(currentRequest.history)...,
			)

			for _, call := range currentRequest.resultCalls {
				var respString string = passedFunction(call)
				// Generate dsr message
				tmpMessage := dsr.Message{
					Role:       "tool",
					Content:    respString,
					Name:       call.Function.Name,
					ToolCallId: call.Id,
				}
				toolsAnswers = append(toolsAnswers, &tmpMessage)
			}
			// Summary of answer
			log.Printf("Request with tools answer: %s\n", js.JsonAsString(toolsAnswers))

			newReq := &dsr.ChatCompletionsRequest{
				Model:    l.model,
				Stream:   false,
				Messages: toolsAnswers,
				Tools:    nil,
			}
			// Blocking response
			chatRespToTool, err2 := l.clientDS.CallChatCompletionsChat(ctx, newReq)
			if err2 != nil {
				log.Println("Error in Chat handling DS with TOOL")
				log.Println(err2)
			}

			log.Printf("DS Answer after tools:\n%s\n", js.JsonAsString(chatRespToTool))
		}
	}
	currentRequest.duration = time.Now().Sub(currentRequest.startTime)
	currentRequest.finished = true
	l.requests.Add(id, currentRequest)
	close(waitCh)
}
