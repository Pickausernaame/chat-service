package websocketstream

import (
	"context"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/Pickausernaame/chat-service/internal/middlewares"
	eventstream "github.com/Pickausernaame/chat-service/internal/services/event-stream"
	"github.com/Pickausernaame/chat-service/internal/types"
)

const (
	writeTimeout = time.Second
)

type eventStream interface {
	Subscribe(ctx context.Context, userID types.UserID) (<-chan eventstream.Event, error)
}

//go:generate options-gen -out-filename=handler_options.gen.go -from-struct=Options
type Options struct {
	pingPeriod time.Duration `default:"3s" validate:"omitempty,min=100ms,max=30s"`

	eventStream  eventStream     `option:"mandatory" validate:"required"`
	eventAdapter EventAdapter    `option:"mandatory" validate:"required"`
	eventWriter  EventWriter     `option:"mandatory" validate:"required"`
	upgrader     Upgrader        `option:"mandatory" validate:"required"`
	shutdownCh   <-chan struct{} `option:"mandatory" validate:"required"`
}

type HTTPHandler struct {
	Options
	lg *zap.Logger
}

func NewHTTPHandler(opts Options) (*HTTPHandler, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validating opts: %v", err)
	}

	return &HTTPHandler{Options: opts, lg: zap.L().Named("ws-stream")}, nil
}

func (h *HTTPHandler) Serve(eCtx echo.Context) error {
	ws, err := h.upgrader.Upgrade(eCtx.Response(), eCtx.Request(), nil)
	if err != nil {
		return fmt.Errorf("connecting ws: %v", err)
	}

	clientID := middlewares.MustUserID(eCtx)
	ctx, cancel := context.WithCancel(eCtx.Request().Context())
	defer cancel()

	events, err := h.eventStream.Subscribe(ctx, clientID)
	if err != nil {
		return fmt.Errorf("subscribing ws for clientID %s: %v", clientID.String(), err)
	}

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return h.readLoop(ctx, ws)
	})

	eg.Go(func() error {
		return h.writeLoop(ctx, ws, events)
	})

	select {
	case <-h.shutdownCh:
	case <-ctx.Done():
	}
	newWsCloser(ws).Close(websocket.CloseNormalClosure)
	cancel()
	if err = eg.Wait(); err != nil {
		return fmt.Errorf("waiting err group: %v", err)
	}
	return nil
}

// readLoop listen PONGs.
func (h *HTTPHandler) readLoop(_ context.Context, ws Websocket) error {
	ws.SetPongHandler(func(string) error {
		return ws.SetReadDeadline(time.Now().Add(h.pingPeriod))
	})
	for {
		_, _, err := ws.NextReader()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				h.lg.Error("reading error", zap.Error(err))
			}
			break
		}
	}

	return nil
}

// writeLoop listen events and writes them into Websocket.
func (h *HTTPHandler) writeLoop(ctx context.Context, ws Websocket, events <-chan eventstream.Event) error {
	ticker := time.NewTicker(h.pingPeriod)
	for {
		select {
		case <-ctx.Done():
			return nil
		case e := <-events:
			err := ws.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err != nil {
				return fmt.Errorf("setting write deadline for event: %v", err)
			}

			wr, err := ws.NextWriter(websocket.TextMessage)
			if err != nil {
				return fmt.Errorf("getting next writer: %v", err)
			}

			if err = h.eventWriter.Write(e, wr); err != nil {
				return fmt.Errorf("writing message: %v", err)
			}
		case <-ticker.C:
			err := h.Ping(ws)
			if err != nil {
				return err
			}
		}
	}
}

func (h *HTTPHandler) Ping(ws Websocket) error {
	err := ws.SetWriteDeadline(time.Now().Add(writeTimeout))
	if err != nil {
		return fmt.Errorf("setting write deadline for event: %v", err)
	}

	if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
		h.lg.Error("ping writing", zap.Error(err))
		return fmt.Errorf("writing ping: %v", err)
	}
	return nil
}
