package jobresolveproblem

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	messagesrepo "github.com/Pickausernaame/chat-service/internal/repositories/messages"
	eventstream "github.com/Pickausernaame/chat-service/internal/services/event-stream"
	"github.com/Pickausernaame/chat-service/internal/services/outbox"
	"github.com/Pickausernaame/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/job_mock.gen.go -package=jobresolveproblemmocks

const Name = "resolve-problem"

type chatRepository interface {
	ClientIDByID(ctx context.Context, id types.ChatID) (types.UserID, error)
}

type messageRepository interface {
	GetMessageByID(ctx context.Context, msgID types.MessageID) (*messagesrepo.Message, error)
}

type eventStream interface {
	Publish(ctx context.Context, userID types.UserID, event eventstream.Event) error
}

type managerLoad interface {
	CanManagerTakeProblem(ctx context.Context, managerID types.UserID) (bool, error)
}

//go:generate options-gen -out-filename=job_options.gen.go -from-struct=Options
type Options struct {
	msgRepo     messageRepository `option:"mandatory" validate:"required"`
	chatRepo    chatRepository    `option:"mandatory" validate:"required"`
	managerLoad managerLoad       `option:"mandatory" validate:"required"`
	eventStream eventStream       `option:"mandatory" validate:"required"`
}

type Job struct {
	Options
	lg *zap.Logger
	outbox.DefaultJob
}

func New(opts Options) (*Job, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validations opts: %v", err)
	}
	return &Job{Options: opts, lg: zap.L().Named(Name)}, nil
}

func (j *Job) Name() string {
	return Name
}

func (j *Job) Handle(ctx context.Context, payload string) error {
	req := &request{}

	if err := json.Unmarshal([]byte(payload), req); err != nil {
		return fmt.Errorf("unmarshaling payload: %v", err)
	}

	msg, err := j.msgRepo.GetMessageByID(ctx, req.MessageID)
	if err != nil {
		return fmt.Errorf("getting msg by id: %v", err)
	}

	ok, err := j.managerLoad.CanManagerTakeProblem(ctx, req.ManagerID)
	if err != nil {
		return fmt.Errorf("getting manager takes problem: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)
	defer func() {
		if err := eg.Wait(); err != nil {
			j.lg.Error("error group error", zap.Error(err))
		}
	}()

	eg.Go(func() error {
		event := eventstream.NewChatClosedEvent(types.NewEventID(), req.RequestID, ok, req.ChatID)
		err := j.eventStream.Publish(ctx, req.ManagerID, event)
		if err != nil {
			return fmt.Errorf("publishing closed chat event: %v", err)
		}
		return nil
	})

	clientID, err := j.chatRepo.ClientIDByID(ctx, msg.ChatID)
	if err != nil {
		return fmt.Errorf("getting client by chatID: %v", err)
	}

	event := eventstream.NewNewMessageEvent(types.NewEventID(), msg.InitialRequestID,
		msg.ChatID, msg.ID, msg.AuthorID, msg.CreatedAt, msg.Body, msg.IsService)
	err = j.eventStream.Publish(ctx, clientID, event)
	if err != nil {
		return fmt.Errorf("publishing event: %v", err)
	}
	return nil
}
