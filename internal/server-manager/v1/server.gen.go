// Package managerv1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.5-0.20230506011706-29ebe3262399 DO NOT EDIT.
package managerv1

import (
	"fmt"
	"net/http"
	"time"

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
	ErrorCodeProblemNotExistError   ErrorCode = 5001
)

// Chat defines model for Chat.
type Chat struct {
	ChatId   types.ChatID `json:"chatId"`
	ClientId types.UserID `json:"clientId"`
}

// ChatId defines model for ChatId.
type ChatId struct {
	ChatId types.ChatID `json:"chatId"`
}

// ChatList defines model for ChatList.
type ChatList struct {
	Chats []Chat `json:"chats"`
}

// CloseChatResponse defines model for CloseChatResponse.
type CloseChatResponse struct {
	Data  *map[string]interface{} `json:"data,omitempty"`
	Error *Error                  `json:"error,omitempty"`
}

// Error defines model for Error.
type Error struct {
	// Code contains HTTP error codes and specific business logic error codes (the last must be >= 1000).
	Code    ErrorCode `json:"code"`
	Details *string   `json:"details,omitempty"`
	Message string    `json:"message"`
}

// ErrorCode contains HTTP error codes and specific business logic error codes (the last must be >= 1000).
type ErrorCode int

// FreeHandsResponse defines model for FreeHandsResponse.
type FreeHandsResponse struct {
	Data  *string `json:"data"`
	Error *Error  `json:"error,omitempty"`
}

// GetChatHistoryRequest defines model for GetChatHistoryRequest.
type GetChatHistoryRequest struct {
	ChatId   types.ChatID `json:"chatId"`
	Cursor   *string      `json:"cursor,omitempty"`
	PageSize *int         `json:"pageSize,omitempty"`
}

// GetChatHistoryResponse defines model for GetChatHistoryResponse.
type GetChatHistoryResponse struct {
	Data  *MessagesPage `json:"data,omitempty"`
	Error *Error        `json:"error,omitempty"`
}

// GetChatsResponse defines model for GetChatsResponse.
type GetChatsResponse struct {
	Data  *ChatList `json:"data,omitempty"`
	Error *Error    `json:"error,omitempty"`
}

// GetFreeHandsBtnAvailabilityResponse defines model for GetFreeHandsBtnAvailabilityResponse.
type GetFreeHandsBtnAvailabilityResponse struct {
	Data  *ManagerAvailability `json:"data,omitempty"`
	Error *Error               `json:"error,omitempty"`
}

// ManagerAvailability defines model for ManagerAvailability.
type ManagerAvailability struct {
	Available *bool `json:"available,omitempty"`
}

// Message defines model for Message.
type Message struct {
	AuthorId  types.UserID    `json:"authorId"`
	Body      string          `json:"body"`
	CreatedAt time.Time       `json:"createdAt"`
	Id        types.MessageID `json:"id"`
}

// MessageWithoutBody defines model for MessageWithoutBody.
type MessageWithoutBody struct {
	AuthorId  types.UserID    `json:"authorId"`
	CreatedAt time.Time       `json:"createdAt"`
	Id        types.MessageID `json:"id"`
}

// MessagesPage defines model for MessagesPage.
type MessagesPage struct {
	Messages []Message `json:"messages"`
	Next     string    `json:"next"`
}

// SendMessageRequest defines model for SendMessageRequest.
type SendMessageRequest struct {
	ChatId      types.ChatID `json:"chatId"`
	MessageBody string       `json:"messageBody"`
}

// SendMessageResponse defines model for SendMessageResponse.
type SendMessageResponse struct {
	Data  *MessageWithoutBody `json:"data,omitempty"`
	Error *Error              `json:"error,omitempty"`
}

// XRequestIDHeader defines model for XRequestIDHeader.
type XRequestIDHeader = types.RequestID

// PostCloseChatParams defines parameters for PostCloseChat.
type PostCloseChatParams struct {
	XRequestID XRequestIDHeader `json:"X-Request-ID"`
}

