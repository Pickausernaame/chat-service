package messagesrepo

import (
	"context"
	"errors"
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

	switch {
	case isCursorExist:
		msgs = r.messagesByClientIDQuery(ctx, clientID).Where(
			message.And(
				message.CreatedAtLT(cursor.LastCreatedAt),
				message.IsVisibleForClient(true),
			),
		).Order(message.ByCreatedAt(sql.OrderDesc())).Limit(cursor.PageSize).AllX(ctx)
	case isPageSizeExist:
		msgs = r.messagesByClientIDQuery(ctx, clientID).Where(message.IsVisibleForClient(true)).
			Order(message.ByCreatedAt(sql.OrderDesc())).Limit(pageSize).AllX(ctx)
	}

	res := make([]*Message, 0, len(msgs))
	for _, m := range msgs {
		res = append(res, adaptStoreMessage(m))
	}

	cursor = nil
	if len(msgs) != 0 {
		c := r.messagesByClientIDQuery(ctx, clientID).
			Where(message.CreatedAtLT(msgs[len(msgs)-1].CreatedAt)).CountX(ctx)
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

func (r *Repo) messagesByClientIDQuery(ctx context.Context, clientID types.UserID) *store.MessageQuery {
	return r.db.Message(ctx).Query().Where(func(s *sql.Selector) {
		ch := sql.Table(chat.Table)
		s.Join(ch).On(s.C(message.ChatColumn), ch.C(chat.FieldID))
		s.Where(sql.EQ(ch.C(chat.FieldClientID), clientID))
	})
}
