package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jesseduffield/horcrux/pkg/commands"
)

type Secrets struct {
	TOKEN      string `json:"TOKEN"`
	CHANNEL_ID string `json:"CHANNEL_ID"`
}

func main() {
	var userChoice string

	fmt.Print("Enter 1 for split + upload, 2 for download + bind: ")
	fmt.Scanln(&userChoice)

	if userChoice == "1" {
		splitAndUpload()
	} else if userChoice == "2" {
		var userFile string
		fmt.Print("Enter file name without extension for download: ")
		fmt.Scanln(&userFile)

		downloadAndBind(userFile)
	}
}

func downloadAndBind(userInput string) {
	var secrets Secrets
	json_file, err := os.Open("secrets.json")
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return
	}

	decoder := json.NewDecoder(json_file)
	err = decoder.Decode(&secrets)
	if err != nil {
		fmt.Println("Error decoding JSON data:", err)
		return
	}

	dg, err := discordgo.New("Bot " + secrets.TOKEN)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return
	}

	messages, err := dg.ChannelMessages(secrets.CHANNEL_ID, 100, "", "", "")
	if err != nil {
		fmt.Println("Error fetching channel messages:", err)
		return
	}

	for _, message := range messages {
		for _, attachment := range message.Attachments {
			if strings.Contains(attachment.Filename, userInput) {
				err := downloadAttachment(attachment.URL, attachment.Filename)
				if err != nil {
					fmt.Println("Error downloading horcrux:", err)
				} else {
					fmt.Println("Attachment horcrux:", attachment.Filename)
				}
			}
		}
	}

	paths, err := commands.GetHorcruxPathsInDir("./")
	overwrite := true
	if err != nil {
		fmt.Println("Error getting horcrux paths: ", err)
	}

	if err := commands.Bind(paths, "", overwrite); err != nil {
		fmt.Println("Error binding: ", err)
	}

	dg.Close()
}

func downloadAttachment(url, filename string) error {
	// Open a file for writing the attachment.
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get the attachment content.
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the attachment content to the file.
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func splitAndUpload() {
	var file_input string

	fmt.Print("Enter file name: ")
	fmt.Scanln(&file_input)
	secretFile, err := os.Open(file_input)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	fileInfo, err := secretFile.Stat()
	if err != nil {
		fmt.Println("Error retrieving file info: ", err)
	}

	fileSize := float64(fileInfo.Size()) / (1024.0 * 1024.0)
	splitAmount := 0
	maxSize := 8
	for {
		splitAmount++
		tempSize := fileSize / float64(splitAmount)
		if tempSize <= float64(maxSize) {
			break
		}
	}

	var horcruxNames []string
	commands.Split(file_input, "./", splitAmount, splitAmount, &horcruxNames)

	var secrets Secrets
	json_file, err := os.Open("secrets.json")
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return
	}

	decoder := json.NewDecoder(json_file)
	err = decoder.Decode(&secrets)
	if err != nil {
		fmt.Println("Error decoding JSON data:", err)
		return
	}

	dg, err := discordgo.New("Bot " + secrets.TOKEN)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return
	}

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection:", err)
		return
	}

	for _, value := range horcruxNames {
		file, err := os.Open(value)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}

		_, err = dg.ChannelFileSend(secrets.CHANNEL_ID, file.Name(), file)
		if err != nil {
			fmt.Println("Error sending file:", err)
			return
		} else {
			fmt.Println("Uploading", value)
		}
	}

	dg.Close()
}
