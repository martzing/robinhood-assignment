package validate_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"robinhood-assignment/helpers"
	"robinhood-assignment/internal/core/ports"
	"robinhood-assignment/internal/dto"
	"robinhood-assignment/internal/validate"
	"testing"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type testAuthValidate struct {
	authValidate ports.AuthValidate
}

func newTestAuthValidate(t *testing.T) testAuthValidate {
	authValidate := validate.NewAuthValidate()
	return testAuthValidate{authValidate}
}

func TestValidateLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	govalidator.SetFieldsRequiredByDefault(true)
	type requestBody struct {
		Username string
		Password string
	}
	t.Run("validate login success", func(t *testing.T) {
		username := "samart"
		password := "1234567890"
		body := requestBody{
			Username: username,
			Password: password,
		}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request, _ = http.NewRequest("POST", "http://example.com", &buf)
		tvalid := newTestAuthValidate(t)
		got, err := tvalid.authValidate.ValidateLogin(ctx)
		expected := &dto.LoginRequest{
			Username: username,
			Password: password,
		}
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})
	t.Run("validate login error when username is missing", func(t *testing.T) {
		body := requestBody{
			Password: "1234567890",
		}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request, _ = http.NewRequest("POST", "http://example.com", &buf)
		tvalid := newTestAuthValidate(t)
		got, err := tvalid.authValidate.ValidateLogin(ctx)
		expected := helpers.NewCustomError(http.StatusBadRequest, "username: Missing required field")
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
	t.Run("validate login error when password is missing", func(t *testing.T) {
		body := requestBody{
			Username: "samart",
		}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request, _ = http.NewRequest("POST", "http://example.com", &buf)
		tvalid := newTestAuthValidate(t)
		got, err := tvalid.authValidate.ValidateLogin(ctx)
		expected := helpers.NewCustomError(http.StatusBadRequest, "password: Missing required field")
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
}

func TestValidateCreateStaff(t *testing.T) {
	gin.SetMode(gin.TestMode)
	govalidator.SetFieldsRequiredByDefault(true)
	type requestBody struct {
		Name     string
		Email    string
		Username string
		Password string
		ImageUrl string
		Role     string
	}
	name := "Samart"
	email := "samart@gmail.com"
	username := "samart"
	password := "1234567890"
	imageUrl := "https://image-url.com"
	role := "STAFF"
	t.Run("validate create staff success", func(t *testing.T) {
		body := requestBody{
			Name:     name,
			Email:    email,
			Username: username,
			Password: password,
			ImageUrl: imageUrl,
			Role:     role,
		}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request, _ = http.NewRequest("POST", "http://example.com", &buf)
		tvalid := newTestAuthValidate(t)
		got, err := tvalid.authValidate.ValidateCreateStaff(ctx)
		expected := &dto.CreateStaffRequest{
			Name:     name,
			Email:    email,
			Username: username,
			Password: password,
			ImageUrl: imageUrl,
			Role:     role,
		}
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})
	t.Run("validate create error when name is missing", func(t *testing.T) {
		body := requestBody{
			Email:    email,
			Username: username,
			Password: password,
			ImageUrl: imageUrl,
			Role:     role,
		}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request, _ = http.NewRequest("POST", "http://example.com", &buf)
		tvalid := newTestAuthValidate(t)
		got, err := tvalid.authValidate.ValidateCreateStaff(ctx)
		expected := helpers.NewCustomError(http.StatusBadRequest, "name: Missing required field")
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
	t.Run("validate create error when email is missing", func(t *testing.T) {
		body := requestBody{
			Name:     name,
			Username: username,
			Password: password,
			ImageUrl: imageUrl,
			Role:     role,
		}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request, _ = http.NewRequest("POST", "http://example.com", &buf)
		tvalid := newTestAuthValidate(t)
		got, err := tvalid.authValidate.ValidateCreateStaff(ctx)
		expected := helpers.NewCustomError(http.StatusBadRequest, "email: Missing required field")
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
	t.Run("validate create error when username is missing", func(t *testing.T) {
		body := requestBody{
			Name:     name,
			Email:    email,
			Password: password,
			ImageUrl: imageUrl,
			Role:     role,
		}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request, _ = http.NewRequest("POST", "http://example.com", &buf)
		tvalid := newTestAuthValidate(t)
		got, err := tvalid.authValidate.ValidateCreateStaff(ctx)
		expected := helpers.NewCustomError(http.StatusBadRequest, "username: Missing required field")
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
	t.Run("validate create error when password is missing", func(t *testing.T) {
		body := requestBody{
			Name:     name,
			Email:    email,
			Username: username,
			ImageUrl: imageUrl,
			Role:     role,
		}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request, _ = http.NewRequest("POST", "http://example.com", &buf)
		tvalid := newTestAuthValidate(t)
		got, err := tvalid.authValidate.ValidateCreateStaff(ctx)
		expected := helpers.NewCustomError(http.StatusBadRequest, "password: Missing required field")
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
	t.Run("validate create error when image url is missing", func(t *testing.T) {
		body := requestBody{
			Name:     name,
			Email:    email,
			Username: username,
			Password: password,
			Role:     role,
		}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request, _ = http.NewRequest("POST", "http://example.com", &buf)
		tvalid := newTestAuthValidate(t)
		got, err := tvalid.authValidate.ValidateCreateStaff(ctx)
		expected := helpers.NewCustomError(http.StatusBadRequest, "imageUrl: Missing required field")
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
	t.Run("validate create error when role is missing", func(t *testing.T) {
		body := requestBody{
			Name:     name,
			Email:    email,
			Username: username,
			Password: password,
			ImageUrl: imageUrl,
		}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request, _ = http.NewRequest("POST", "http://example.com", &buf)
		tvalid := newTestAuthValidate(t)
		got, err := tvalid.authValidate.ValidateCreateStaff(ctx)
		expected := helpers.NewCustomError(http.StatusBadRequest, "role: Missing required field")
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
	t.Run("validate create error when email is invalid format", func(t *testing.T) {
		body := requestBody{
			Name:     name,
			Email:    "samart.com",
			Username: username,
			Password: password,
			ImageUrl: imageUrl,
			Role:     role,
		}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request, _ = http.NewRequest("POST", "http://example.com", &buf)
		tvalid := newTestAuthValidate(t)
		got, err := tvalid.authValidate.ValidateCreateStaff(ctx)
		expected := helpers.NewCustomError(http.StatusBadRequest, "email: samart.com does not validate as email")
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
	t.Run("validate create error when image url is invalid format", func(t *testing.T) {
		body := requestBody{
			Name:     name,
			Email:    email,
			Username: username,
			Password: password,
			ImageUrl: "image-url",
			Role:     role,
		}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request, _ = http.NewRequest("POST", "http://example.com", &buf)
		tvalid := newTestAuthValidate(t)
		got, err := tvalid.authValidate.ValidateCreateStaff(ctx)
		expected := helpers.NewCustomError(http.StatusBadRequest, "imageUrl: image-url does not validate as url")
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
	t.Run("validate create error when role not in ecepted list", func(t *testing.T) {
		body := requestBody{
			Name:     name,
			Email:    email,
			Username: username,
			Password: password,
			ImageUrl: imageUrl,
			Role:     "user",
		}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request, _ = http.NewRequest("POST", "http://example.com", &buf)
		tvalid := newTestAuthValidate(t)
		got, err := tvalid.authValidate.ValidateCreateStaff(ctx)
		expected := helpers.NewCustomError(http.StatusBadRequest, "role: user does not validate as in(STAFF|ADMIN)")
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
}
