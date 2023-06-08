package websocketstream

import (
	"sync"
	"time"

	gorillaws "github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	closeDeadline = 5 * time.Second
	graceTimeout  = 1 * time.Second
)

type wsCloser struct {
	once   sync.Once
	logger *zap.Logger
	ws     Websocket
}

func newWsCloser(ws Websocket) *wsCloser {
	return &wsCloser{
		ws:     ws,
		logger: zap.L().Named("ws-closer"),
		once:   sync.Once{},
	}
}

func (c *wsCloser) Close(code int) {
	c.once.Do(func() {
		c.logger.Debug("close connection")

		_ = c.ws.WriteControl(
			gorillaws.CloseMessage,
			gorillaws.FormatCloseMessage(code, ""),
			time.Now().Add(closeDeadline),
		)

		time.Sleep(graceTimeout)
		_ = c.ws.Close()
	})
}
