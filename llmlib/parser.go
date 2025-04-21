package llmlib

import (
	"log"

	"github.com/go-deepseek/deepseek/response"
	js "github.com/serg-2/libs-go/jsonlib"
)

func parseResult(currentRequest *request, resp *response.Choice) {
	switch resp.FinishReason {
	case "stop":
		currentRequest.result = resp.Message.Content
		currentRequest.resultCalls = []SystemToolCalls{}
	case "tool_calls":
		if resp.Message.ToolCalls == nil {
			log.Printf("Received empty tool calls:\n%s\n", js.JsonAsString(resp))
		} else {
			currentRequest.result = "Answer is tool request"
			currentRequest.resultCalls = reparseToolCalls(resp.Message.ToolCalls)
		}
	default:
		log.Printf("Received unparsed choise:\n%s\n", js.JsonAsString(resp))
	}
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
