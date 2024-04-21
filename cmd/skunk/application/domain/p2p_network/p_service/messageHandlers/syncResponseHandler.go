package messageHandlers

import (
	"encoding/json"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_model"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
)

type SyncResponseHandler struct {
}

func (s *SyncResponseHandler) HandleMessage(message network.Message) error {
	/*

		syncResponse := network.Message{
			Id:        uuid.New().String(),
			Timestamp: time.Now().UnixNano(),
			Content:   string(externalMessagesBytes),
			FromUser:  chatMessageRepo.GetUsername(),
			ChatID:    message.ChatID,
			Operation: network.SYNC_RESPONSE,
		}

		{4b195e69-636c-4c7a-b0a5-481460637052
		1713729547132972295
			[{"Id":"internalMessage123!",
			"Timestamp":1633029448,
			"Content":"LOOOOOOOOOOOOOOOOOOOOOOOL",
			"FromUser":"as23d",
			"ChatID":"asdf",
			"Operation":1},
			{"Id":"internalMessage2",
			"Timestamp":1633029448,
			"Content":"WOOW",
			"FromUser":"as23d",
			"ChatID":"asdf",
			"Operation":1}]
		(Username)
		asdf
		2}
	*/

	chatRepo := p_model.GetNetworkChatsInstance()
	chatMessageRepo := chatRepo.GetChat(message.ChatID)

	var receivedMessages []network.Message
	err := json.Unmarshal([]byte(message.Content), &receivedMessages)
	if err != nil {
		return err
	}

	for _, message := range receivedMessages {
		chatMessageRepo.AddMessage(message)
	}

	return nil
}