// PostFreeHandsParams defines parameters for PostFreeHands.
type PostFreeHandsParams struct {
	XRequestID XRequestIDHeader `json:"X-Request-ID"`
}

// PostGetChatHistoryParams defines parameters for PostGetChatHistory.
type PostGetChatHistoryParams struct {
	XRequestID XRequestIDHeader `json:"X-Request-ID"`
}

// PostGetChatsParams defines parameters for PostGetChats.
type PostGetChatsParams struct {
	XRequestID XRequestIDHeader `json:"X-Request-ID"`
}

// PostGetFreeHandsBtnAvailabilityParams defines parameters for PostGetFreeHandsBtnAvailability.
type PostGetFreeHandsBtnAvailabilityParams struct {
	XRequestID XRequestIDHeader `json:"X-Request-ID"`
}

// PostSendMessageParams defines parameters for PostSendMessage.
type PostSendMessageParams struct {
	XRequestID XRequestIDHeader `json:"X-Request-ID"`
}

// PostCloseChatJSONRequestBody defines body for PostCloseChat for application/json ContentType.
type PostCloseChatJSONRequestBody = ChatId

// PostGetChatHistoryJSONRequestBody defines body for PostGetChatHistory for application/json ContentType.
type PostGetChatHistoryJSONRequestBody = GetChatHistoryRequest

// PostSendMessageJSONRequestBody defines body for PostSendMessage for application/json ContentType.
type PostSendMessageJSONRequestBody = SendMessageRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /closeChat)
	PostCloseChat(ctx echo.Context, params PostCloseChatParams) error

	// (POST /freeHands)
	PostFreeHands(ctx echo.Context, params PostFreeHandsParams) error

	// (POST /getChatHistory)
	PostGetChatHistory(ctx echo.Context, params PostGetChatHistoryParams) error

	// (POST /getChats)
	PostGetChats(ctx echo.Context, params PostGetChatsParams) error

	// (POST /getFreeHandsBtnAvailability)
	PostGetFreeHandsBtnAvailability(ctx echo.Context, params PostGetFreeHandsBtnAvailabilityParams) error

	// (POST /sendMessage)
	PostSendMessage(ctx echo.Context, params PostSendMessageParams) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// PostCloseChat converts echo context to params.
