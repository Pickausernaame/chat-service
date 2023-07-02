package messagesrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/Pickausernaame/chat-service/internal/store"
	"github.com/Pickausernaame/chat-service/internal/store/chat"
	"github.com/Pickausernaame/chat-service/internal/store/message"
	"github.com/Pickausernaame/chat-service/internal/types"
)

var (
	ErrInvalidPageSize = errors.New("invalid page size")
	ErrInvalidCursor   = errors.New("invalid cursor")
)

const (
	minPageSize = 10
	maxPageSize = 100
)

type Cursor struct {
	LastCreatedAt time.Time
	PageSize      int
}

// GetClientChatMessages returns Nth page of messages in the chat for client side.
func (r *Repo) GetClientChatMessages(
	ctx context.Context,
	clientID types.UserID,
	pageSize int,
	cursor *Cursor,
) ([]*Message, *Cursor, error) {
	query := r.messagesByClientIDQuery(ctx, clientID)
	return r.getMessages(ctx, query, pageSize, cursor)
}

// GetProblemMessages returns Nth page of messages in the chat for manager side (specific problem).
func (r *Repo) GetProblemMessages(
	ctx context.Context,
	problemID types.ProblemID,
	pageSize int,
	cursor *Cursor,
) ([]*Message, *Cursor, error) {
	query := r.messagesByProblemIDQuery(ctx, problemID)
	return r.getMessages(ctx, query, pageSize, cursor)
}

func (r *Repo) messagesByClientIDQuery(ctx context.Context, clientID types.UserID) *store.MessageQuery {
	return r.db.Message(ctx).Query().Where(func(s *sql.Selector) {
		ch := sql.Table(chat.Table)
		s.Join(ch).On(s.C(message.ChatColumn), ch.C(chat.FieldID))
		s.Where(sql.EQ(ch.C(chat.FieldClientID), clientID))
	}).Where(message.IsVisibleForClient(true))
}

func (r *Repo) messagesByProblemIDQuery(ctx context.Context, problemID types.ProblemID) *store.MessageQuery {
	return r.db.Message(ctx).Query().Where(message.ProblemID(problemID), message.IsVisibleForManager(true))
}

func (r *Repo) getMessages(
	ctx context.Context,
	query *store.MessageQuery,
	pageSize int,
	cursor *Cursor,
) ([]*Message, *Cursor, error) {
	var msgs []*store.Message
	isPageSizeExist := false
	isCursorExist := false
	cursorPageSize := 0

	// check pageSize
	if pageSize != 0 {
		if pageSize < minPageSize || pageSize > maxPageSize {
			return nil, nil, ErrInvalidPageSize
		}
		isPageSizeExist = true
	}

	// check cursor
	if cursor != nil {
		if cursor.PageSize < minPageSize || cursor.PageSize > maxPageSize {
			return nil, nil, ErrInvalidCursor
		}

		if cursor.LastCreatedAt.Equal(time.Time{}) {
			return nil, nil, ErrInvalidCursor
		}

		isCursorExist = true
		cursorPageSize = cursor.PageSize
	}

	var err error
	switch {
	case isCursorExist:
		msgs, err = query.Where(message.CreatedAtLT(cursor.LastCreatedAt)).
			Order(message.ByCreatedAt(sql.OrderDesc())).Limit(cursor.PageSize).All(ctx)
	case isPageSizeExist:
		msgs, err = query.Order(message.ByCreatedAt(sql.OrderDesc())).Limit(pageSize).All(ctx)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("getting list: %v", err)
	}

	res := make([]*Message, 0, len(msgs))
	for _, m := range msgs {
		res = append(res, adaptStoreMessage(m))
	}

	cursor = nil
	if len(msgs) != 0 {
		c := query.Where(message.CreatedAtLT(msgs[len(msgs)-1].CreatedAt)).CountX(ctx)
		if c > 0 {
			cursor = &Cursor{
				LastCreatedAt: msgs[len(msgs)-1].CreatedAt,
				PageSize:      pageSize,
			}
			if cursorPageSize != 0 {
				cursor.PageSize = cursorPageSize
			}
		}
	}

	return res, cursor, nil
}
