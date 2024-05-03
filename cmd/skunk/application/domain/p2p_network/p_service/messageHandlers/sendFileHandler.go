package messageHandlers

import (
	"encoding/json"
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/chat"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
)

// A Peer sends a file to a chat
type SendFileHandler struct {
	userChatLogic chat.ChatLogic
}

func (s *SendFileHandler) HandleMessage(message network.Message) error {

	// Structure of the message:
	/*
		{
			"fileContent": []byte("asdfasdfasdf"),
		}
	*/

	var content struct {
		FileContent []byte `json:"fileContent"`
	}

	err := json.Unmarshal([]byte(message.Content), &content)
	if err != nil {
		fmt.Println("Error unmarshalling message content")
		return err
	}

	s.userChatLogic.ReceiveFile(message.FromUser, message.ChatID, string(content.FileContent))
	//TODO that's not the right way to do it
	// The file should be stored in a file system
	// Also "ReceiveFile" should receive the file path instead of the file content

	return nil
}
