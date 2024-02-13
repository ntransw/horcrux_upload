package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/jesseduffield/horcrux/pkg/commands"
)

type Secrets struct {
	TOKEN      string `json:"TOKEN"`
	CHANNEL_ID string `json:"CHANNEL_ID"`
}

func main() {
	splitAndUpload()
}

func downloadAttachment(s *discordgo.Session, attachment *discordgo.MessageAttachment) {
	// Create a new file to save the attachment
	file, err := os.Create(attachment.Filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Download the attachment content
	resp, err := s.Client.Get(attachment.URL)
	if err != nil {
		fmt.Println("Error downloading attachment:", err)
		return
	}
	defer resp.Body.Close()

	// Copy the attachment content to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("Error copying attachment content:", err)
		return
	}

	fmt.Printf("Attachment '%s' downloaded successfully.\n", attachment.Filename)
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
	splitAmount := 1
	maxSize := 8
	for {
		splitAmount++
		fileSize /= float64(splitAmount)
		if fileSize <= float64(maxSize) {
			break
		}
	}

	var horcruxNames []string
	commands.Split(file_input, "./", splitAmount, splitAmount, &horcruxNames)

	json_file, err := os.Open("secrets.json")
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return
	}

	var secrets Secrets
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
		}
	}

	dg.Close()
}
