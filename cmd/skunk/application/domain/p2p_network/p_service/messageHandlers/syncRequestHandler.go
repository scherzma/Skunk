package messageHandlers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_model"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"time"
)

type SyncRequestHandler struct {
}

func (s *SyncRequestHandler) HandleMessage(message network.Message) error {

	fmt.Println("SyncRequestHandler")

	chatRepo := p_model.GetNetworkChatsInstance()
	chatMessageRepo := chatRepo.GetChat(message.ChatID)

	// Parse the content of the message
	/*
		{
		  "existingMessageIds": [
			"<message id 1>",
			"<message id 2>",
			...
		  ]
		}
	*/

	var content struct {
		ExistingMessageIDs []string `json:"existingMessageIds"`
	}
	err := json.Unmarshal([]byte(message.Content), &content)
	if err != nil {
		fmt.Println("Error unmarshalling message content")
		return err
	}

	// Find difference between "message" already known messages and own messages that the other peer does not know
	missingExternalMessages := chatMessageRepo.GetMissingExternalMessages(content.ExistingMessageIDs)
	missingInternalMessages := chatMessageRepo.GetMissingInternalMessageIDs(content.ExistingMessageIDs)

	// Convert missingExternalMessages to a JSON string
	externalMessagesBytes, err := json.Marshal(missingExternalMessages)
	if err != nil {
		fmt.Println("Error marshalling missing external messages")
		return err
	}

	// Convert missingInternalMessages to a JSON string
	internalMessagesBytes, err := json.Marshal(missingInternalMessages)
	if err != nil {
		fmt.Println("Error marshalling missing internal messages")
		return err
	}

	// Send the sync response to the other peer
	syncResponse := network.Message{
		Id:        uuid.New().String(),
		Timestamp: time.Now().UnixNano(),
		Content:   string(externalMessagesBytes),
		FromUser:  chatMessageRepo.GetUsername(),
		ChatID:    message.ChatID,
		Operation: network.SYNC_RESPONSE,
	}
	fmt.Println("syncResponse: ", syncResponse)
	fmt.Println("syncResponse.Content: ", syncResponse.Content)

	// Send sync request to other peer to get the difference between the messages that the other peer knows this peer does not know
	syncRequest := network.Message{
		Id:        uuid.New().String(),
		Timestamp: time.Now().UnixNano(),
		Content:   string(internalMessagesBytes),
		FromUser:  chatMessageRepo.GetUsername(),
		ChatID:    message.ChatID,
		Operation: network.SYNC_REQUEST,
	}
	fmt.Println("syncRequest.Content: ", syncRequest.Content)
	fmt.Println("syncRequest: ", syncRequest)

	peer := GetPeerInstance()
	peer.SendMessageToNetworkPeer("addressResponse", syncResponse)
	peer.SendMessageToNetworkPeer("addressRequest", syncRequest)

	return nil
}
