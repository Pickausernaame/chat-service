// Package clientevents provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.5-0.20230506011706-29ebe3262399 DO NOT EDIT.
package clientevents

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/Pickausernaame/chat-service/internal/types"
	"github.com/getkin/kin-openapi/openapi3"
)

// Defines values for BaseEventEventType.
const (
	BaseEventEventTypeMessageBlockedEvent BaseEventEventType = "MessageBlockedEvent"
	BaseEventEventTypeMessageSentEvent    BaseEventEventType = "MessageSentEvent"
	BaseEventEventTypeNewMessageEvent     BaseEventEventType = "NewMessageEvent"
)

// Defines values for NewMessageEventEventType.
const (
	NewMessageEventEventTypeMessageBlockedEvent NewMessageEventEventType = "MessageBlockedEvent"
	NewMessageEventEventTypeMessageSentEvent    NewMessageEventEventType = "MessageSentEvent"
	NewMessageEventEventTypeNewMessageEvent     NewMessageEventEventType = "NewMessageEvent"
)

// BaseEvent defines model for BaseEvent.
type BaseEvent struct {
	// EventId Unique identifier for the event
	EventId *types.EventID `json:"eventId,omitempty"`

	// EventType Type of the event
	EventType *BaseEventEventType `json:"eventType,omitempty"`

	// MessageId Unique identifier for the message
	MessageId *types.MessageID `json:"messageId,omitempty"`

	// RequestId Unique identifier for the request
	RequestId *types.RequestID `json:"requestId,omitempty"`
}

// BaseEventEventType Type of the event
type BaseEventEventType string

// MessageBlockedEvent defines model for MessageBlockedEvent.
type MessageBlockedEvent = BaseEvent

// MessageSentEvent defines model for MessageSentEvent.
type MessageSentEvent = BaseEvent

// NewMessageEvent defines model for NewMessageEvent.
type NewMessageEvent struct {
	// AuthorId Unique identifier for the author
	AuthorId *types.UserID `json:"authorId,omitempty"`

	// Body Body of the message
	Body *string `json:"body,omitempty"`

	// CreatedAt Date and time of event creation
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	// EventId Unique identifier for the event
	EventId *types.EventID `json:"eventId,omitempty"`

	// EventType Type of the event
	EventType *NewMessageEventEventType `json:"eventType,omitempty"`

	// IsService Indicates if the event is a service event
	IsService *bool `json:"isService,omitempty"`

	// MessageId Unique identifier for the message
	MessageId *types.MessageID `json:"messageId,omitempty"`

	// RequestId Unique identifier for the request
	RequestId *types.RequestID `json:"requestId,omitempty"`
}

// NewMessageEventEventType Type of the event
type NewMessageEventEventType string

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7RUS0/jPBT9K9b9vqXbBLFB3k1hFkjz0jCsEAvXuWnu4Eewbwqo6n8f2QkDE7qhEqvE",
	"x76Pc+5jBya4Pnj0nEDtIJkOnS6/K53w8xY950MfQ4+RCcsVZviyyb8NJhOpZwoeFFx7uh9QUIOeqSWM",
	"og1RcIeimICENkSnGRQMAzUggZ96BAWJI/kNSHhcbMJiAvMnLUsOlxev7xbk+hDHxDR3oGBD3A3rpQmu",
	"+kHmTg8Jo9faYWU6zYuEcUsGK/KccVsV17Dfy5HLrxJwziajIrT/5I9+cKBu4Bs+fMWU9GbSSMJ0vELP",
	"M2hlg7nDZkRv55z3Etz47n2KTkZHaTol9nGqRrwfML2zRyajoxj9nAJ+EKPM6VA11Q60td9bUDc7+D9i",
	"Cwr+q15mqpoGqnqZpv3t/kCzHOto3ofH+JHz+dYDdyG+r3ijzVG1u04YP64V16F5ektkFZqn59l+mSSn",
	"H7+g3+SIp3VdS3Dkn4GTA4NrImrG5hO/DXChGYX2jWByZYuUDSKKRX7xSqhGMy7yMzgQgtLVyO1tiEvf",
	"kNGMSdCrJSUoCS0mQf7urcnvOgSL2hdhJiisf6Mp3ZRB8m3IkZjY5tuV9nfiauhzOcR5p1mcW8pRSu8k",
	"kLDFmMZ8tic54dCj1z2BgtPlybIGWUqYQPnBWgk5MYyp9OdMMdyiDb3L3sdXIGGIFhQ8JFVVNhhtu5BY",
	"ndVndfWQcs5/AgAA//8YHbTTwQYAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
