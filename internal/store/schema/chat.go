package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/Pickausernaame/chat-service/internal/types"
)

// Chat holds the schema definition for the Chat entity.
type Chat struct {
	ent.Schema
}

// Fields of the Chat.
func (Chat) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", types.ChatID{}).Immutable().Unique().Default(types.NewChatID),
		field.UUID("client_id", types.UserID{}).Immutable().Unique(),
		field.Time("created_at").Immutable().Default(time.Now),
	}
}

// Edges of the Chat.
func (Chat) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("problems", Problem.Type),
		edge.To("messages", Message.Type),
	}
}
