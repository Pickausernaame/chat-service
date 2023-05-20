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

// Error defines model for Error.
type Error struct {
	// Code contains HTTP error codes and specific business logic error codes (the last must be >= 1000).
	Code    ErrorCode `json:"code"`
	Details *string   `json:"details,omitempty"`
	Message string    `json:"message"`
}

// ErrorCode contains HTTP error codes and specific business logic error codes (the last must be >= 1000).
type ErrorCode = int

// ManagerAvailability defines model for ManagerAvailability.
type ManagerAvailability struct {
	Available *bool `json:"available,omitempty"`
}

// PostGetFreeHandsBtnAvailabilityResponse defines model for PostGetFreeHandsBtnAvailabilityResponse.
type PostGetFreeHandsBtnAvailabilityResponse struct {
	Data  *ManagerAvailability `json:"data,omitempty"`
	Error *Error               `json:"error,omitempty"`
}

// XRequestIDHeader defines model for XRequestIDHeader.
type XRequestIDHeader = types.RequestID

// PostGetFreeHandsBtnAvailabilityParams defines parameters for PostGetFreeHandsBtnAvailability.
type PostGetFreeHandsBtnAvailabilityParams struct {
	XRequestID XRequestIDHeader `json:"X-Request-ID"`
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /getFreeHandsBtnAvailability)
	PostGetFreeHandsBtnAvailability(ctx echo.Context, params PostGetFreeHandsBtnAvailabilityParams) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
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

	router.POST(baseURL+"/getFreeHandsBtnAvailability", wrapper.PostGetFreeHandsBtnAvailability)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var SwaggerSpec = []string{

	"H4sIAAAAAAAC/4xUzW7bOBB+FWJ2D7uAbCmbPQQCeshP07hAASMJ0ACpD2NqLLGRSIYcGg0MvXtByo4V",
	"2G16IsiZIef74WxAms4aTZo9lBuw6LAjJpd2D7f0HMjz7OqGsCIXz5SGEpphm4HGjqCEh8k2czK7ggwc",
	"PQflqIKSXaAMvGyow1i9Mq5DhhJCUBVkwC821nt2SteQwY9JbSbbw7j46WsL4+hEddY4HjrmBkqoFTdh",
	"OZWmy+dKPmHw5DRiR7lskCee3FpJypXmeN7m6XLo+77ftZcQf3TOJJjWGUuOFaVjaSqK69+OVlDCX/me",
	"tXxbnafSy5jYZ1ARo2pT7VuIfQYdeY81HYn1Y+oeXxOz4f1Fn8H+kXIDFXnplGVloibSaEalvbi5v58L",
	"ioki1nmBuhLeklQrJcUyeKXJe9GaWsk3ef9wQ6JFz6ILnsWSxLdQFKf0QZwURfHvFDLolFZd6KD8vyhe",
	"1Yuk1uQiti+osSZ3vkbV4lK1il8O2cQh2o4pWBrTEmpIgsyN50/E147oBnXlL1iPb7wlb432dHhzhYzv",
	"6XSsxT4D2in/rsY715AMTvHLXYwNzy8JHbnzEB25213vHP/56z1svZYAp+j+CzTMdkCv9MokZhRHiuAC",
	"9ZO4CzY6Xlw2yGILQZzPZ5DBmpwfLLA+iUiMJY1WQQmn02J6Cln6I6nBvP41q4lN4/nQWZcNySexckSi",
	"iYUCR2XRFVECjMmzCsr31Ev97KfM43HG9yn5wRTqF/GfDB5IsP4riuGXaiadAKC1rZKpp/y7jyg2oyn0",
	"O4n/1HtJqmM8jckRjnxoeZosM/JMgj12y+MigopTakfK26uvaE2tsR1pFkMWZBBcuzVOmeetkdg2xnN5",
	"Vpyd5NEKi/5nAAAA//+qCrSO3wUAAA==",
}

func GetSwagger() (swagger *openapi3.T, err error) {
	return server.GetSwagger(SwaggerSpec)
}
