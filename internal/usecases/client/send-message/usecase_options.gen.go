// Code generated by options-gen. DO NOT EDIT.
package sendmessage

import (
	fmt461e464ebed9 "fmt"

	errors461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/errors"
	validator461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/validator"
	"go.uber.org/zap"
)

type OptOptionsSetter func(o *Options)

func NewOptions(
	chatRepo chatsRepository,
	msgRepo messagesRepository,
	problemRepo problemsRepository,
	txr transactor,
	options ...OptOptionsSetter,
) Options {
	o := Options{}

	// Setting defaults from field tag (if present)

	o.chatRepo = chatRepo
	o.msgRepo = msgRepo
	o.problemRepo = problemRepo
	o.txr = txr

	for _, opt := range options {
		opt(&o)
	}
	return o
}

func WithLg(opt *zap.Logger) OptOptionsSetter {
	return func(o *Options) {
		o.lg = opt
	}
}

func (o *Options) Validate() error {
	errs := new(errors461e464ebed9.ValidationErrors)
	errs.Add(errors461e464ebed9.NewValidationError("chatRepo", _validate_Options_chatRepo(o)))
	errs.Add(errors461e464ebed9.NewValidationError("msgRepo", _validate_Options_msgRepo(o)))
	errs.Add(errors461e464ebed9.NewValidationError("problemRepo", _validate_Options_problemRepo(o)))
	errs.Add(errors461e464ebed9.NewValidationError("txr", _validate_Options_txr(o)))
	return errs.AsError()
}

func _validate_Options_chatRepo(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.chatRepo, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `chatRepo` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_msgRepo(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.msgRepo, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `msgRepo` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_problemRepo(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.problemRepo, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `problemRepo` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_txr(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.txr, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `txr` did not pass the test: %w", err)
	}
	return nil
}