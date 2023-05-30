// Package managerv1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.5-0.20230506011706-29ebe3262399 DO NOT EDIT.
package managerv1

import (
	"fmt"
	"net/http"

	"github.com/Pickausernaame/chat-service/internal/server"
	"github.com/Pickausernaame/chat-service/internal/types"
	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

const (
	BearerAuthScopes = "bearerAuth.Scopes"
)

// Defines values for ErrorCode.
const (
	ErrorCodeManagerOverloadedError ErrorCode = 5000
)

// Error defines model for Error.
type Error struct {
	// Code contains HTTP error codes and specific business logic error codes (the last must be >= 1000).
	Code    ErrorCode `json:"code"`
	Details *string   `json:"details,omitempty"`
	Message string    `json:"message"`
}

// ErrorCode contains HTTP error codes and specific business logic error codes (the last must be >= 1000).
type ErrorCode int

// ManagerAvailability defines model for ManagerAvailability.
type ManagerAvailability struct {
	Available *bool `json:"available,omitempty"`
}

// PostFreeHandsResponse defines model for PostFreeHandsResponse.
type PostFreeHandsResponse struct {
	Data  *string `json:"data"`
	Error *Error  `json:"error,omitempty"`
}

// PostGetFreeHandsBtnAvailabilityResponse defines model for PostGetFreeHandsBtnAvailabilityResponse.
type PostGetFreeHandsBtnAvailabilityResponse struct {
	Data  *ManagerAvailability `json:"data,omitempty"`
	Error *Error               `json:"error,omitempty"`
}

// XRequestIDHeader defines model for XRequestIDHeader.
type XRequestIDHeader = types.RequestID

// PostFreeHandsParams defines parameters for PostFreeHands.
type PostFreeHandsParams struct {
	XRequestID XRequestIDHeader `json:"X-Request-ID"`
}

// PostGetFreeHandsBtnAvailabilityParams defines parameters for PostGetFreeHandsBtnAvailability.
type PostGetFreeHandsBtnAvailabilityParams struct {
	XRequestID XRequestIDHeader `json:"X-Request-ID"`
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /freeHands)
	PostFreeHands(ctx echo.Context, params PostFreeHandsParams) error

	// (POST /getFreeHandsBtnAvailability)
	PostGetFreeHandsBtnAvailability(ctx echo.Context, params PostGetFreeHandsBtnAvailabilityParams) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// PostFreeHands converts echo context to params.
func (w *ServerInterfaceWrapper) PostFreeHands(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params PostFreeHandsParams

	headers := ctx.Request().Header
	// ------------- Required header parameter "X-Request-ID" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Request-ID")]; found {
		var XRequestID XRequestIDHeader
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for X-Request-ID, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Request-ID", runtime.ParamLocationHeader, valueList[0], &XRequestID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter X-Request-ID: %s", err))
		}

		params.XRequestID = XRequestID
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Header parameter X-Request-ID is required, but not found"))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostFreeHands(ctx, params)
	return err
}

// PostGetFreeHandsBtnAvailability converts echo context to params.
func (w *ServerInterfaceWrapper) PostGetFreeHandsBtnAvailability(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params PostGetFreeHandsBtnAvailabilityParams

	headers := ctx.Request().Header
	// ------------- Required header parameter "X-Request-ID" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Request-ID")]; found {
		var XRequestID XRequestIDHeader
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for X-Request-ID, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Request-ID", runtime.ParamLocationHeader, valueList[0], &XRequestID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter X-Request-ID: %s", err))
		}

		params.XRequestID = XRequestID
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Header parameter X-Request-ID is required, but not found"))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostGetFreeHandsBtnAvailability(ctx, params)
	return err
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router server.EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router server.EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/freeHands", wrapper.PostFreeHands)
	router.POST(baseURL+"/getFreeHandsBtnAvailability", wrapper.PostGetFreeHandsBtnAvailability)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var SwaggerSpec = []string{

	"H4sIAAAAAAAC/9RVX2/jNgz/KgK3hw1wYnfdgIOBPfTa3doBw4prgR3Q5YGRmVirLPkkKrui8HcfKDtN",
	"guR6u8c9BRZJkfz9UZ5B+673jhxHqJ+hx4AdMYX89eE9fUwU+ebqmrChIGfGQQ3t+FmAw46ghg+zKXN2",
	"cwUFBPqYTKAGag6JCoi6pQ6leuVDhww1pGQaKICfeqmPHIxbQwGfZms/mw7lJ85fRtiPzkzX+8DjxNxC",
	"DWvDbVrOte/KW6MfMUUKDrGjUrfIs0hhYzSVxrGc2zJfDsMwDNvx8sa/hODzmn3wPQU2lI+1b0h+vw20",
	"ghq+KXeolVN1mUsvJXEooCFGY3Pt4YpDAR3FiGs6ERv2oXt4SSzG/ouhgF2T+hkaijqYno0XTrR3jMZF",
	"dX1/f6tIEpXURYWuUbEnbVZGq2WKxlGMyvq10Qd533FLymJk1aXIaknqr1RV5/SzOquq6vs5FEAudVA/",
	"/FRV1aKAzjjTycGPVfXCpUC8zuL4NJP02QaDyCTKSi/z/44O1xT+2FCwHhtqRuRlxyl0sUFjcWms4adj",
	"RnCM2n0Yl95bQgeZ1Fsf+V0gukbXxPcUe+8iHd/TIGdhumSn60bJHrFGW2V8UQO7/r/SboS37PY3+vJE",
	"r/U5BdFXzyjKJ52C4ac7iY3tl4SBwkUSV22/3m1d+9uf9zD5JQOeozsbt8z9uL1xK5+ZMSyYwlt0j+ou",
	"9eJaddkiq2kFdXF7AwVsKMRRxpsz2cT35LA3UMP5vJqfQ5F9ngcsV1tIM3Y+8rEX7oijwj101MoH5egf",
	"pa0RNETLAjpKwU0D9aFecr/dS/hwGtFdSnn0Ug4L8fLIcR70h6oaXxLH5PLI2PfW6DxB+XeUuZ/3XsrX",
	"KDyt7Qz8V+KgAsVkeT7poVx/XrGfR/uyJf0YlfCiWqk8aHka6les8T8A/78Y+wQdE1AHhOwTsOfIvPe+",
	"Fx8WspX8j21RObz7ijZkfd+RYzVmQQEp2MmWdVlar9G2PnL9pnpzVorRFsO/AQAA///V5knNAQgAAA==",
}

func GetSwagger() (swagger *openapi3.T, err error) {
	return server.GetSwagger(SwaggerSpec)
}