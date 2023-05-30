package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/Pickausernaame/chat-service/internal/types"
)

// JobMaxAttempts is some limit as protection from endless retries of outbox jobs.
const JobMaxAttempts = 30

type Job struct {
	ent.Schema
}

func (Job) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", types.JobID{}).Default(types.NewJobID).Unique().Immutable(),
		field.String("name").Immutable().NotEmpty(),
		field.String("payload").Immutable().NotEmpty(),
		field.Int("attempts").Default(0).Max(JobMaxAttempts),
		field.Time("available_at").Immutable().Optional(),
		field.Time("reserved_until").Optional(),
		field.Time("created_at").Immutable().Default(time.Now),
	}
}

func (Job) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("attempts", "reserved_until"),
	}
}

type FailedJob struct {
	ent.Schema
}

func (FailedJob) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", types.FailedJobID{}).Default(types.NewFailedJobID).Unique().Immutable(),
		field.String("name").Immutable().NotEmpty(),
		field.String("payload").Immutable().NotEmpty(),
		field.String("reason").Immutable().NotEmpty(),
		field.Time("created_at").Immutable().Default(time.Now),
	}
}
