package outbox

import (
	"context"
	"fmt"
	"time"

	"github.com/Pickausernaame/chat-service/internal/types"
)

func (s *Service) Put(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error) {
	id, err := s.jobsRepo.CreateJob(ctx, name, payload, availableAt)
	if err != nil {
		return types.JobIDNil, fmt.Errorf("creating job: %v", err)
	}
	return id, nil
}
