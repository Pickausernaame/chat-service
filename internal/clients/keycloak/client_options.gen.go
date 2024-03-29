// Code generated by options-gen. DO NOT EDIT.
package keycloakclient

import (
	fmt461e464ebed9 "fmt"

	errors461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/errors"
	validator461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/validator"
)

type OptOptionsSetter func(o *Options)

func NewOptions(
	basePath string,
	realmName string,
	clientID string,
	clientSecret string,
	options ...OptOptionsSetter,
) Options {
	o := Options{}

	// Setting defaults from field tag (if present)

	o.basePath = basePath
	o.realmName = realmName
	o.clientID = clientID
	o.clientSecret = clientSecret

	for _, opt := range options {
		opt(&o)
	}
	return o
}

func WithDebugMode(opt bool) OptOptionsSetter {
	return func(o *Options) {
		o.debugMode = opt
	}
}

func (o *Options) Validate() error {
	errs := new(errors461e464ebed9.ValidationErrors)
	errs.Add(errors461e464ebed9.NewValidationError("basePath", _validate_Options_basePath(o)))
	errs.Add(errors461e464ebed9.NewValidationError("realmName", _validate_Options_realmName(o)))
	errs.Add(errors461e464ebed9.NewValidationError("clientID", _validate_Options_clientID(o)))
	errs.Add(errors461e464ebed9.NewValidationError("clientSecret", _validate_Options_clientSecret(o)))
	return errs.AsError()
}

func _validate_Options_basePath(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.basePath, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `basePath` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_realmName(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.realmName, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `realmName` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_clientID(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.clientID, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `clientID` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_clientSecret(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.clientSecret, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `clientSecret` did not pass the test: %w", err)
	}
	return nil
}
