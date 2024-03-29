// Code generated by options-gen. DO NOT EDIT.
package getassignedproblems

import (
	fmt461e464ebed9 "fmt"

	errors461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/errors"
	validator461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/validator"
)

type OptOptionsSetter func(o *Options)

func NewOptions(
	problemRepo problemRepository,
	chatRepo chatRepository,
	options ...OptOptionsSetter,
) Options {
	o := Options{}

	// Setting defaults from field tag (if present)

	o.problemRepo = problemRepo
	o.chatRepo = chatRepo

	for _, opt := range options {
		opt(&o)
	}
	return o
}

func (o *Options) Validate() error {
	errs := new(errors461e464ebed9.ValidationErrors)
	errs.Add(errors461e464ebed9.NewValidationError("problemRepo", _validate_Options_problemRepo(o)))
	errs.Add(errors461e464ebed9.NewValidationError("chatRepo", _validate_Options_chatRepo(o)))
	return errs.AsError()
}

func _validate_Options_problemRepo(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.problemRepo, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `problemRepo` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_chatRepo(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.chatRepo, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `chatRepo` did not pass the test: %w", err)
	}
	return nil
}
