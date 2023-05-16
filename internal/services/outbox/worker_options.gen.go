// Code generated by options-gen. DO NOT EDIT.
package outbox

import (
	fmt461e464ebed9 "fmt"
	"time"

	errors461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/errors"
	validator461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/validator"
	"go.uber.org/zap"
)

type OptWorkerOptionsSetter func(o *WorkerOptions)

func NewWorkerOptions(
	idleTime time.Duration,
	reserveFor time.Duration,
	jobsRepo jobsRepository,
	jobsReg jobsRegistry,
	txtr transactor,
	lg *zap.Logger,
	options ...OptWorkerOptionsSetter,
) WorkerOptions {
	o := WorkerOptions{}

	// Setting defaults from field tag (if present)

	o.idleTime = idleTime
	o.reserveFor = reserveFor
	o.jobsRepo = jobsRepo
	o.jobsReg = jobsReg
	o.txtr = txtr
	o.lg = lg

	for _, opt := range options {
		opt(&o)
	}
	return o
}

func (o *WorkerOptions) Validate() error {
	errs := new(errors461e464ebed9.ValidationErrors)
	errs.Add(errors461e464ebed9.NewValidationError("idleTime", _validate_WorkerOptions_idleTime(o)))
	errs.Add(errors461e464ebed9.NewValidationError("reserveFor", _validate_WorkerOptions_reserveFor(o)))
	errs.Add(errors461e464ebed9.NewValidationError("jobsRepo", _validate_WorkerOptions_jobsRepo(o)))
	errs.Add(errors461e464ebed9.NewValidationError("jobsReg", _validate_WorkerOptions_jobsReg(o)))
	errs.Add(errors461e464ebed9.NewValidationError("txtr", _validate_WorkerOptions_txtr(o)))
	errs.Add(errors461e464ebed9.NewValidationError("lg", _validate_WorkerOptions_lg(o)))
	return errs.AsError()
}

func _validate_WorkerOptions_idleTime(o *WorkerOptions) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.idleTime, "min=100ms,max=10s"); err != nil {
		return fmt461e464ebed9.Errorf("field `idleTime` did not pass the test: %w", err)
	}
	return nil
}

func _validate_WorkerOptions_reserveFor(o *WorkerOptions) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.reserveFor, "min=1s,max=10m"); err != nil {
		return fmt461e464ebed9.Errorf("field `reserveFor` did not pass the test: %w", err)
	}
	return nil
}

func _validate_WorkerOptions_jobsRepo(o *WorkerOptions) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.jobsRepo, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `jobsRepo` did not pass the test: %w", err)
	}
	return nil
}

func _validate_WorkerOptions_jobsReg(o *WorkerOptions) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.jobsReg, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `jobsReg` did not pass the test: %w", err)
	}
	return nil
}

func _validate_WorkerOptions_txtr(o *WorkerOptions) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.txtr, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `txtr` did not pass the test: %w", err)
	}
	return nil
}

func _validate_WorkerOptions_lg(o *WorkerOptions) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.lg, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `lg` did not pass the test: %w", err)
	}
	return nil
}
