package problemsrepo

import (
	"context"
	"time"

	"github.com/Pickausernaame/chat-service/internal/store"
	"github.com/Pickausernaame/chat-service/internal/store/message"
	"github.com/Pickausernaame/chat-service/internal/store/problem"
	"github.com/Pickausernaame/chat-service/internal/types"
)

func (r *Repo) GetUnassignedProblems(ctx context.Context) ([]*store.Problem, error) {
	return r.db.Problem(ctx).Query().Where(problem.ResolveAtIsNil(), problem.ManagerIDIsNil()).Order(problem.ByCreatedAt()).All(ctx)
}

func (r *Repo) AssignManager(ctx context.Context, problemID types.ProblemID, managerID types.UserID) error {
	return r.db.Problem(ctx).Update().SetManagerID(managerID).Where(problem.ID(problemID)).Exec(ctx)
}

func (r *Repo) ResolveProblem(ctx context.Context, problemID types.ProblemID) error {
	return r.db.Problem(ctx).Update().SetResolveAt(time.Now()).Where(problem.ID(problemID)).Exec(ctx)
}

func (r *Repo) GetRequestID(ctx context.Context, problemID types.ProblemID) (types.RequestID, error) {
	p, err := r.db.Problem(ctx).Get(ctx, problemID)
	if err != nil {
		return types.RequestIDNil, err
	}
	msg, err := p.QueryMessages().Order(message.ByCreatedAt()).Only(ctx)
	if err != nil {
		return types.RequestIDNil, err
	}
	return msg.InitialRequestID, nil
}
