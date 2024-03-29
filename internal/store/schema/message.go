package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/Pickausernaame/chat-service/internal/types"
)

// Message holds the schema definition for the Message entity.
type Message struct {
	ent.Schema
}

// Fields of the Message.
func (Message) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", types.MessageID{}).Immutable().Unique().Default(types.NewMessageID),
		field.UUID("chat_id", types.ChatID{}).Immutable(),
		field.UUID("problem_id", types.ProblemID{}).Immutable(),
		field.UUID("author_id", types.UserID{}).Immutable().Optional(),
		field.UUID("initial_request_id", types.RequestID{}).Unique().Optional(),
		field.Bool("is_visible_for_client").Immutable().Default(false),
		field.Bool("is_visible_for_manager").Default(false),
		field.String("body").MaxLen(3000),
		field.Time("checked_at").Optional(),
		field.Bool("is_blocked").Default(false),
		field.Bool("is_service").Default(false),
		field.Time("created_at").Immutable().Default(time.Now),
	}
}

// Edges of the Message.
func (Message) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("problem", Problem.Type).Immutable().
			Ref("messages").Unique().Field("problem_id").Required(),
		edge.From("chat", Chat.Type).Immutable().
			Ref("messages").Unique().Field("chat_id").Required(),
	}
}

func (Message) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("created_at"),
		index.Fields("chat_id"),
	}
}