func (w *ServerInterfaceWrapper) PostCloseChat(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params PostCloseChatParams

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
	err = w.Handler.PostCloseChat(ctx, params)
	return err
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

// PostGetChatHistory converts echo context to params.
func (w *ServerInterfaceWrapper) PostGetChatHistory(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params PostGetChatHistoryParams

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
	err = w.Handler.PostGetChatHistory(ctx, params)
	return err
}

// PostGetChats converts echo context to params.
func (w *ServerInterfaceWrapper) PostGetChats(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params PostGetChatsParams

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
	err = w.Handler.PostGetChats(ctx, params)
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

// PostSendMessage converts echo context to params.
func (w *ServerInterfaceWrapper) PostSendMessage(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params PostSendMessageParams

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
	err = w.Handler.PostSendMessage(ctx, params)
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

	router.POST(baseURL+"/closeChat", wrapper.PostCloseChat)
	router.POST(baseURL+"/freeHands", wrapper.PostFreeHands)
	router.POST(baseURL+"/getChatHistory", wrapper.PostGetChatHistory)
	router.POST(baseURL+"/getChats", wrapper.PostGetChats)
	router.POST(baseURL+"/getFreeHandsBtnAvailability", wrapper.PostGetFreeHandsBtnAvailability)
	router.POST(baseURL+"/sendMessage", wrapper.PostSendMessage)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var SwaggerSpec = []string{

	"H4sIAAAAAAAC/9xYbW/bNhD+KwS3DxsgR/KyAoWBfUhfk6Fdg6ZDC2T+QEtniS1FqryTm6zwfx9ISrIc",
	"yU6aviDZp0R8vXvuuYd3/sxTU1ZGgybks8+8ElaUQGD917vX8LEGpJMnxyAysG5Maj7jRfiMuBYl8Bl/",
	"N2lWTk6e8Ihb+FhLCxmfka0h4pgWUAq3e2lsKYjPeF3LjEecLiu3H8lKnfOIX0xyM2kG3R886Ezoz05k",
	"WRlLwWIq+Iznkop6cZCaMj6V6QdRI1gtRAlxWgiaINiVTCGWmty4iv3hfL1er1vzvMePCxFOtaYCSxL8",
	"qDviJLud/e7E72V8xFMlQd/atr8R7HcDtk+D8xbCnsXzzkKzeA8p8XXk4Q/O3I8AjDs5bzx5IXEHmfw/",
	"kqD0//xsYcln/Kd4k4txw8nYE3Ld+SmsFZd87F4M1yqD4Pa8BqyMRhjenwnyqTgAH6w19jp7nvpF3vOn",
	"7for/pkMbnTKY7dwHfEMSEiFPZuaaK4jXgKiyGFk7goE7cIo3D9v7XvcWJMBplZWJI3Tr9RoElIjO37z",
	"5pR5x5nbh0zojGEFqVzKlC1qlBoQmTK5TLfW/UIFMCWQWFkjsQWwf+okOYQ/2DRJkl8PeMRB1yWfnT9I",
	"kiR6kCTTecRLqWXpRn9Pki6mjlW5V9OLidszWQnrdBWdX50TL4UWOdhXK7DKiAyyAH/PzVNrFgrKvww9",
	"vZBIYd4B8cwCHAud4fWk0LVSYqGgVe5BQL6UJM+BHB2PJZKxl42Y3yuFrS0GjwdYVCKHM/mvh7MUFyGy",
	"UxfZLs7TYZj3iMZVsK6L1r4YvAwJgacuK24fOPw6KzoZvJ0FHXMfkT5aCanEQipJXwtNSKX+gbexb+yY",
	"gT0izKq+hC2MUSB0c8pG4a7srKkw9i4+7RFfmOxyNCVSC4IgO6ItqzNBMCFZAh9RFHlLDxvgftDT7o3y",
	"bvedjDZRmm9i+VZSYWp61IB0f8L6/45eh3zf0V7YglQOAtZUFjcv2dqMHlRtEddwQTeuZZA3G5yNZ6Cz",
	"5uDeIyqUerXks/PrRfgk4+toh2stUUtx8QJ07oA/TJpnrB2YRjcz2p81LO0HLnyDp62fZl8s367tg7S2",
	"ki7P3Fy4fQHCgj2qncft17OW23++fcObZtGruJ/dkL0gqgLzpF4aH2VJTvf5I6E/sLO6ctxmLhiseTrY",
	"0ekJj/gKLIa6dDV1npgKtKgkn/HDg+TgkEc+G7yBcdoW+B46E2iwXdy6HkDqnLkMcXWog1e4OSc6/NQg",
	"dV2CP3rT8e/g0WZJPPhFYD0PNADs9M4V16ADP6tKydRfHr9HZ93n3o8BN+HsNs1cUeoHAn08JL8lybe7",
	"ddA+eQO28W3oxxoROWjoFC/bYmV3aM6AkIlewcCWxjINn1hoi3E8Xl0Z9K3i9Z3gGzYaI/BdiwGzgLWi",
	"Dtd8qyreDe5zIM95VoSV41hu19h3NwHGG6cfnA87GpLdSYFMSRyEDvcHzekdq0L3ikwgylxDxsiwMujk",
	"vkDe9ZwY9FEj4LkOiZmlZy+yT5KKDo4ekruaoT0vQQHpB2ROmFjhdm7l3c782HnRnUf62n5xBPwGpC1F",
	"uqJAuKlc9mm7zryMNZWQ4y8VsOcZ7hVEd1eHRgrPHyxCY3XjTZ/lXpnnUe0XeOdzh5lrIVrMtw98AitQ",
	"pipBEwureMRrq5pabxbHyqRCFQZp9jB5OI1d9TZf/xcAAP//cvVt7FMZAAA=",
}

func GetSwagger() (swagger *openapi3.T, err error) {
	return server.GetSwagger(SwaggerSpec)
}
