package frontendAdapter

import (
	"fmt"
	"sync"

	"github.com/scherzma/Skunk/cmd/skunk/application/port/frontend"
)

// FrontendAdapter implements the frontend.Frontend interface
type FrontendAdapter struct {
	mu        sync.RWMutex
	observers map[frontend.FrontendObserver]struct{}
}

// NewFrontendAdapter creates a new FrontendAdapter
func NewFrontendAdapter() *FrontendAdapter {
	return &FrontendAdapter{
		observers: make(map[frontend.FrontendObserver]struct{}),
	}
}

// SubscribeToFrontend subscribes an observer to receive frontend messages
func (fa *FrontendAdapter) SubscribeToFrontend(observer frontend.FrontendObserver) error {
	fa.mu.Lock()
	defer fa.mu.Unlock()

	if _, exists := fa.observers[observer]; exists {
		return fmt.Errorf("observer already subscribed")
	}

	fa.observers[observer] = struct{}{}
	return nil
}

// UnsubscribeFromFrontend unsubscribes an observer from receiving frontend messages
func (fa *FrontendAdapter) UnsubscribeFromFrontend(observer frontend.FrontendObserver) error {
	fa.mu.Lock()
	defer fa.mu.Unlock()

	if _, exists := fa.observers[observer]; !exists {
		return fmt.Errorf("observer not found")
	}

	delete(fa.observers, observer)
	return nil
}

// SendToFrontend sends a message to all subscribed observers
func (fa *FrontendAdapter) SendToFrontend(message frontend.FrontendMessage) error {
	fa.mu.RLock()
	defer fa.mu.RUnlock()

	var err error
	for observer := range fa.observers {
		if notifyErr := observer.Notify(message); notifyErr != nil {
			err = notifyErr // capture the last error, but try to notify all observers
		}
	}
	return err
}
