package llmlib

import (
	"context"
	"log"
	"time"

	"github.com/cohesion-org/deepseek-go"
	js "github.com/serg-2/libs-go/jsonlib"
)

func waitForAnswerDS(
	chatRequest *deepseek.ChatCompletionRequest,
	ctx context.Context,
	waitCh chan struct{},
	l *LLMClient,
	id string,
	passedFunction PassedFunction,
	sourceChatId int64,
	senderId int64,
) {
	log.Printf("DEBUG: Full Request: %s\n", js.JsonAsString(chatRequest))
	firstResponse, err := l.clientDS.CreateChatCompletion(ctx, chatRequest)
	log.Printf("DEBUG: Full first Response: %s\n", js.JsonAsString(firstResponse))

	currentRequest := l.requests.Get(id).(request)

	if err != nil {
		log.Println("Error in Chat handling DS. Try 1")
		log.Println(err)

		// RETRY
		currentRequest.numberRetries += 1
		time.Sleep(1 * time.Second)
		firstResponse, err = l.clientDS.CreateChatCompletion(ctx, chatRequest)
		log.Printf("DEBUG: Full first Response. TRY 2: %s\n", js.JsonAsString(firstResponse))
		// END RETRY
		if err != nil {
			currentRequest.result = "Error in Chat handling. Last Error: " + err.Error()
		}
	}

	if err == nil {
		parseResult(&currentRequest, firstResponse.Choices[0], chatRequest.Messages)
		// Check need to make functions
		if currentRequest.resultCalls != nil {
			// Call functions
			var toolsAnswers []deepseek.ChatCompletionMessage

			toolsAnswers = append(
				toolsAnswers,
				SystemToDS(currentRequest.history)...,
			)

			for _, call := range currentRequest.resultCalls {
				respString, ok := passedFunction(call, sourceChatId, senderId)

				// Error.
				if ok == -1 {
					log.Printf("Skipping tool request as function response: %s\n", respString)
					currentRequest.result = "DONE WITH ERROR: " + respString
					currentRequest.duration = time.Now().Sub(currentRequest.startTime)
					currentRequest.finished = true
					l.requests.Add(id, currentRequest)
					close(waitCh)
					return
				}
				// Special case empty ok answer
				if ok == 1 {
					currentRequest.result = "DONE with empty answer"
					currentRequest.duration = time.Now().Sub(currentRequest.startTime)
					currentRequest.finished = true
					l.requests.Add(id, currentRequest)
					close(waitCh)
					return
				}

				// strict ok answer
				if ok == 2 {
					currentRequest.result = respString
					currentRequest.duration = time.Now().Sub(currentRequest.startTime)
					currentRequest.finished = true
					l.requests.Add(id, currentRequest)
					close(waitCh)
					return
				}

				// Ask ai for more...

				// Generate dsr message
				tmpMessage := deepseek.ChatCompletionMessage{
					Role:       "tool",
					Content:    respString,
					ToolCallID: call.Id,
				}
				toolsAnswers = append(toolsAnswers, tmpMessage)
			}
			// Summary of answer
			log.Printf("DEBUG: Request with tools answer: %s\n", js.JsonAsString(toolsAnswers))

			requestAfterTools := &deepseek.ChatCompletionRequest{
				Model:    l.model,
				Messages: toolsAnswers,
				Tools:    nil,
			}
			// Blocking response
			chatRespToTool, err2 := l.clientDS.CreateChatCompletion(ctx, requestAfterTools)
			if err2 != nil {
				log.Println("Error in Chat handling DS with TOOL")
				log.Println(err2)
			}

			parseResult(&currentRequest, chatRespToTool.Choices[0], requestAfterTools.Messages)

			log.Printf("DEBUG: Answer after tools answer:\n%s\n", js.JsonAsString(chatRespToTool))
		}
	}

	currentRequest.duration = time.Now().Sub(currentRequest.startTime)
	currentRequest.finished = true
	l.requests.Add(id, currentRequest)
	close(waitCh)
}
