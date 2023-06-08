package problemsrepo

import (
	"context"
	"errors"
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/Pickausernaame/chat-service/internal/store"
	"github.com/Pickausernaame/chat-service/internal/store/problem"
	"github.com/Pickausernaame/chat-service/internal/types"
)

var ErrProblemNotFound = errors.New("problem not found")

func (r *Repo) CreateIfNotExists(ctx context.Context, chatID types.ChatID) (types.ProblemID, error) {
	return r.db.Problem(ctx).Create().SetChatID(chatID).
		OnConflict(
			sql.ConflictColumns(problem.ChatColumn),
			sql.ConflictWhere(sql.IsNull(problem.FieldResolveAt)),
		).Ignore().ID(ctx)
}

func (r *Repo) GetManagerOpenProblemsCount(ctx context.Context, managerID types.UserID) (int, error) {
	return r.db.Problem(ctx).Query().
		Where(problem.And(
			problem.ManagerID(managerID),
			problem.ResolveAtIsNil()),
		).Count(ctx)
}

func (r *Repo) GetAssignedUnsolvedProblems(ctx context.Context, managerID types.UserID) ([]*Problem, error) {
	ps, err := r.db.Problem(ctx).Query().
		Where(problem.And(
			problem.ResolveAtIsNil(),
			problem.ManagerID(managerID)),
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

func (r *Repo) GetManagerIDByChatID(ctx context.Context, chatID types.ChatID) (types.UserID, error) {
	p, err := r.db.Problem(ctx).Query().
		Where(problem.And(problem.ChatID(chatID), problem.ResolveAtIsNil())).First(ctx)
	if err != nil {
		return types.UserIDNil, err
	}
	return p.ManagerID, nil
}

func (r *Repo) GetProblemByChatAndManagerIDs(ctx context.Context, chatID types.ChatID,
	managerID types.UserID,
) (*Problem, error) {
	p, err := r.db.Problem(ctx).Query().Where(problem.And(problem.ChatID(chatID), problem.ManagerID(managerID),
		problem.ResolveAtIsNil())).First(ctx)
	if err != nil {
		if store.IsNotFound(err) {
			return nil, ErrProblemNotFound
		}
		return nil, err
	}
	return adaptStoreProblem(p), nil
}

func (r *Repo) ResolveProblem(ctx context.Context, problemID types.ProblemID, managerID types.UserID) error {
	predicate := problem.And(problem.ID(problemID), problem.ManagerID(managerID), problem.ResolveAtIsNil())
	exist, err := r.db.Problem(ctx).Query().Where(predicate).Exist(ctx)
	if err != nil {
		return err
	}
	if !exist {
		return ErrProblemNotFound
	}
	return r.db.Problem(ctx).Update().SetResolveAt(time.Now()).Where(predicate).Exec(ctx)
}
