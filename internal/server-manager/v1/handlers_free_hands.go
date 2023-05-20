package managerv1

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Pickausernaame/chat-service/internal/middlewares"
	canreceiveproblems "github.com/Pickausernaame/chat-service/internal/usecases/manager/can-receive-problems"
	"github.com/Pickausernaame/chat-service/pkg/pointer"
)

func (h Handlers) PostGetFreeHandsBtnAvailability(eCtx echo.Context, params PostGetFreeHandsBtnAvailabilityParams) error {
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
		PostGetFreeHandsBtnAvailabilityResponse{
			Data: &ManagerAvailability{
				Available: pointer.PtrWithZeroAsNil(resp.Result),
			},
		})
}
