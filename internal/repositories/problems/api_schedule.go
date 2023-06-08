package problemsrepo

import (
	"context"

	"github.com/Pickausernaame/chat-service/internal/store/problem"
	"github.com/Pickausernaame/chat-service/internal/types"
)

func (r *Repo) GetUnassignedProblems(ctx context.Context) ([]*Problem, error) {
	ps, err := r.db.Problem(ctx).Query().
		Where(problem.And(
			problem.ResolveAtIsNil(),
			problem.ManagerIDIsNil()),
		).Order(problem.ByCreatedAt()).All(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]*Problem, 0, len(ps))
	for _, p := range ps {
		res = append(res, adaptStoreProblem(p))
	}
	return res, nil
}

func (r *Repo) AssignManager(ctx context.Context, problemID types.ProblemID,
	managerID types.UserID,
) error {
	return r.db.Problem(ctx).Update().SetManagerID(managerID).
		Where(problem.ID(problemID)).Exec(ctx)
}
