package problemsrepo

import (
	"context"

	"entgo.io/ent/dialect/sql"

	"github.com/Pickausernaame/chat-service/internal/store/problem"
	"github.com/Pickausernaame/chat-service/internal/types"
)

func (r *Repo) CreateIfNotExists(ctx context.Context, chatID types.ChatID) (types.ProblemID, error) {
	return r.db.Problem(ctx).Create().SetChatID(chatID).
		OnConflict(
			sql.ConflictColumns(problem.ChatColumn),
			sql.ConflictWhere(sql.IsNull(problem.FieldResolveAt)),
		).Ignore().ID(ctx)
}

func (r *Repo) GetManagerOpenProblemsCount(ctx context.Context, managerID types.UserID) (int, error) {
	return r.db.Problem(ctx).Query().Where(problem.And(problem.ManagerID(managerID), problem.ResolveAtIsNil())).Count(ctx)
}

func (r *Repo) GetAssignedUnsolvedProblems(ctx context.Context, managerID types.UserID) ([]*Problem, error) {
	ps, err := r.db.Problem(ctx).Query().Where(problem.ResolveAtIsNil(), problem.ManagerID(managerID)).Order(problem.ByCreatedAt()).All(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]*Problem, 0, len(ps))
	for _, p := range ps {
		res = append(res, adaptStoreProblem(p))
	}
	return res, nil
}
