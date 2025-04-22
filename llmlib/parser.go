package llmlib

import (
	"log"

	dsr "github.com/go-deepseek/deepseek/request"
	"github.com/go-deepseek/deepseek/response"
	js "github.com/serg-2/libs-go/jsonlib"
)

func parseResult(
	currentRequest *request,
	resp *response.Choice,
	previousMessages []*dsr.Message,
) {
	switch resp.FinishReason {
	case "stop":
		currentRequest.result = resp.Message.Content
		currentRequest.resultCalls = nil
		currentRequest.history = DStoSystem(previousMessages)
	case "tool_calls":
		if resp.Message.ToolCalls == nil {
			log.Printf("Received empty tool calls:\n%s\n", js.JsonAsString(resp))
		} else {
			currentRequest.result = "Answer is tool request"
			currentRequest.resultCalls = reparseToolCalls(resp.Message.ToolCalls)
			currentRequest.history = DStoSystem(previousMessages)
		}
	default:
		log.Printf("Received unparsed choise:\n%s\n", js.JsonAsString(resp))
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
		result = append(result,
			&dsr.Message{
				Role:       mess.Role,
				Content:    mess.Content,
				Name:       mess.Name,
				ToolCallId: mess.ToolCallId,
			})
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
