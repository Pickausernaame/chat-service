package problemsrepo

import (
	"time"

	"github.com/Pickausernaame/chat-service/internal/store"
	"github.com/Pickausernaame/chat-service/internal/types"
)

type Problem struct {
	ID        types.ProblemID
	ChatID    types.ChatID
	ManagerID types.UserID
	ResolveAt time.Time
	CreatedAt time.Time
}

type ProblemAndClientID struct {
	*Problem
	ClientID types.UserID
}

func adaptStoreProblem(m *store.Problem) *Problem {
	return &Problem{
		ID:        m.ID,
		ChatID:    m.ChatID,
		ManagerID: m.ManagerID,
		ResolveAt: m.ResolveAt,
		CreatedAt: m.CreatedAt,
	}
}
