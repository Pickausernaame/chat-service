package managerv1_test

import (
	"errors"
	"net/http"

	srverr "github.com/Pickausernaame/chat-service/internal/errors"
	managerv1 "github.com/Pickausernaame/chat-service/internal/server-manager/v1"
	"github.com/Pickausernaame/chat-service/internal/types"
	canreceiveproblems "github.com/Pickausernaame/chat-service/internal/usecases/manager/can-receive-problems"
	setreadyreceiveproblems "github.com/Pickausernaame/chat-service/internal/usecases/manager/set-ready-receive-problems"
)

func (s *HandlersSuite) TestGetFreeHandsBtnAvailability_Usecase_Error() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/getFreeHandsBtnAvailability", "")
	s.canReceiveProblemsUseCase.EXPECT().Handle(eCtx.Request().Context(), canreceiveproblems.Request{
		ID:        reqID,
		ManagerID: s.managerID,
	}).Return(canreceiveproblems.Response{}, errors.New("something went wrong"))

	// Action.
	err := s.handlers.PostGetFreeHandsBtnAvailability(eCtx, managerv1.PostGetFreeHandsBtnAvailabilityParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestGetFreeHandsBtnAvailability_Usecase_Success() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/getFreeHandsBtnAvailability", "")
	s.canReceiveProblemsUseCase.EXPECT().Handle(eCtx.Request().Context(), canreceiveproblems.Request{
		ID:        reqID,
		ManagerID: s.managerID,
	}).Return(canreceiveproblems.Response{Result: true}, nil)

	// Action.
	err := s.handlers.PostGetFreeHandsBtnAvailability(eCtx, managerv1.PostGetFreeHandsBtnAvailabilityParams{XRequestID: reqID})

	// Assert.
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.Code)
	s.JSONEq(`
{
    "data":
    {
        "available": true
    }
}`, resp.Body.String())
}

func (s *HandlersSuite) TestPostFreeHands_Usecase_OverloadError() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/getFreeHandsBtnAvailability", "")
	s.setReadyReceiveProblems.EXPECT().
		Handle(eCtx.Request().Context(), setreadyreceiveproblems.Request{
			ID:        reqID,
			ManagerID: s.managerID,
		}).
		Return(setreadyreceiveproblems.Response{}, setreadyreceiveproblems.ErrManagerOverload)

	// Action.
	err := s.handlers.PostFreeHands(eCtx, managerv1.PostFreeHandsParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	srvErr := &srverr.ServerError{}
	s.ErrorAs(err, &srvErr)
	s.Require().Equal(5000, srvErr.Code)
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestPostFreeHands_Usecase_SomeError() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/getFreeHandsBtnAvailability", "")
	s.setReadyReceiveProblems.EXPECT().
		Handle(eCtx.Request().Context(), setreadyreceiveproblems.Request{
			ID:        reqID,
			ManagerID: s.managerID,
		}).
		Return(setreadyreceiveproblems.Response{}, errors.New("some error"))

	// Action.
	err := s.handlers.PostFreeHands(eCtx, managerv1.PostFreeHandsParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestPostFreeHands_Usecase_Success() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/getFreeHandsBtnAvailability", "")
	s.setReadyReceiveProblems.EXPECT().
		Handle(eCtx.Request().Context(), setreadyreceiveproblems.Request{
			ID:        reqID,
			ManagerID: s.managerID,
		}).
		Return(setreadyreceiveproblems.Response{}, nil)

	// Action.
	err := s.handlers.PostFreeHands(eCtx, managerv1.PostFreeHandsParams{XRequestID: reqID})

	// Assert.
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.Code)
	s.JSONEq(`
{
    "data": null
}`, resp.Body.String())
}
