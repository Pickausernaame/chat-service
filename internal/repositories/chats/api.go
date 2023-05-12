package chatsrepo

import (
	"context"

	"github.com/Pickausernaame/chat-service/internal/store/chat"
	"github.com/Pickausernaame/chat-service/internal/types"
)

func (r *Repo) CreateIfNotExists(ctx context.Context, userID types.UserID) (types.ChatID, error) {
	return r.db.Chat(ctx).Create().
		SetClientID(userID).
		OnConflictColumns(chat.FieldClientID).Ignore().ID(ctx)
}
