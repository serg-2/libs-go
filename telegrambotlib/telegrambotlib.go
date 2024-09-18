package telegrambotlib

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	cl "github.com/serg-2/libs-go/commonlib"
)

const MAXIMUM_TRACK_IN_PLAYLIST_SIZE = 7

// CheckAndGetCommand - Checks command and get it.
func CheckAndGetCommand(message string, commandsP *[]CommandStruct) (string, []string) {
	var command string
	var arguments []string = make([]string, 0)

	// Check too long message
	if len(message) > MAX_INPUT_MESSAGE_SIZE_FOR_COMMAND {
		return "", nil
	}

	// too many spaces
	if len(strings.Split(message, " ")) > 20 {
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
	// Check command exists
	if command == "" {
		log.Printf("Command not found ***%s***\n", command)
		return "", nil
	}

	// find arguments and existence in map
	argsRange, found := getCommandRange(command, commandsP)

	// Check command exists
	if !found {
		log.Printf("Command range not found ***%s***\n", command)
		return "", nil
	}

	// Check number of arguments is ok
	if len(arguments) < argsRange.Min || len(arguments) > argsRange.Max {
		log.Printf("Number of arguments(%d) for command ***%s*** is not ok.\n", len(arguments), command)
		return "", nil
	}

	// Log command to console
	log.Printf("Parsing command: %s\n", command)

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

func GetCommandHelpFromAliases(s string, commandsP *[]CommandStruct) string {
	for _, command := range *commandsP {
		// Check Main command
		if command.Name == s {
			return command.Help
		}
		// Check Aliases
		for _, alias := range command.Aliases {
			if alias == s {
				return command.Help
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

func SendLowLevelDelayed(targetChatId int64, message string, bot *tgbotapi.BotAPI, delay time.Duration) {
	go func() {
		time.Sleep(delay)
		SendLowLevel(targetChatId, message, bot)
	}()
}

func SendToGroup(message string, userList []int64, bot *tgbotapi.BotAPI) {
	for _, user := range userList {
		SendLowLevel(user, message, bot)
	}
}

func SendPictureUrl(chatId int64, link string, pictureName string, caption string, bot *tgbotapi.BotAPI) error {
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

	// Add caption
	if caption != "" {
		photoUpload.Caption = caption
	}

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

// SendPictureFile - send picture
func SendPictureFile(chatId int64, filename string, pictureName string, bot *tgbotapi.BotAPI) {
	photoBytes, err := os.ReadFile(filename)
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

// SendVideoFile - Send video
func SendVideoFile(chatId int64, filename string, videoCaption string, bot *tgbotapi.BotAPI) {
	videoBytes, err := os.ReadFile(filename)
	if err != nil {
		cl.ChkNonFatal(err)
	}
	videoFileBytes := tgbotapi.FileBytes{
		Name:  videoCaption,
		Bytes: videoBytes,
	}

	videoUpload := tgbotapi.NewVideo(chatId, videoFileBytes)
	bot.Send(videoUpload)
}

func readBytes(filename string) tgbotapi.FileBytes {
	fileBytes, err := os.ReadFile(filename)
	cl.ChkNonFatal(err)

	tgFileBytes := tgbotapi.FileBytes{
		Name:  "",
		Bytes: fileBytes,
	}
	return tgFileBytes
}

// SendMediaFiles - Send media files in group
func SendMediaFiles(bot *tgbotapi.BotAPI, chatId int64, mediaType string, filenames []string) {
	var result []interface{}

	if mediaType == "audio" {
		// Send by MAXIMUM_TRACK_IN_PLAYLIST_SIZE
		for len(filenames) > MAXIMUM_TRACK_IN_PLAYLIST_SIZE {
			tmp_filenames := filenames[:MAXIMUM_TRACK_IN_PLAYLIST_SIZE]
			filenames = filenames[MAXIMUM_TRACK_IN_PLAYLIST_SIZE:]
			for _, filename := range tmp_filenames {
				bytes := readBytes(filename)
				audio := tgbotapi.NewInputMediaAudio(bytes)

				// Change default performer and title(as needed cp1251 -> UTF8)
				audio.Performer, audio.Title = ReadTag(filename)

				result = append(result, audio)
			}
			mediaGroup := tgbotapi.NewMediaGroup(chatId, result)
			bot.Send(mediaGroup)
			result = nil
		}
		// Send rest
		for _, filename := range filenames {
			bytes := readBytes(filename)
			audio := tgbotapi.NewInputMediaAudio(bytes)

			// Change default performer and title(as needed cp1251 -> UTF8)
			audio.Performer, audio.Title = ReadTag(filename)

			result = append(result, audio)
		}
		mediaGroup := tgbotapi.NewMediaGroup(chatId, result)
		bot.Send(mediaGroup)
	}
}

// SendByteArrayFile - Send Byte array to file
func SendByteArrayFile(chatId int64, byteArray []byte, filename string, bot *tgbotapi.BotAPI) {
	// Creating filename and payload
	fileInterface := tgbotapi.FileBytes{
		Name:  filename,
		Bytes: byteArray,
	}

	// Sending answer
	docUpload := tgbotapi.NewDocument(chatId, fileInterface)
	bot.Send(docUpload)
}

// BotInitialize - initialize new bot
func BotInitialize(config BotConfig) (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel) {

	bot, err := tgbotapi.NewBotAPI(config.Token)
	cl.ChkFatal(err)

	//bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Old style init webHook
	// _, err = bot.SetWebhook(tgbotapi.NewWebhookWithCert("https://"+config.Host+":"+fmt.Sprintf("%d", config.Port)+"/"+bot.Token, config.Certificate))
	// New Style init webHook
	webHook, _ := tgbotapi.NewWebhookWithCert("https://"+config.Host+":"+fmt.Sprintf("%d", config.Port)+"/"+bot.Token, tgbotapi.FilePath(config.Certificate))
	_, err = bot.Request(webHook)
	if err != nil {
		switch err.(type) {
		case *tgbotapi.Error:
			// Try one more time
			log.Printf("Internal server error, trying one more time...\n")
			_, err = bot.Request(webHook)
		default:
			cl.ChkFatal(err)
		}
	}

	cl.ChkFatal(err)

	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServeTLS(config.ListenHost+":"+fmt.Sprintf("%d", config.Port), config.Certificate, config.Key, nil)

	return bot, updates
}
