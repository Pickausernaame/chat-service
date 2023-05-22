package websocketstream

import (
	"io"
	"net/http"
	"time"

	gorillaws "github.com/gorilla/websocket"
)

type Websocket interface {
	SetWriteDeadline(t time.Time) error
	NextWriter(messageType int) (io.WriteCloser, error)
	WriteMessage(messageType int, data []byte) error
	WriteControl(messageType int, data []byte, deadline time.Time) error

	SetPongHandler(h func(appData string) error)
	SetReadDeadline(t time.Time) error
	NextReader() (messageType int, r io.Reader, err error)

	Close() error
}

type Upgrader interface {
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (Websocket, error)
}

type upgraderImpl struct {
	upgrader *gorillaws.Upgrader
}

func NewUpgrader(allowOrigins []string, secWsProtocol string) Upgrader {
	upgrader := &gorillaws.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			for _, allowedOrigin := range allowOrigins {
				if origin == allowedOrigin {
					return true
				}
			}
			return false
		},
		HandshakeTimeout: time.Second * 30,
		//ReadBufferSize:  3000 * 4 * 2,
		//WriteBufferSize: 3000 * 4 * 2,
		Subprotocols: []string{secWsProtocol},
	}
	return &upgraderImpl{
		upgrader: upgrader,
	}
}

func (u *upgraderImpl) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (Websocket, error) {
	return u.upgrader.Upgrade(w, r, responseHeader)
}
