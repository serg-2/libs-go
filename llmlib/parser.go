package llmlib

import (
	"log"

	"github.com/cohesion-org/deepseek-go"
	js "github.com/serg-2/libs-go/jsonlib"
)

func parseResult(
	currentRequest *request,
	choice0 deepseek.Choice,
	previousMessages []deepseek.ChatCompletionMessage,
) {
	switch choice0.FinishReason {
	case "stop":
		currentRequest.result = choice0.Message.Content
		currentRequest.resultCalls = nil
		currentRequest.history = DStoSystem(previousMessages)
	case "tool_calls":
		if choice0.Message.ToolCalls == nil {
			log.Printf("Received empty tool calls:\n%s\n", js.JsonAsString(choice0))
		} else {
			log.Println("Parsing tool calls...")
			currentRequest.result = "Answer is tool request"
			currentRequest.resultCalls = reparseToolCalls(choice0.Message.ToolCalls)
			currentRequest.history = DStoSystem(previousMessages)
			currentRequest.history = append(currentRequest.history,
				SystemMessages{
					Role:                "assistant",
					internalDSToolCalls: choice0.Message.ToolCalls,
				},
			)
		}
	default:
		log.Printf("Received unparsed choise:\n%s\n", js.JsonAsString(choice0))
	}
}

func DStoSystem(previousMessages []deepseek.ChatCompletionMessage) []SystemMessages {
	var result []SystemMessages
	for _, mess := range previousMessages {
		result = append(result,
			SystemMessages{
				Role:       mess.Role,
				Content:    mess.Content,
				ToolCallId: mess.ToolCallID,
			})
	}
	return result
}

func SystemToDS(previousMessages []SystemMessages) []deepseek.ChatCompletionMessage {
	var result []deepseek.ChatCompletionMessage
	for _, mess := range previousMessages {
		if mess.internalDSToolCalls != nil {
			result = append(result,
				deepseek.ChatCompletionMessage{
					Role:      "assistant",
					Content:   "Запрашиваю выполнение функций...",
					ToolCalls: mess.internalDSToolCalls,
				})
		} else {
			result = append(result,
				deepseek.ChatCompletionMessage{
					Role:       mess.Role,
					Content:    mess.Content,
				})
		}
	}
	return result
}

func reparseToolCalls(toolCall []deepseek.ToolCall) []SystemToolCalls {
	var result []SystemToolCalls
	for _, call := range toolCall {
		result = append(result,
			SystemToolCalls{
				Id:   call.ID,
				Type: call.Type,
				Function: SystemToolFunction{
					Name:      call.Function.Name,
					Arguments: call.Function.Arguments,
				},
			},
		)
	}
	return result
}
