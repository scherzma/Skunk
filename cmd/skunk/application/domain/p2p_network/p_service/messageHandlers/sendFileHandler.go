package messageHandlers

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/chat"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/store"
	"io"
	"os"
	"path/filepath"
)

// A Peer sends a file to a chat
type sendFileHandler struct {
	userChatLogic         chat.ChatLogic
	chatInvitationStorage store.ChatInvitationStoragePort
}

func NewSendFileHandler(userChatLogic chat.ChatLogic, chatInvitationStorage store.ChatInvitationStoragePort) *sendFileHandler {
	return &sendFileHandler{
		userChatLogic:         userChatLogic,
		chatInvitationStorage: chatInvitationStorage,
	}
}

func (s *sendFileHandler) HandleMessage(message network.Message) error {

	// Structure of the message (fileContent is base64 encoded):
	/*
		{
			"fileName": "file_name",
			"fileExtension": "file_extension",
			"fileContent": "YmFzZTY0c3RyaW5n",
		}
	*/

	var content struct {
		FileName      string `json:"fileName"`
		FileExtension string `json:"fileExtension"`
		FileContent   string `json:"fileContent"`
	}

	err := json.Unmarshal([]byte(message.Content), &content)
	if err != nil {
		fmt.Println("Error unmarshalling message content")
		return err
	}

	// Decode the base64-encoded file content
	fileData, err := base64.StdEncoding.DecodeString(content.FileContent)
	if err != nil {
		fmt.Println("Error decoding file content")
		return err
	}

	// Create a directory for storing the files if it doesn't exist
	fileDir := "./stored_files"
	if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
		fmt.Println("Error creating file directory")
		return err
	}

	// Create the file path
	filePath := filepath.Join(fileDir, fmt.Sprintf("%s.%s", content.FileName, content.FileExtension))

	// Check if the file already exists
	if _, err := os.Stat(filePath); err == nil {
		// File exists, compare hashes
		existingHash, err := computeFileHash(filePath)
		if err != nil {
			fmt.Println("Error computing hash of existing file")
			return err
		}

		newFileHash := computeHash(fileData)
		if existingHash == newFileHash {
			// Hashes are the same, use the existing file
			s.userChatLogic.ReceiveFile(message.SenderID, message.ChatID, filePath)
			return nil
		} else {
			// Hashes are different, modify the file name
			filePath = getUniqueFilePath(fileDir, content.FileName, content.FileExtension)
		}
	}

	// Write the file data to the file
	if err := os.WriteFile(filePath, fileData, os.ModePerm); err != nil {
		fmt.Println("Error writing file to filesystem")
		return err
	}

	// Notify the chat logic of the received file with the file path
	s.userChatLogic.ReceiveFile(message.SenderID, message.ChatID, filePath)

	return nil
}

// computeFileHash computes the SHA-256 hash of a file
func computeFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// computeHash computes the SHA-256 hash of the given data
func computeHash(data []byte) string {
	hash := sha256.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

// getUniqueFilePath generates a unique file path by appending a counter to the file name
func getUniqueFilePath(dir, baseName, ext string) string {
	counter := 1
	for {
		newFileName := fmt.Sprintf("%s_%d.%s", baseName, counter, ext)
		newFilePath := filepath.Join(dir, newFileName)
		if _, err := os.Stat(newFilePath); os.IsNotExist(err) {
			return newFilePath
		}
		counter++
	}
}
