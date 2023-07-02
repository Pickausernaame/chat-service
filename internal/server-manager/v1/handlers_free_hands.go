package managerv1

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	srverr "github.com/Pickausernaame/chat-service/internal/errors"
	"github.com/Pickausernaame/chat-service/internal/middlewares"
	canreceiveproblems "github.com/Pickausernaame/chat-service/internal/usecases/manager/can-receive-problems"
	setreadyreceiveproblems "github.com/Pickausernaame/chat-service/internal/usecases/manager/set-ready-receive-problems"
	"github.com/Pickausernaame/chat-service/pkg/pointer"
)

func (h Handlers) PostGetFreeHandsBtnAvailability(eCtx echo.Context,
	params PostGetFreeHandsBtnAvailabilityParams,
) error {
	ctx := eCtx.Request().Context()
	managerID := middlewares.MustUserID(eCtx)

	resp, err := h.canReceiveProblemsUseCase.Handle(ctx, canreceiveproblems.Request{
		ID:        params.XRequestID,
		ManagerID: managerID,
	})
	if err != nil {
		return err
	}

	return eCtx.JSON(http.StatusOK,
		GetFreeHandsBtnAvailabilityResponse{
			Data: &ManagerAvailability{
				Available: pointer.PtrWithZeroAsNil(resp.Result),
			},
		})
}

func (h Handlers) PostFreeHands(eCtx echo.Context, params PostFreeHandsParams) error {
	ctx := eCtx.Request().Context()
	managerID := middlewares.MustUserID(eCtx)

	_, err := h.setReadyReceiveProblemsUseCase.Handle(ctx, setreadyreceiveproblems.Request{
		ID:        params.XRequestID,
		ManagerID: managerID,
	})
	if err != nil {
		if errors.Is(err, setreadyreceiveproblems.ErrManagerOverload) {
			return srverr.NewServerError(int(ErrorCodeManagerOverloadedError), "manager overload", err)
		}
		return err
	}

	return eCtx.JSON(http.StatusOK,
		FreeHandsResponse{Data: nil})
}
