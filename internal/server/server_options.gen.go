// Code generated by options-gen. DO NOT EDIT.
package server

import (
	fmt461e464ebed9 "fmt"

	keycloakclient "github.com/Pickausernaame/chat-service/internal/clients/keycloak"
	"github.com/getkin/kin-openapi/openapi3"
	errors461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/errors"
	validator461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/validator"
	"github.com/labstack/echo/v4"
)

type OptOptionsSetter func(o *Options)

func NewOptions(
	serverName string,
	addr string,
	allowOrigins []string,
	v1Swagger *openapi3.T,
	reg func(v *echo.Group),
	keycloakClient *keycloakclient.Client,
	resource string,
	role string,
	secWsProtocol string,
	eventSubscriber eventStream,
	errHandler echo.HTTPErrorHandler,
	adapter adapter,
	options ...OptOptionsSetter,
) Options {
	o := Options{}

	// Setting defaults from field tag (if present)

	o.serverName = serverName
	o.addr = addr
	o.allowOrigins = allowOrigins
	o.v1Swagger = v1Swagger
	o.reg = reg
	o.keycloakClient = keycloakClient
	o.resource = resource
	o.role = role
	o.secWsProtocol = secWsProtocol
	o.eventSubscriber = eventSubscriber
	o.errHandler = errHandler
	o.adapter = adapter

	for _, opt := range options {
		opt(&o)
	}
	return o
}

func (o *Options) Validate() error {
	errs := new(errors461e464ebed9.ValidationErrors)
	errs.Add(errors461e464ebed9.NewValidationError("serverName", _validate_Options_serverName(o)))
	errs.Add(errors461e464ebed9.NewValidationError("addr", _validate_Options_addr(o)))
	errs.Add(errors461e464ebed9.NewValidationError("allowOrigins", _validate_Options_allowOrigins(o)))
	errs.Add(errors461e464ebed9.NewValidationError("v1Swagger", _validate_Options_v1Swagger(o)))
	errs.Add(errors461e464ebed9.NewValidationError("reg", _validate_Options_reg(o)))
	errs.Add(errors461e464ebed9.NewValidationError("keycloakClient", _validate_Options_keycloakClient(o)))
	errs.Add(errors461e464ebed9.NewValidationError("resource", _validate_Options_resource(o)))
	errs.Add(errors461e464ebed9.NewValidationError("role", _validate_Options_role(o)))
	errs.Add(errors461e464ebed9.NewValidationError("secWsProtocol", _validate_Options_secWsProtocol(o)))
	errs.Add(errors461e464ebed9.NewValidationError("eventSubscriber", _validate_Options_eventSubscriber(o)))
	errs.Add(errors461e464ebed9.NewValidationError("errHandler", _validate_Options_errHandler(o)))
	errs.Add(errors461e464ebed9.NewValidationError("adapter", _validate_Options_adapter(o)))
	return errs.AsError()
}

func _validate_Options_serverName(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.serverName, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `serverName` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_addr(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.addr, "required,hostname_port"); err != nil {
		return fmt461e464ebed9.Errorf("field `addr` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_allowOrigins(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.allowOrigins, "min=1"); err != nil {
		return fmt461e464ebed9.Errorf("field `allowOrigins` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_v1Swagger(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.v1Swagger, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `v1Swagger` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_reg(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.reg, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `reg` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_keycloakClient(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.keycloakClient, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `keycloakClient` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_resource(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.resource, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `resource` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_role(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.role, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `role` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_secWsProtocol(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.secWsProtocol, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `secWsProtocol` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_eventSubscriber(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.eventSubscriber, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `eventSubscriber` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_errHandler(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.errHandler, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `errHandler` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_adapter(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.adapter, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `adapter` did not pass the test: %w", err)
	}
	return nil
}
