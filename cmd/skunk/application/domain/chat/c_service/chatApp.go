package c_service

import (
	"github.com/scherzma/Skunk/cmd/skunk/application/port/frontend"
	"sync"
)

type ChatApp struct {
	frontends []frontend.Frontend
}

var (
	chatService *ChatApp
	once        sync.Once
)

func GetChatServiceInstance() *ChatApp {
	once.Do(func() {
		chatService = &ChatApp{
			frontends: []frontend.Frontend{},
		}
	})

	return chatService
}

// FrontendObserver: Gets notified when a frontend sends a message to the chat
func (c *ChatApp) Notify(message frontend.FrontendMessage) error {
	// TODO: Implement
	return nil
}

func (c *ChatApp) AddFrontend(frontend frontend.Frontend) {
	c.frontends = append(c.frontends, frontend)
}

func (c *ChatApp) RemoveFrontend(frontend frontend.Frontend) {
	for i, f := range c.frontends {
		if f == frontend {
			c.frontends = append(c.frontends[:i], c.frontends[i+1:]...)
			return
		}
	}
}

func (c *ChatApp) SendMessageToAllFrontends(message frontend.FrontendMessage) {
	for _, f := range c.frontends {
		f.SendToFrontend(message)
	}
}

func (c *ChatApp) ProcessMessageForUser(message frontend.FrontendMessage) error {
	// TODO: Implement
	return nil
}
