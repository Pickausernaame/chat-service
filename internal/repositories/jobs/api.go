package jobsrepo

import (
	"context"
	"errors"
	"time"

	"github.com/Pickausernaame/chat-service/internal/store"
	"github.com/Pickausernaame/chat-service/internal/store/schema"
	"github.com/Pickausernaame/chat-service/internal/types"
)

var ErrNoJobs = errors.New("no jobs found")

type Job struct {
	ID       types.JobID
	Name     string
	Payload  string
	Attempts int
}

const findAndReserveJobQuery = `WITH cte AS (
  SELECT *
  FROM jobs
  WHERE available_at <= NOW()
    AND (reserved_until < NOW() OR reserved_until IS NULL)
    AND attempts < $1
  LIMIT 1
  FOR UPDATE
)
UPDATE jobs
SET attempts = attempts + 1,
    reserved_until = $2
WHERE id = (SELECT id FROM cte)
RETURNING id, name, payload, attempts;`

func (r *Repo) FindAndReserveJob(ctx context.Context, until time.Time) (Job, error) {
	rows, err := r.db.Job(ctx).QueryContext(ctx, findAndReserveJobQuery, schema.JobMaxAttempts, until)
	if err != nil {
		if store.IsNotFound(err) {
			return Job{}, ErrNoJobs
		}
		return Job{}, err
	}

	defer func() {
		_ = rows.Close()
	}()
	res := Job{}

	for rows.Next() {
		err = rows.Scan(&res.ID, &res.Name, &res.Payload, &res.Attempts)
		if err != nil {
			return Job{}, err
		}
	}
	if err = rows.Err(); err != nil {
		return Job{}, err
	}

	nilJob := Job{}
	if res == nilJob {
		return Job{}, ErrNoJobs
	}
	return res, nil
}

func (r *Repo) CreateJob(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error) {
	j, err := r.db.Job(ctx).Create().SetName(name).SetPayload(payload).
		SetAvailableAt(availableAt).SetReservedUntil(availableAt).Save(ctx)
	if err != nil {
		return types.JobIDNil, err
	}
	return j.ID, nil
}

func (r *Repo) CreateFailedJob(ctx context.Context, name, payload, reason string) error {
	_, err := r.db.FailedJob(ctx).Create().SetName(name).SetPayload(payload).SetReason(reason).Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) DeleteJob(ctx context.Context, jobID types.JobID) error {
	return r.db.Job(ctx).DeleteOneID(jobID).Exec(ctx)
}
