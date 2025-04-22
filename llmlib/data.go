package llmlib

import (
	"time"

	"github.com/go-deepseek/deepseek"
	"github.com/ollama/ollama/api"
	cl "github.com/serg-2/libs-go/commonlib"
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
	startTime       time.Time
	duration        time.Duration
	finishedChannel chan struct{}

	result      string
	resultCalls []SystemToolCalls

	history []SystemMessages
}

type SystemMessages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	// Could be tool_call.name
	Name       string `json:"name"`
	ToolCallId string `json:"tool_call_id"`
}

type SystemToolCalls struct {
	Id       string             `json:"id"`
	Type     string             `json:"type"`
	Function SystemToolFunction `json:"function"`
}

type SystemToolFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type PassedFunction func(call SystemToolCalls) string
