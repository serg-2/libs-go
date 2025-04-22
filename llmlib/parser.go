package llmlib

import (
	"log"

	dsr "github.com/go-deepseek/deepseek/request"
	"github.com/go-deepseek/deepseek/response"
	js "github.com/serg-2/libs-go/jsonlib"
)

func parseResult(
	currentRequest *request,
	choice0 *response.Choice,
	previousMessages []*dsr.Message,
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

func DStoSystem(previousMessages []*dsr.Message) []SystemMessages {
	var result []SystemMessages
	for _, mess := range previousMessages {
		result = append(result,
			SystemMessages{
				Role:       mess.Role,
				Content:    mess.Content,
				Name:       mess.Name,
				ToolCallId: mess.ToolCallId,
			})
	}
	return result
}

func SystemToDS(previousMessages []SystemMessages) []*dsr.Message {
	var result []*dsr.Message
	for _, mess := range previousMessages {
		if mess.internalDSToolCalls != nil {
			for _, tc := range mess.internalDSToolCalls {
				result = append(result,
					&dsr.Message{
						Role:       "assistant",
						Content:    "Запрашиваю выполнение функции: " + tc.Function.Name,
						Name:       tc.Function.Name,
						ToolCallId: tc.Id,
					})
			}
		} else {
			result = append(result,
				&dsr.Message{
					Role:       mess.Role,
					Content:    mess.Content,
					Name:       mess.Name,
					ToolCallId: mess.ToolCallId,
				})
		}
	}
	return result
}

func reparseToolCalls(toolCall []*response.ToolCall) []SystemToolCalls {
	var result []SystemToolCalls
	for _, call := range toolCall {
		result = append(result,
			SystemToolCalls{
				Id:   call.Id,
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
