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

	tag "github.com/dhowden/tag"
	cl "github.com/serg-2/libs-go/commonlib"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

const MAXIMUM_TRACK_IN_PLAYLIST_SIZE = 7

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

// SendPictureFile - send picture
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

// SendVideoFile - Send video
func SendVideoFile(chatId int64, filename string, videoCaption string, bot *tgbotapi.BotAPI) {
	videoBytes, err := ioutil.ReadFile(filename)
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
	fileBytes, err := ioutil.ReadFile(filename)
	cl.ChkNonFatal(err)

	tgFileBytes := tgbotapi.FileBytes{
		Name:  "",
		Bytes: fileBytes,
	}
	return tgFileBytes
}

// Read MP3 tag by filename
func readTag(filename string) (string, string) {
	f, err := os.Open(filename)
	cl.ChkNonFatal(err)
	defer f.Close()

	m, err := tag.ReadFrom(f)
	cl.ChkNonFatal(err)

	// fmt.Println("Format", m.Format())
	//fmt.Printf("RAW: %v\n", m.Raw()["TPE1"])

	// DETECT ENCODING
	// import "github.com/mikkyang/id3-go"
	// mp3File, _ := id3.Open(filename)
	// defer mp3File.Close()
	// fmt.Println("Encoding:", mp3File.Frame("TPE1").(*v2.TextFrame).Encoding())

	// Using Artist and Title field
	if m.Format() == tag.ID3v1 {
		return cp1251ToUtf8(m.Artist()), cp1251ToUtf8(m.Title())
	} else if m.Format() == tag.ID3v2_3 {
		return v23convert(m.Artist()), v23convert(m.Title())
	} else {
		return m.Artist(), m.Title()
	}
}

// MP3 id3v2_3 convert
func v23convert(stringFrom string) string {
	// Transform from LATIN1. Because library TAG already decoded frame as LATIN1, what was wrong
	latin1, _, err := transform.String(charmap.ISO8859_1.NewEncoder(), stringFrom)
	// Check error of encode
	if err != nil {
		latin1 = stringFrom
	}

	// Decode from windows1251
	decoder := charmap.Windows1251.NewDecoder()
	res, _ := decoder.String(latin1)

	// If error res will be ""
	if len(res) == 0 {
		return latin1
	}
	return res
}

// Converting charset cp1251 to UTF8
func cp1251ToUtf8(stringFrom string) string {
	// MAIN
	decoder := charmap.Windows1251.NewDecoder()
	res, _ := decoder.String(stringFrom)
	return res
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
				audio.Performer, audio.Title = readTag(filename)

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
			audio.Performer, audio.Title = readTag(filename)

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
