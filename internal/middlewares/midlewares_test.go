package middlewares_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"robinhood-assignment/config"
	"robinhood-assignment/internal/core/domains"
	"robinhood-assignment/internal/core/ports"
	"robinhood-assignment/internal/core/ports/mocks"
	"robinhood-assignment/internal/dto"
	"robinhood-assignment/internal/middlewares"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testMiddlewares struct {
	myJWT      *mocks.MyJWT
	middleware ports.Middlewares
}

func newMiddlewares(t *testing.T) testMiddlewares {
	myJWT := mocks.NewMyJWT(t)
	middleware := middlewares.NewMidlewares(myJWT)
	return testMiddlewares{myJWT, middleware}
}

func TestAdminMiddleware(t *testing.T) {
	t.Setenv("JWT_SECRET", "mock-jwt-secret")
	config.New()
	gin.SetMode(gin.TestMode)
	mockJWT := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	t.Run("Authorization is missing", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = &http.Request{
			Header: make(http.Header),
		}
		res := &dto.ErrorResponse{
			StatusCode: http.StatusUnauthorized,
			Error:      "Authorization is missing",
		}
		tmid := newMiddlewares(t)
		tmid.middleware.AdminMiddleware(ctx)
		expected, _ := json.Marshal(res)
		got := w.Body.Bytes()
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, expected, got)
	})

	t.Run("Token expire or invalid token", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = &http.Request{
			Header: make(http.Header),
		}
		ctx.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", mockJWT))
		tmid := newMiddlewares(t)
		claims := &domains.Claims{}
		tmid.myJWT.On("ParseWithClaims", mockJWT, claims, mock.Anything).Return(nil, errors.New("Some error"))
		tmid.middleware.AdminMiddleware(ctx)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Don't have permission", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = &http.Request{
			Header: make(http.Header),
		}
		ctx.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", mockJWT))
		res := &dto.ErrorResponse{
			StatusCode: http.StatusForbidden,
			Error:      "You don't have permission for this API",
		}
		tmid := newMiddlewares(t)
		claims := &domains.Claims{}
		tmid.myJWT.On("ParseWithClaims", mockJWT, claims, mock.Anything).Return(&jwt.Token{}, nil)
		tmid.middleware.AdminMiddleware(ctx)
		expected, _ := json.Marshal(res)
		got := w.Body.Bytes()
		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Equal(t, expected, got)
	})
}

func TestStaffMiddleware(t *testing.T) {
	t.Setenv("JWT_SECRET", "mock-jwt-secret")
	config.New()
	gin.SetMode(gin.TestMode)
	mockJWT := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	t.Run("Authorization is missing", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = &http.Request{
			Header: make(http.Header),
		}
		res := &dto.ErrorResponse{
			StatusCode: http.StatusUnauthorized,
			Error:      "Authorization is missing",
		}
		tmid := newMiddlewares(t)
		tmid.middleware.StaffMiddleware(ctx)
		expected, _ := json.Marshal(res)
		got := w.Body.Bytes()
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, expected, got)
	})

	t.Run("Token expire or invalid token", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = &http.Request{
			Header: make(http.Header),
		}
		ctx.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", mockJWT))
		tmid := newMiddlewares(t)
		claims := &domains.Claims{}
		tmid.myJWT.On("ParseWithClaims", mockJWT, claims, mock.Anything).Return(nil, errors.New("Some error"))
		tmid.middleware.StaffMiddleware(ctx)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Don't have permission", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = &http.Request{
			Header: make(http.Header),
		}
		ctx.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", mockJWT))
		res := &dto.ErrorResponse{
			StatusCode: http.StatusForbidden,
			Error:      "You don't have permission for this API",
		}
		tmid := newMiddlewares(t)
		claims := &domains.Claims{}
		tmid.myJWT.On("ParseWithClaims", mockJWT, claims, mock.Anything).Return(&jwt.Token{}, nil)
		tmid.middleware.StaffMiddleware(ctx)
		expected, _ := json.Marshal(res)
		got := w.Body.Bytes()
		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Equal(t, expected, got)
	})
}
