package messagesrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/Pickausernaame/chat-service/internal/store"
	"github.com/Pickausernaame/chat-service/internal/store/message"
	"github.com/Pickausernaame/chat-service/internal/types"
)

var ErrMsgNotFound = errors.New("message not found")

func (r *Repo) GetMessageByRequestID(ctx context.Context, reqID types.RequestID) (*Message, error) {
	m, err := r.db.Message(ctx).Query().Where(message.InitialRequestID(reqID)).Only(ctx)
	if err != nil {
		if store.IsNotFound(err) {
			return nil, ErrMsgNotFound
		}
		return nil, err
	}
	return adaptStoreMessage(m), nil
}

func (r *Repo) GetMessageByID(ctx context.Context, id types.MessageID) (*Message, error) {
	m, err := r.db.Message(ctx).Get(ctx, id)
	if err != nil {
		if store.IsNotFound(err) {
			return nil, ErrMsgNotFound
		}
		return nil, err
	}
	return adaptStoreMessage(m), nil
}

// CreateClientVisible creates a message that is visible only to the client.
func (r *Repo) CreateClientVisible(
	ctx context.Context,
	reqID types.RequestID,
	problemID types.ProblemID,
	chatID types.ChatID,
	authorID types.UserID,
	msgBody string,
) (*Message, error) {
	m, err := r.db.Message(ctx).Create().
		SetInitialRequestID(reqID).
		SetProblemID(problemID).
		SetChatID(chatID).
		SetAuthorID(authorID).
		SetBody(msgBody).
		SetIsVisibleForClient(true).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating client visible msg: %v", err)
	}
	return adaptStoreMessage(m), nil
}

func (r *Repo) MessageForManagerByChatID(ctx context.Context, id types.ChatID) (*Message, error) {
	m, err := r.db.Message(ctx).Query().
		Where(
			message.ChatID(id),
			message.IsVisibleForManager(true),
			message.IsBlocked(false)).
		Order(message.ByCreatedAt()).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting message visible msg: %v", err)
	}
	return adaptStoreMessage(m), nil
}

func (r *Repo) CreateProblemAssignedMessage(
	ctx context.Context,
	id types.ChatID,
	managerID types.UserID,
	problemID types.ProblemID,
) (*store.Message, error) {
	return r.db.Message(ctx).Create().
		SetChatID(id).
		SetIsVisibleForClient(true).
		SetIsService(true).
		SetProblemID(problemID).
		SetBody(fmt.Sprintf("Manager %s will answer you", managerID.String())).
		Save(ctx)
}
