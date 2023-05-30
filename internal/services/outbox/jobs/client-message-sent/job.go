package clientmessagesentjob

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	messagesrepo "github.com/Pickausernaame/chat-service/internal/repositories/messages"
	eventstream "github.com/Pickausernaame/chat-service/internal/services/event-stream"
	"github.com/Pickausernaame/chat-service/internal/services/outbox"
	"github.com/Pickausernaame/chat-service/internal/types"
)

const Name = "client-message-sent"

type messageRepository interface {
	GetMessageByID(ctx context.Context, msgID types.MessageID) (*messagesrepo.Message, error)
}

type eventStream interface {
	Publish(ctx context.Context, userID types.UserID, event eventstream.Event) error
}

//go:generate options-gen -out-filename=job_options.gen.go -from-struct=Options
type Options struct {
	msgRepo     messageRepository `option:"mandatory" validate:"required"`
	eventStream eventStream       `option:"mandatory" validate:"required"`
}

type Job struct {
	Options
	lg *zap.Logger
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
	var id types.MessageID

	if err := id.UnmarshalText([]byte(payload)); err != nil {
		return fmt.Errorf("unmarshaling payload: %v", err)
	}

	msg, err := j.msgRepo.GetMessageByID(ctx, id)
	if err != nil {
		return fmt.Errorf("getting msg by id: %v", err)
	}

	msgSentEv := eventstream.NewMessageSentEvent(types.NewEventID(), msg.InitialRequestID, msg.ID)
	err = j.eventStream.Publish(ctx, msg.AuthorID, msgSentEv)
	if err != nil {
		return fmt.Errorf("publishing event: %v", err)
	}

	newMessageEv := eventstream.NewNewMessageEvent(types.NewEventID(), msg.InitialRequestID,
		msg.ChatID, msg.ID, msg.AuthorID, msg.CreatedAt, msg.Body, msg.IsService)

	err = j.eventStream.Publish(ctx, msg.AuthorID, newMessageEv)
	if err != nil {
		return fmt.Errorf("publishing event: %v", err)
	}

	return nil
}

func (j *Job) ExecutionTimeout() time.Duration {
	return outbox.DefaultJob{}.ExecutionTimeout()
}

func (j *Job) MaxAttempts() int {
	return outbox.DefaultJob{}.MaxAttempts()
}
