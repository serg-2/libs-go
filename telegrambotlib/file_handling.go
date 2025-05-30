package telegrambotlib

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func DownloadFile(fileId string, fileName string, directoryName string, bot *tgbotapi.BotAPI) error {
	fileConfig := tgbotapi.FileConfig{FileID: fileId}

	file, err := bot.GetFile(fileConfig)
	if err != nil {
		return err
	}

	url := file.Link(bot.Token)

	err = lowLevelDownloadFile(url, fileName, directoryName)
	if err != nil {
		return err
	}
	return nil
}

func lowLevelDownloadFile(url string, fileName string, directoryName string) error {
	// Create the downloads directory if it doesn't exist
	if err := os.MkdirAll(directoryName, 0755); err != nil {
		return err
	}

	// Create the file
	filePath := filepath.Join(directoryName, fileName)
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
