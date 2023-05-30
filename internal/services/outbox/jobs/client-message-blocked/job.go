package clientmessageblockedjob

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

//go:generate mockgen -source=$GOFILE -destination=mocks/job_mock.gen.go -package=sendclientmessagejobmocks

const Name = "client-message-blocked"

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

	event := eventstream.NewMessageBlockedEvent(types.NewEventID(), msg.InitialRequestID, msg.ID)
	err = j.eventStream.Publish(ctx, msg.AuthorID, event)
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
