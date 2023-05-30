package outbox

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	"go.uber.org/zap"

	jobsrepo "github.com/Pickausernaame/chat-service/internal/repositories/jobs"
)

type jobsRegistry interface {
	get(name string) (Job, error)
}

//go:generate options-gen -out-filename=worker_options.gen.go -from-struct=WorkerOptions
type WorkerOptions struct {
	idleTime   time.Duration `option:"mandatory" validate:"min=100ms,max=10s"`
	reserveFor time.Duration `option:"mandatory" validate:"min=1s,max=10m"`

	jobsRepo jobsRepository `option:"mandatory" validate:"required"`
	jobsReg  jobsRegistry   `option:"mandatory" validate:"required"`
	txtr     transactor     `option:"mandatory" validate:"required"`
}

type worker struct {
	WorkerOptions
	lg *zap.Logger
}

func newWorkers(options WorkerOptions, count int) ([]*worker, error) {
	err := options.Validate()
	if err != nil {
		return nil, fmt.Errorf("creating worker: %v", err)
	}
	res := make([]*worker, 0, count)
	for i := 0; i < count; i++ {
		res = append(res, &worker{WorkerOptions: options, lg: zap.L().Named(fmt.Sprintf("outbox-worker-%d", i))})
	}

	return res, nil
}

func (w *worker) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:

			func() {
				var jobMeta jobsrepo.Job
				var err error
				// job processing
				defer func() {
					pErr := w.postJobProcessing(ctx, jobMeta, err)
					if pErr != nil {
						w.lg.Error("post job processing error", zap.Error(pErr), zap.Any("job", jobMeta))
					}
				}()

				// try to reserve free job
				jobMeta, err = w.jobsRepo.FindAndReserveJob(ctx, time.Now().Add(w.reserveFor))
				if err != nil {
					if errors.Is(err, jobsrepo.ErrNoJobs) {
						time.Sleep(w.idleTime)
						return
					}
					w.lg.Error("finding job error", zap.Error(err))
					return
				}

				// try to get job from registry
				var j Job
				j, err = w.jobsReg.get(jobMeta.Name)
				if err != nil {
					w.lg.Error("finding job in registry error", zap.Error(err))
					return
				}

				// try to execute job

				for {
					err = func() error {
						jCtx, cancel := context.WithTimeout(ctx, j.ExecutionTimeout())
						defer cancel()
						return j.Handle(jCtx, jobMeta.Payload)
					}()

					if nil == err {
						break
					}

					jobMeta.Attempts++
					if j.MaxAttempts() < jobMeta.Attempts {
						err = fmt.Errorf("attempts limit exceeded. job = %q attempts=%d", jobMeta.ID, jobMeta.Attempts)
						break
					}
				}

				if err != nil {
					w.lg.Error("executing job error", zap.Error(err))
					return
				}
			}()
		}
	}
}

// postJobProcessing - it is post job processing
// if job exists - delete it
// if error exists - mark this job like failed.
func (w *worker) postJobProcessing(ctx context.Context, jobMeta jobsrepo.Job, jobProcessingErr error) error {
	//nolint:nestif
	if jobMeta.Name != "" || jobProcessingErr != nil {
		err := w.txtr.RunInTx(ctx, func(ctx context.Context) error {
			if jobMeta.Name != "" {
				err := w.jobsRepo.DeleteJob(ctx, jobMeta.ID)
				if err != nil {
					return fmt.Errorf("deleting job: %v", err)
				}
			}

			if jobProcessingErr != nil && jobMeta.Name != "" {
				err := w.jobsRepo.CreateFailedJob(ctx, jobMeta.Name, jobMeta.Payload, jobProcessingErr.Error())
				if err != nil {
					log.Info("try to set failed job", zap.Any("job", jobMeta))
					return fmt.Errorf("creating faailed job: %v", err)
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}
