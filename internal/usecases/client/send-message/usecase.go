package sendmessage

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	messagesrepo "github.com/Pickausernaame/chat-service/internal/repositories/messages"
	"github.com/Pickausernaame/chat-service/internal/types"
	"github.com/Pickausernaame/chat-service/internal/validator"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=sendmessagemocks

var (
	ErrInvalidRequest    = errors.New("invalid request")
	ErrChatNotCreated    = errors.New("chat not created")
	ErrProblemNotCreated = errors.New("problem not created")
)

type chatsRepository interface {
	CreateIfNotExists(ctx context.Context, userID types.UserID) (types.ChatID, error)
}

type messagesRepository interface {
	GetMessageByRequestID(ctx context.Context, reqID types.RequestID) (*messagesrepo.Message, error)
	CreateClientVisible(
		ctx context.Context,
		reqID types.RequestID,
		problemID types.ProblemID,
		chatID types.ChatID,
		authorID types.UserID,
		msgBody string,
	) (*messagesrepo.Message, error)
}

type problemsRepository interface {
	CreateIfNotExists(ctx context.Context, chatID types.ChatID) (types.ProblemID, error)
}

type transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	chatRepo    chatsRepository    `option:"mandatory" validate:"required"`
	msgRepo     messagesRepository `option:"mandatory" validate:"required"`
	problemRepo problemsRepository `option:"mandatory" validate:"required"`
	txr         transactor         `option:"mandatory" validate:"required"`
	lg          *zap.Logger
}

type UseCase struct {
	Options
}

func New(opts Options) (UseCase, error) {
	return UseCase{Options: opts}, opts.Validate()
}

func (u UseCase) Handle(ctx context.Context, req Request) (Response, error) {
	if err := validator.Validator.Struct(req); err != nil {
		return Response{}, fmt.Errorf("validate request: %w", ErrInvalidRequest)
	}

	var msg *messagesrepo.Message
	var err error
	err = u.txr.RunInTx(ctx, func(ctx context.Context) error {
		msg, err = u.msgRepo.GetMessageByRequestID(ctx, req.ID)
		if err != nil && !errors.Is(err, messagesrepo.ErrMsgNotFound) {
			return err
		}

		// if msg exist - return it
		if nil == err && msg != nil {
			return nil
		}

		// if msg not exist - create it!
		chatID, err := u.chatRepo.CreateIfNotExists(ctx, req.ClientID)
		if err != nil {
			u.lg.Error("chat.CreateIfNotExists error", zap.Error(err),
				zap.Stringer("userID", req.ClientID),
				zap.Stringer("reqID", req.ID))
			return fmt.Errorf("chat.CreateIfNotExists error: %w", ErrChatNotCreated)
		}

		problemID, err := u.problemRepo.CreateIfNotExists(ctx, chatID)
		if err != nil {
			u.lg.Error("problem.CreateIfNotExists error", zap.Error(err),
				zap.Stringer("userID", req.ClientID),
				zap.Stringer("reqID", req.ID),
				zap.Stringer("chatID", chatID))
			return fmt.Errorf("problem.CreateIfNotExists error: %w", ErrProblemNotCreated)
		}

		msg, err = u.msgRepo.CreateClientVisible(ctx, req.ID, problemID, chatID, req.ClientID, req.MessageBody)
		if err != nil {
			return fmt.Errorf("message.CreateClientVisible error: %v", err)
		}
		return nil
	})
	if err != nil {
		return Response{}, err
	}

	return Response{
		AuthorID:  msg.AuthorID,
		MessageID: msg.ID,
		CreatedAt: msg.CreatedAt,
	}, nil
}
