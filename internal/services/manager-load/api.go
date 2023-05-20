package managerload

import (
	"context"
	"fmt"

	"github.com/Pickausernaame/chat-service/internal/types"
)

func (s *Service) CanManagerTakeProblem(ctx context.Context, managerID types.UserID) (bool, error) {
	count, err := s.problemsRepo.GetManagerOpenProblemsCount(ctx, managerID)
	if err != nil {
		return false, fmt.Errorf("geeting open problems count: %v", err)
	}
	return s.maxProblemsAtTime > count, nil
}
