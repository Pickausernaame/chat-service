package getchathistory

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/Pickausernaame/chat-service/internal/cursor"
	messagesrepo "github.com/Pickausernaame/chat-service/internal/repositories/messages"
	problemsrepo "github.com/Pickausernaame/chat-service/internal/repositories/problems"
	"github.com/Pickausernaame/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=getchathistorymocks

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrInvalidCursor  = errors.New("invalid cursor")
)

type messagesRepository interface {
	GetProblemMessages(
		ctx context.Context,
		problemID types.ProblemID,
		pageSize int,
		cursor *messagesrepo.Cursor,
	) ([]*messagesrepo.Message, *messagesrepo.Cursor, error)
}

type problemsRepository interface {
	GetProblemByChatAndManagerIDs(ctx context.Context, chatID types.ChatID,
		managerID types.UserID) (*problemsrepo.Problem, error)
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	msgRepo messagesRepository `option:"mandatory" validate:"required"`
	prbRepo problemsRepository `option:"mandatory" validate:"required"`
}

type UseCase struct {
	Options
	lg *zap.Logger
}

func New(opts Options) (UseCase, error) {
	if err := opts.Validate(); err != nil {
		return UseCase{}, fmt.Errorf("validating: %v", err)
	}
	return UseCase{Options: opts, lg: zap.L().Named("get-chat-history-usecase")}, nil
}

func (u UseCase) Handle(ctx context.Context, req Request) (Response, error) {
	if err := req.Validate(); err != nil {
		return Response{}, fmt.Errorf("request validation: %v %w", err, ErrInvalidRequest)
	}
	var cur *messagesrepo.Cursor
	if req.Cursor != "" {
		cur = &messagesrepo.Cursor{}
		if err := cursor.Decode(req.Cursor, cur); err != nil {
			return Response{}, ErrInvalidCursor
		}
	}

	p, err := u.prbRepo.GetProblemByChatAndManagerIDs(ctx, req.ChatID, req.ManagerID)
	if err != nil {
		return Response{}, fmt.Errorf("getting problems by chat and manager ids: %v", err)
	}

	msgs, next, err := u.msgRepo.GetProblemMessages(ctx, p.ID, req.PageSize, cur)
	if err != nil {
		if errors.Is(err, messagesrepo.ErrInvalidCursor) {
			return Response{}, fmt.Errorf("getting chat messages: %v %w", err, ErrInvalidCursor)
		}
		return Response{}, fmt.Errorf("getting chat messages: %v", err)
	}

	r := Response{
		Messages:   make([]*Message, 0, len(msgs)),
		NextCursor: "",
	}

	if next != nil {
		r.NextCursor, err = cursor.Encode(next)
		if err != nil {
			return Response{}, err
		}
	}

	for _, m := range msgs {
		r.Messages = append(r.Messages, toDTOMessage(m))
	}

	return r, nil
}
