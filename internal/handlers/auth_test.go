package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"robinhood-assignment/helpers"
	"robinhood-assignment/internal/core/ports"
	"robinhood-assignment/internal/core/ports/mocks"
	"robinhood-assignment/internal/dto"
	"robinhood-assignment/internal/handlers"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type testAuthHandler struct {
	authService  *mocks.AuthServie
	authValidate *mocks.AuthValidate
	handler      ports.AuthHandler
}

func newTestAuthHandler(t *testing.T) testAuthHandler {
	authService := mocks.NewAuthServie(t)
	authValidate := mocks.NewAuthValidate(t)
	handler := handlers.NewAuthHandler(authService, authValidate)
	return testAuthHandler{authService, authValidate, handler}
}

func TestCreateStaff(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("create staff success", func(t *testing.T) {
		req := dto.CreateStaffRequest{
			Name:     "User name 1",
			Email:    "User email 1",
			Username: "Username 1",
			Password: "Password 1",
			ImageUrl: "https://image-url.com",
			Role:     "STAFF",
		}
		res := dto.BaseResponse{
			StatusCode: http.StatusOK,
			Message:    "success",
		}
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		thld := newTestAuthHandler(t)
		thld.authValidate.On("ValidateCreateStaff", ctx).Return(&req, nil)
		thld.authService.On("CreateStaff", ctx, &req).Return(nil)
		thld.handler.CreateStaff(ctx)
		expected, _ := json.Marshal(res)
		got := w.Body.Bytes()
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expected, got)
	})
	t.Run("create staff error when validate fail", func(t *testing.T) {
		errMsg := "Invalid input parameter"
		res := &dto.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      errMsg,
		}
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		thld := newTestAuthHandler(t)
		thld.authValidate.On("ValidateCreateStaff", ctx).Return(nil, helpers.NewCustomError(http.StatusBadRequest, errMsg))
		thld.handler.CreateStaff(ctx)
		expected, _ := json.Marshal(res)
		got := w.Body.Bytes()
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expected, got)
	})
	t.Run("create staff error when call service fail", func(t *testing.T) {
		req := dto.CreateStaffRequest{
			Name:     "User name 1",
			Email:    "User email 1",
			Username: "Username 1",
			Password: "Password 1",
			ImageUrl: "https://image-url.com",
			Role:     "STAFF",
		}
		errMsg := "Duplicate username"
		res := &dto.ErrorResponse{
			StatusCode: http.StatusConflict,
			Error:      errMsg,
		}
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		thld := newTestAuthHandler(t)
		thld.authValidate.On("ValidateCreateStaff", ctx).Return(&req, nil)
		thld.authService.On("CreateStaff", ctx, &req).Return(helpers.NewCustomError(http.StatusConflict, "Duplicate username"))
		thld.handler.CreateStaff(ctx)
		expected, _ := json.Marshal(res)
		got := w.Body.Bytes()
		assert.Equal(t, http.StatusConflict, w.Code)
		assert.Equal(t, expected, got)
	})
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("login success", func(t *testing.T) {
		req := dto.LoginRequest{
			Username: "Username",
			Password: "Password",
		}
		token := "jwt-token"
		res := dto.LoginResponse{
			StatusCode: http.StatusOK,
			Token:      token,
		}
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		thld := newTestAuthHandler(t)
		thld.authValidate.On("ValidateLogin", ctx).Return(&req, nil)
		thld.authService.On("Login", ctx, &req).Return(token, nil)
		thld.handler.Login(ctx)
		expected, _ := json.Marshal(res)
		got := w.Body.Bytes()
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expected, got)
	})
	t.Run("login error when validate fail", func(t *testing.T) {
		errMsg := "Invalid input parameter"
		res := &dto.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      errMsg,
		}
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		thld := newTestAuthHandler(t)
		thld.authValidate.On("ValidateLogin", ctx).Return(nil, helpers.NewCustomError(http.StatusBadRequest, errMsg))
		thld.handler.Login(ctx)
		expected, _ := json.Marshal(res)
		got := w.Body.Bytes()
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expected, got)
	})
	t.Run("login error when call service fail", func(t *testing.T) {
		req := dto.LoginRequest{
			Username: "Username",
			Password: "Password",
		}
		errMsg := "Password is incorrect"
		res := &dto.ErrorResponse{
			StatusCode: http.StatusUnauthorized,
			Error:      errMsg,
		}
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		thld := newTestAuthHandler(t)
		thld.authValidate.On("ValidateLogin", ctx).Return(&req, nil)
		thld.authService.On("Login", ctx, &req).Return("", helpers.NewCustomError(http.StatusUnauthorized, errMsg))
		thld.handler.Login(ctx)
		expected, _ := json.Marshal(res)
		got := w.Body.Bytes()
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, expected, got)
	})
}
