package inmemmanagerpool

import (
	"context"
	"errors"
	"sync"
	"time"

	managerpool "github.com/Pickausernaame/chat-service/internal/services/manager-pool"
	"github.com/Pickausernaame/chat-service/internal/types"
)

const (
	serviceName = "manager-pool"
	managersMax = 1000
)

type Manager struct {
	ID      types.UserID
	AddedAt time.Time
}

type Service struct {
	managerSet  map[types.UserID]struct{}
	managerPool []types.UserID
	mtx         sync.RWMutex
}

func New() *Service {
	return &Service{
		managerSet:  map[types.UserID]struct{}{},
		managerPool: make([]types.UserID, 0, managersMax),
		mtx:         sync.RWMutex{},
	}
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) Get(_ context.Context) (types.UserID, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if len(s.managerPool) == 0 {
		return types.UserIDNil, managerpool.ErrNoAvailableManagers
	}

	id := s.managerPool[0]
	// delete this manager
	s.managerPool = s.managerPool[1:]
	delete(s.managerSet, id)
	return id, nil
}

func (s *Service) Put(_ context.Context, managerID types.UserID) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, ok := s.managerSet[managerID]; ok {
		return nil
	}

	if len(s.managerPool) == managersMax {
		return errors.New("manager limit exceeded")
	}

	s.managerSet[managerID] = struct{}{}
	s.managerPool = append(s.managerPool, managerID)

	return nil
}

func (s *Service) Contains(_ context.Context, managerID types.UserID) (bool, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	_, ok := s.managerSet[managerID]
	return ok, nil
}

func (s *Service) Size() int {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return len(s.managerPool)
}
