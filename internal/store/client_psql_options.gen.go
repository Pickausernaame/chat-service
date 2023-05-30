// Code generated by options-gen. DO NOT EDIT.
package store

import (
	fmt461e464ebed9 "fmt"

	errors461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/errors"
	validator461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/validator"
)

type OptPSQLOptionsSetter func(o *PSQLOptions)

func NewPSQLOptions(
	address string,
	username string,
	password string,
	database string,
	options ...OptPSQLOptionsSetter,
) PSQLOptions {
	o := PSQLOptions{}

	// Setting defaults from field tag (if present)

	o.address = address
	o.username = username
	o.password = password
	o.database = database

	for _, opt := range options {
		opt(&o)
	}
	return o
}

func WithDebug(opt bool) OptPSQLOptionsSetter {
	return func(o *PSQLOptions) {
		o.debug = opt
	}
}

func (o *PSQLOptions) Validate() error {
	errs := new(errors461e464ebed9.ValidationErrors)
	errs.Add(errors461e464ebed9.NewValidationError("address", _validate_PSQLOptions_address(o)))
	errs.Add(errors461e464ebed9.NewValidationError("username", _validate_PSQLOptions_username(o)))
	errs.Add(errors461e464ebed9.NewValidationError("password", _validate_PSQLOptions_password(o)))
	errs.Add(errors461e464ebed9.NewValidationError("database", _validate_PSQLOptions_database(o)))
	return errs.AsError()
}

func _validate_PSQLOptions_address(o *PSQLOptions) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.address, "required,hostname_port"); err != nil {
		return fmt461e464ebed9.Errorf("field `address` did not pass the test: %w", err)
	}
	return nil
}

func _validate_PSQLOptions_username(o *PSQLOptions) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.username, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `username` did not pass the test: %w", err)
	}
	return nil
}

func _validate_PSQLOptions_password(o *PSQLOptions) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.password, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `password` did not pass the test: %w", err)
	}
	return nil
}

func _validate_PSQLOptions_database(o *PSQLOptions) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.database, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `database` did not pass the test: %w", err)
	}
	return nil
}
