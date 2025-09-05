package llmlib

import (
	"time"

	"github.com/cohesion-org/deepseek-go"
	"github.com/ollama/ollama/api"
	cl "github.com/serg-2/libs-go/commonlib"
)

type LLMClient struct {
	clientOllama          *api.Client
	clientDS              *deepseek.Client
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
	numberRetries   int

	result      string
	resultCalls []SystemToolCalls

	history []SystemMessages
}

type SystemMessages struct {
	Role       string `json:"role"`
	Content    string `json:"content"`
	ToolCallId string `json:"tool_call_id"`

	internalDSToolCalls []deepseek.ToolCall
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

type PassedFunction func(call SystemToolCalls, sourceChatId int64, senderId int64) (string, int)
