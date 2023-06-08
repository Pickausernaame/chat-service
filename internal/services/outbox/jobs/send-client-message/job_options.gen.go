// Code generated by options-gen. DO NOT EDIT.
package sendclientmessagejob

import (
	fmt461e464ebed9 "fmt"

	errors461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/errors"
	validator461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/validator"
)

type OptOptionsSetter func(o *Options)

func NewOptions(
	msgProd messageProducer,
	msgRepo messageRepository,
	eventStream eventStream,
	options ...OptOptionsSetter,
) Options {
	o := Options{}

	// Setting defaults from field tag (if present)

	o.msgProd = msgProd
	o.msgRepo = msgRepo
	o.eventStream = eventStream

	for _, opt := range options {
		opt(&o)
	}
	return o
}

func (o *Options) Validate() error {
	errs := new(errors461e464ebed9.ValidationErrors)
	errs.Add(errors461e464ebed9.NewValidationError("msgProd", _validate_Options_msgProd(o)))
	errs.Add(errors461e464ebed9.NewValidationError("msgRepo", _validate_Options_msgRepo(o)))
	errs.Add(errors461e464ebed9.NewValidationError("eventStream", _validate_Options_eventStream(o)))
	return errs.AsError()
}

func _validate_Options_msgProd(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.msgProd, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `msgProd` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_msgRepo(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.msgRepo, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `msgRepo` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_eventStream(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.eventStream, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `eventStream` did not pass the test: %w", err)
	}
	return nil
}
