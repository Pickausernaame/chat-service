package outbox

import (
	"context"
	"fmt"
	"time"

	"github.com/sourcegraph/conc"
	"go.uber.org/zap"

	jobsrepo "github.com/Pickausernaame/chat-service/internal/repositories/jobs"
	"github.com/Pickausernaame/chat-service/internal/types"
)

const serviceName = "outbox"

type jobsRepository interface {
	CreateJob(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error)
	FindAndReserveJob(ctx context.Context, until time.Time) (jobsrepo.Job, error)
	DeleteJob(ctx context.Context, jobID types.JobID) error
	CreateFailedJob(ctx context.Context, name, payload, reason string) error
}

type transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

//go:generate options-gen -out-filename=service_options.gen.go -from-struct=Options
type Options struct {
	workers    int           `option:"mandatory" validate:"min=1,max=32"`
	idleTime   time.Duration `option:"mandatory" validate:"min=100ms,max=10s"`
	reserveFor time.Duration `option:"mandatory" validate:"min=1s,max=10m"`

	jobsRepo jobsRepository `option:"mandatory" validate:"required"`
	txtr     transactor     `option:"mandatory" validate:"required"`
	lg       *zap.Logger    `option:"mandatory" validate:"required"`
}

type Service struct {
	workers  []*Worker
	reg      *registry
	jobsRepo jobsRepository
}

func New(opts Options) (*Service, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("creating service %s: %v", serviceName, err)
	}
	reg := newRegistry()

	ws, err := newWorkers(
		NewWorkerOptions(opts.idleTime, opts.reserveFor, opts.jobsRepo, reg, opts.txtr, opts.lg), opts.workers)
	if err != nil {
		return nil, fmt.Errorf("creating service %s: %v", serviceName, err)
	}

	return &Service{
		workers:  ws,
		reg:      reg,
		jobsRepo: opts.jobsRepo,
	}, nil
}

func (s *Service) RegisterJob(job Job) error {
	return s.reg.set(job)
}

func (s *Service) MustRegisterJob(job Job) {
	if err := s.reg.set(job); err != nil {
		panic(err)
	}
}

func (s *Service) Run(ctx context.Context) error {
	wg := conc.NewWaitGroup()

	// starting workers
	for i, w := range s.workers {
		w := w
		i := i
		wg.Go(func() {
			w.Run(ctx, i)
		})
	}

	wg.Wait()

	return nil
}
