package sendclientmessagejob

import (
	"context"
	"fmt"

	messagesrepo "github.com/Pickausernaame/chat-service/internal/repositories/messages"
	eventstream "github.com/Pickausernaame/chat-service/internal/services/event-stream"
	msgproducer "github.com/Pickausernaame/chat-service/internal/services/msg-producer"
	"github.com/Pickausernaame/chat-service/internal/services/outbox"
	"github.com/Pickausernaame/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/job_mock.gen.go -package=sendclientmessagejobmocks

const Name = "send-client-message"

type messageProducer interface {
	ProduceMessage(ctx context.Context, message msgproducer.Message) error
}

type messageRepository interface {
	GetMessageByID(ctx context.Context, msgID types.MessageID) (*messagesrepo.Message, error)
}

type eventStream interface {
	Publish(ctx context.Context, userID types.UserID, event eventstream.Event) error
}

//go:generate options-gen -out-filename=job_options.gen.go -from-struct=Options
type Options struct {
	msgProd     messageProducer   `option:"mandatory" validate:"required"`
	msgRepo     messageRepository `option:"mandatory" validate:"required"`
	eventStream eventStream       `option:"mandatory" validate:"required"`
}

type Job struct {
	Options
	outbox.DefaultJob
}

func New(opts Options) (*Job, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validations opts: %v", err)
	}

	return &Job{Options: opts}, nil
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

	err = j.msgProd.ProduceMessage(ctx, msgproducer.Message{
		ID:         msg.ID,
		ChatID:     msg.ChatID,
		Body:       msg.Body,
		FromClient: !msg.IsService,
	})
	if err != nil {
		return fmt.Errorf("producing msg: %v", err)
	}

	event := eventstream.NewNewMessageEvent(types.NewEventID(), msg.InitialRequestID,
		msg.ChatID, msg.ID, msg.AuthorID, msg.CreatedAt, msg.Body, msg.IsService)
	err = j.eventStream.Publish(ctx, msg.AuthorID, event)
	if err != nil {
		return fmt.Errorf("publishing event: %v", err)
	}

	return nil
}
