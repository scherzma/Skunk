package frontendAdapter

import (
	"fmt"
	"sync"

	"github.com/scherzma/Skunk/cmd/skunk/application/port/frontend"
)

// FrontendAdapter implements the frontend.Frontend interface
type FrontendAdapter struct {
	mu       sync.RWMutex
	observer frontend.FrontendObserver
}

// NewFrontendAdapter creates a new FrontendAdapter
func NewFrontendAdapter() *FrontendAdapter {
	return &FrontendAdapter{}
}

// SubscribeToFrontend subscribes an observer to receive frontend messages
func (fa *FrontendAdapter) SubscribeToFrontend(observer frontend.FrontendObserver) error {
	fa.mu.Lock()
	defer fa.mu.Unlock()

	if fa.observer != nil {
		return fmt.Errorf("observer already subscribed")
	}

	fa.observer = observer
	return nil
}

// UnsubscribeFromFrontend unsubscribes an observer from receiving frontend messages
func (fa *FrontendAdapter) UnsubscribeFromFrontend(observer frontend.FrontendObserver) error {
	fa.mu.Lock()
	defer fa.mu.Unlock()

	if fa.observer == nil {
		return fmt.Errorf("no observer subscribed")
	}

	if fa.observer != observer {
		return fmt.Errorf("observer not found")
	}

	fa.observer = nil
	return nil
}

// SendToFrontend sends a message to the subscribed observer
func (fa *FrontendAdapter) SendToFrontend(message frontend.FrontendMessage) error {
	fa.mu.RLock()
	defer fa.mu.RUnlock()

	if fa.observer == nil {
		return fmt.Errorf("no observer subscribed")
	}

	return fa.observer.Notify(message)
}
