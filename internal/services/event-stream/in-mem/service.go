package inmemeventstream

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"

	eventstream "github.com/Pickausernaame/chat-service/internal/services/event-stream"
	"github.com/Pickausernaame/chat-service/internal/types"
)

const serviceName = "event-stream"

type Service struct {
	mtx       sync.RWMutex
	subs      map[types.UserID]map[chan eventstream.Event]struct{}
	closeChan chan struct{}
	lg        *zap.Logger
}

func New() *Service {
	return &Service{
		mtx:       sync.RWMutex{},
		subs:      map[types.UserID]map[chan eventstream.Event]struct{}{},
		closeChan: make(chan struct{}),
		lg:        zap.L().Named(serviceName),
	}
}

func (s *Service) Subscribe(ctx context.Context, userID types.UserID) (<-chan eventstream.Event, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.subs[userID] == nil {
		s.subs[userID] = make(map[chan eventstream.Event]struct{})
	}

	ch := make(chan eventstream.Event)

	s.subs[userID][ch] = struct{}{}

	go func() {
		<-ctx.Done()
		s.mtx.Lock()
		defer s.mtx.Unlock()
		if _, ok := s.subs[userID]; ok {
			delete(s.subs[userID], ch)
			close(ch)
		}
	}()

	return ch, nil
}

func (s *Service) Publish(ctx context.Context, userID types.UserID, event eventstream.Event) error {
	if err := event.Validate(); err != nil {
		return fmt.Errorf("validating event: %v", err)
	}
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	if len(s.subs[userID]) == 0 {
		return nil
	}

	if subs, ok := s.subs[userID]; ok {
		for eventChan := range subs {
			select {
			case <-ctx.Done():
				return nil
			case eventChan <- event:
			}

			// кажется, что если не важен порядок, то можно отправлять эвенты в так же в горутине
			// чтобы другие подписчики не ждали
			// eventChan := eventChan
			// go func() {
			//	eventChan <- event
			// }()
		}
	}

	return nil
}

func (s *Service) Close() error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	for userID, chans := range s.subs {
		for eventChan := range chans {
			close(eventChan)
		}
		delete(s.subs, userID)
	}
	return nil
}
