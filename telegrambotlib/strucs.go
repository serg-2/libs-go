package telegrambotlib

// MAX_INPUT_MESSAGE_SIZE_FOR_COMMAND - max allowed user input with command (In bytes)
const MAX_INPUT_MESSAGE_SIZE_FOR_COMMAND int = 400

// MAX_MESSAGE_SIZE - MAX MESSAGE SIZE of telegram
const MAX_MESSAGE_SIZE int = 4000

type RangeStruct struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

// CommandStruct - Structure of commands for load
type CommandStruct struct {
	Name    string      `json:"name"`
	Level   string      `json:"level"`
	Range   RangeStruct `json:"range"`
	Help    string      `json:"help"`
	Aliases []string    `json:"aliases"`
}

type BotConfig struct {
	Token       string  `json:"Token"`
	Host        string  `json:"Host"`
	Port        int16   `json:"Port"`
	ListenHost  string  `json:"ListenHost"`
	Certificate string  `json:"Certificate"`
	Key         string  `json:"Key"`
	Superadmins []int64 `json:"superadmins"`
}
