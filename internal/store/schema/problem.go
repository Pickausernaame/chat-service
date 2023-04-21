package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/Pickausernaame/chat-service/internal/types"
)

// Problem holds the schema definition for the Problem entity.
type Problem struct {
	ent.Schema
}

// Fields of the Problem.
func (Problem) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", types.ProblemID{}).Immutable().Unique().Default(types.NewProblemID),
		field.UUID("chat_id", types.ChatID{}).Immutable(),
		field.UUID("manager_id", types.UserID{}).Nillable(),
		field.Time("resolve_at").Optional(),
		field.Time("created_at").Immutable().Default(time.Now()),
	}
}

// Edges of the Problem.
func (Problem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("chat", Chat.Type).Immutable().
			Ref("problems").Unique().Field("chat_id").Required(),

		edge.To("messages", Message.Type),
	}
}
