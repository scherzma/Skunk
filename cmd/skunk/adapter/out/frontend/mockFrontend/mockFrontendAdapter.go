package mockFrontend

import (
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/frontend"
	"sync"
)

var (
	mockFrontendInstance *MockFrontend
	once                 sync.Once
)

type MockFrontend struct {
	observers []frontend.FrontendObserver
}

func GetMockFrontendInstance() *MockFrontend {
	once.Do(func() {
		mockFrontendInstance = &MockFrontend{
			observers: make([]frontend.FrontendObserver, 0),
		}
	})
	return mockFrontendInstance
}

func (mf *MockFrontend) SubscribeToFrontend(observer frontend.FrontendObserver) error {
	mf.observers = append(mf.observers, observer)
	return nil
}

func (mf *MockFrontend) UnsubscribeFromFrontend(observer frontend.FrontendObserver) error {
	for i, obs := range mf.observers {
		if obs == observer {
			mf.observers = append(mf.observers[:i], mf.observers[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("observer not found")
}

func (mf *MockFrontend) SendToFrontend(message frontend.FrontendMessage) error {
	for _, observer := range mf.observers {
		err := observer.Notify(message)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mf *MockFrontend) SendMockFrontendMessageToSubscribers(message frontend.FrontendMessage) error {
	return mf.SendToFrontend(message)
}
