package telegrambotlib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	cl "github.com/serg-2/libs-go/commonlib"
)

// LoadCommands - Load commands from JSON file
func LoadCommands(file string) []CommandStruct {
	var config []CommandStruct
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		log.Fatalln("Can't open commands file " + file)
		//log.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		log.Fatalln("Bad JSON in " + file)
		//log.Println(err.Error())
	}
	return config
}

// CheckAndGetCommand - Checks command and get it.
func CheckAndGetCommand(message string, commandsP *[]CommandStruct) (string, []string) {
	var command string
	var arguments []string = make([]string, 0)

	// Check too long message
	if len(message) > MAX_INPUT_MESSAGE_SIZE_FOR_COMMAND {
		return "", nil
	}

	// too many spaces
	if len(strings.Split(message, " ")) > 12 {
		return "", nil
	}

	// Split message
	splittedMessage := strings.Split(message, " ")

	// Extract arguments
	if len(splittedMessage) > 1 {
		arguments = splittedMessage[1:]
	}

	command = strings.ToLower(splittedMessage[0])

	if len(command) > 1 {
		if string(command[0]) == "/" {
			command = command[1:]
		}
	}

	// Convert alias to command
	command = getCommandFromAliases(command, commandsP)

	// find arguments and existence in map
	argsRange, found := getCommandRange(command, commandsP)

	// Check command exists
	if !found {
		return "", nil
	}

	// Check number of arguments is ok
	if len(arguments) < argsRange.Min || len(arguments) > argsRange.Max {
		return "", nil
	}
	return command, arguments
}

func getCommandRange(commandName string, commandsP *[]CommandStruct) (RangeStruct, bool) {
	for _, command := range *commandsP {
		if command.Name == commandName {
			return command.Range, true
		}
	}
	return RangeStruct{}, false
}

func getCommandFromAliases(s string, commandsP *[]CommandStruct) string {
	for _, command := range *commandsP {
		// Check Main command
		if command.Name == s {
			return command.Name
		}
		// Check Aliases
		for _, alias := range command.Aliases {
			if alias == s {
				return command.Name
			}
		}
	}
	return ""
}

func SendLowLevel(targetChatId int64, message string, bot *tgbotapi.BotAPI) {
	//	log.Printf("Full Message Size: %d\n", len(message))
	for len(message) > MAX_MESSAGE_SIZE {
		// Warning! Beware of splitting runes!
		var corrector int = 0
		for !utf8.RuneStart(message[MAX_MESSAGE_SIZE-corrector]) {
			corrector += 1
		}

		msg := tgbotapi.NewMessage(targetChatId, message[:MAX_MESSAGE_SIZE-corrector])
		bot.Send(msg)
		message = message[MAX_MESSAGE_SIZE-corrector:]
		//		log.Printf("Rest Message Size: %d\n", len(message))
	}
	msg := tgbotapi.NewMessage(targetChatId, message)
	bot.Send(msg)
}

func SendToGroup(message string, userList []int64, bot *tgbotapi.BotAPI) {
	for _, user := range userList {
		SendLowLevel(user, message, bot)
	}
}

func SendPictureUrl(chatId int64, link string, pictureName string, bot *tgbotapi.BotAPI) error {
	response, e := http.Get(link)
	if e != nil {
		log.Println("Error getting link during sendPictureUrl", e, e.Error())
		return errors.New("error getting link during sendPictureUrl")
	}
	defer response.Body.Close()

	var dst bytes.Buffer

	_, err := io.Copy(&dst, response.Body)
	if err != nil {
		log.Println("Error copying body to buffer in sendPictureUrl", e, e.Error())
		return errors.New("error copying body to buffer in sendPictureUrl")
	}

	photoFileBytes := tgbotapi.FileBytes{
		Name:  pictureName,
		Bytes: dst.Bytes(),
	}

	photoUpload := tgbotapi.NewPhoto(chatId, photoFileBytes)
	_, err = bot.Send(photoUpload)
	if err != nil {
		return err
	}
	return nil
}

func WriteLog(message tgbotapi.Update) {
	receivedAt := message.Message.Time().Format(time.RFC3339)
	//entities := message.Message.Entities
	//log.Printf("Entities: %v\n", entities)
	//log.Printf("Entity type: %v\n", a)
	logMessage := receivedAt + " " + fmt.Sprintf("%d", message.Message.Chat.ID) + " " + message.Message.Text + "\n"

	// Always defined ID + FirstName
	path := fmt.Sprintf("%d", message.Message.From.ID) + "_"

	if message.Message.From.FirstName != "" {
		path += message.Message.From.FirstName
	}
	if message.Message.From.LastName != "" {
		path += message.Message.From.LastName
	}
	if message.Message.From.UserName != "" {
		path += message.Message.From.UserName
	}

	if _, err := os.Stat("users"); os.IsNotExist(err) {
		os.Mkdir("users", os.ModePerm)
	}

	file, err := os.OpenFile("users/"+path+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	cl.ChkFatal(err)
	defer file.Close()

	_, err = file.Write([]byte(logMessage))
	cl.ChkFatal(err)
}

func SendPictureFile(chatId int64, filename string, pictureName string, bot *tgbotapi.BotAPI) {
	photoBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		cl.ChkNonFatal(err)
	}
	photoFileBytes := tgbotapi.FileBytes{
		Name:  pictureName,
		Bytes: photoBytes,
	}

	photoUpload := tgbotapi.NewPhoto(chatId, photoFileBytes)
	bot.Send(photoUpload)
}

func SendBytesFile(chatId int64, buf bytes.Buffer, filename string, bot *tgbotapi.BotAPI) {
	// Creating filename and payload
	fileInterface := tgbotapi.FileBytes{
		Name:  filename,
		Bytes: buf.Bytes(),
	}

	// Sending answer
	docUpload := tgbotapi.NewDocument(chatId, fileInterface)
	bot.Send(docUpload)
}
