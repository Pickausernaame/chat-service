package errhandler

import (
	clientv1 "github.com/Pickausernaame/chat-service/internal/server-client/v1"
)

type Response struct {
	Error clientv1.Error `json:"error"`
}

var ResponseBuilder = func(code int, msg string, details string) any {
	resp := Response{Error: clientv1.Error{Code: code, Message: msg}}
	if details != "" {
		resp.Error.Details = &details
	}
	return resp
}
