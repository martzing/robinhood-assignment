package validate_test

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type testInterviewValidate struct {
	interviewValidate ports.InterviewValidate
}

func newTestInterviewValidate(t *testing.T) testInterviewValidate {
	interviewValidate := validate.NewInterviewValidate()
	return testInterviewValidate{interviewValidate}
}

func TestValidateGetInterviewAppointments(t *testing.T) {
	gin.SetMode(gin.TestMode)
	govalidator.SetFieldsRequiredByDefault(true)
	page := uint32(1)
	limit := uint32(10)
	t.Run("validate get interview appointments success", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		url := fmt.Sprintf("http://example.com/?page=%d&limit=%d", page, limit)
		ctx.Request, _ = http.NewRequest("GET", url, nil)
		tvalid := newTestInterviewValidate(t)
		got, err := tvalid.interviewValidate.ValidateGetInterviewAppointments(ctx)
		expected := &dto.GetInterviewAppointmentsRequest{
			Page:  page,
			Limit: limit,
		}
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})
	t.Run("validate get interview appointments error when invalid page params", func(t *testing.T) {
		newPage := "1x"
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		url := fmt.Sprintf("http://example.com/?page=%s&limit=%d", newPage, limit)
		ctx.Request, _ = http.NewRequest("GET", url, nil)
		tvalid := newTestInterviewValidate(t)
		got, err := tvalid.interviewValidate.ValidateGetInterviewAppointments(ctx)
		expected := helpers.NewCustomError(http.StatusBadRequest, "Invalid page query parameter")
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
	t.Run("validate get interview appointments error when invalid page params", func(t *testing.T) {
		newLimit := "1x"
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		url := fmt.Sprintf("http://example.com/?page=%d&limit=%s", page, newLimit)
		ctx.Request, _ = http.NewRequest("GET", url, nil)
		tvalid := newTestInterviewValidate(t)
		got, err := tvalid.interviewValidate.ValidateGetInterviewAppointments(ctx)
		expected := helpers.NewCustomError(http.StatusBadRequest, "Invalid limit query parameter")
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
}

func TestValidateGetInterviewAppointment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	govalidator.SetFieldsRequiredByDefault(true)
	t.Run("validate get interview appointment success", func(t *testing.T) {
		id := "6476f457e64589e868aac97b"
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = []gin.Param{
			{Key: "id", Value: id},
		}
		tvalid := newTestInterviewValidate(t)
		got, err := tvalid.interviewValidate.ValidateGetInterviewAppointment(ctx)
		expected := id
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})
	t.Run("validate get interview appointment error when id is missing", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		tvalid := newTestInterviewValidate(t)
		got, err := tvalid.interviewValidate.ValidateGetInterviewAppointment(ctx)
		expected := helpers.NewCustomError(http.StatusBadRequest, "id: Missing required field")
		assert.Equal(t, "", got)
		assert.Equal(t, expected, err)
	})
	t.Run("validate get interview appointment error when id is invalid format", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = []gin.Param{
			{Key: "id", Value: "xxxxxxx"},
		}
		tvalid := newTestInterviewValidate(t)
		got, err := tvalid.interviewValidate.ValidateGetInterviewAppointment(ctx)
		expected := helpers.NewCustomError(http.StatusBadRequest, "id in param must be of type bsonobjectid: \"xxxxxxx\"")
		assert.Equal(t, "", got)
		assert.Equal(t, expected, err)
	})
}

func TestValidateCreateInterviewAppointment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	govalidator.SetFieldsRequiredByDefault(true)
	type requestBody struct {
		Title       string
		Description string
	}
	t.Run("validate create interview appointment success", func(t *testing.T) {
		body := requestBody{
			Title:       "title",
			Description: "description",
		}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Set("userId", "6476f457e64589e868aac97b")
		ctx.Request, _ = http.NewRequest("POST", "http://example.com", &buf)

		tvalid := newTestInterviewValidate(t)
		got, err := tvalid.interviewValidate.ValidateCreateInterviewAppointment(ctx)
		expected := &dto.CreateInterviewAppointmentRequest{
			Title:       "title",
			Description: "description",
			CreatedBy:   "6476f457e64589e868aac97b",
		}
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})
	t.Run("validate create interview appointments error when title is missing", func(t *testing.T) {
		body := requestBody{
			Description: "description",
		}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Set("userId", "6476f457e64589e868aac97b")
		ctx.Request, _ = http.NewRequest("POST", "http://example.com", &buf)

		tvalid := newTestInterviewValidate(t)
		got, err := tvalid.interviewValidate.ValidateCreateInterviewAppointment(ctx)
		expected := helpers.NewCustomError(http.StatusBadRequest, "title: Missing required field")
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
	t.Run("validate create interview appointments error when description is missing", func(t *testing.T) {
		body := requestBody{
			Title: "title",
		}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Set("userId", "6476f457e64589e868aac97b")
		ctx.Request, _ = http.NewRequest("POST", "http://example.com", &buf)

		tvalid := newTestInterviewValidate(t)
		got, err := tvalid.interviewValidate.ValidateCreateInterviewAppointment(ctx)
		expected := helpers.NewCustomError(http.StatusBadRequest, "description: Missing required field")
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
}

func TestValidateUpdateInterviewAppointment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	govalidator.SetFieldsRequiredByDefault(true)
	type requestBody struct {
		ID          string
		Title       string
		Description string
		Status      string
	}
	t.Run("validate update interview appointment success", func(t *testing.T) {
		id := "6476f457e64589e868aac97b"
		body := requestBody{Status: "DONE"}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = []gin.Param{
			{Key: "id", Value: id},
		}
		ctx.Request, _ = http.NewRequest("PATCH", "http://example.com", &buf)

		tvalid := newTestInterviewValidate(t)
		got, err := tvalid.interviewValidate.ValidateUpdateInterviewAppointment(ctx)
		expected := &dto.UpdateInterviewAppointmentRequest{
			ID:          id,
			Title:       "",
			Description: "",
			Status:      "DONE",
		}
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})
	t.Run("validate update interview appointment error when not input all field", func(t *testing.T) {
		id := "6476f457e64589e868aac97b"
		body := requestBody{}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = []gin.Param{
			{Key: "id", Value: id},
		}
		ctx.Request, _ = http.NewRequest("PATCH", "http://example.com", &buf)
		tvalid := newTestInterviewValidate(t)
		got, err := tvalid.interviewValidate.ValidateUpdateInterviewAppointment(ctx)
		expected := helpers.NewCustomError(http.StatusBadRequest, "at least one field required")
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
	t.Run("validate update interview appointment error when id is invalid format", func(t *testing.T) {
		id := "xxxxx"
		body := requestBody{Status: "DONE"}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = []gin.Param{
			{Key: "id", Value: id},
		}
		ctx.Request, _ = http.NewRequest("PATCH", "http://example.com", &buf)
		tvalid := newTestInterviewValidate(t)
		got, err := tvalid.interviewValidate.ValidateUpdateInterviewAppointment(ctx)
		expected := helpers.NewCustomError(http.StatusBadRequest, "id in body must be of type bsonobjectid: \"xxxxx\"")
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
}

func TestValidateArchiveInterviewAppointment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	govalidator.SetFieldsRequiredByDefault(true)
	t.Run("validate archive interview appointment success", func(t *testing.T) {
		id := "6476f457e64589e868aac97b"
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = []gin.Param{
			{Key: "id", Value: id},
		}

		tvalid := newTestInterviewValidate(t)
		got, err := tvalid.interviewValidate.ValidateArchiveInterviewAppointment(ctx)
		expected := id
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})
	t.Run("validate archive interview appointment error when id is missing", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		tvalid := newTestInterviewValidate(t)
		got, err := tvalid.interviewValidate.ValidateArchiveInterviewAppointment(ctx)
		expected := helpers.NewCustomError(http.StatusBadRequest, "id: Missing required field")
		assert.Equal(t, "", got)
		assert.Equal(t, expected, err)
	})
	t.Run("validate archive interview appointment error when id is invalid format", func(t *testing.T) {
		id := "xxxxxxx"
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = []gin.Param{
			{Key: "id", Value: id},
		}

		tvalid := newTestInterviewValidate(t)
		got, err := tvalid.interviewValidate.ValidateArchiveInterviewAppointment(ctx)
		expected := helpers.NewCustomError(http.StatusBadRequest, "id in param must be of type bsonobjectid: \"xxxxxxx\"")
		assert.Equal(t, "", got)
		assert.Equal(t, expected, err)
	})
}

func TestValidateAddInterviewComment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	govalidator.SetFieldsRequiredByDefault(true)
	type requestBody struct {
		Comment string
	}
	t.Run("validate add interview comment success", func(t *testing.T) {
		id := "6476f457e64589e868aac97b"
		userId := "64ac6cb9b0a3e8792efc438e"
		body := requestBody{Comment: "comment"}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = []gin.Param{
			{Key: "id", Value: id},
		}
		ctx.Request, _ = http.NewRequest("POST", "http://example.com", &buf)
		ctx.Set("userId", userId)
		tvalid := newTestInterviewValidate(t)
		got, err := tvalid.interviewValidate.ValidateAddInterviewComment(ctx)
		expected := &dto.AddInterviewCommentRequest{
			ID:      id,
			Comment: body.Comment,
			UserID:  userId,
		}
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})
	t.Run("validate add interview comment error when id is missing", func(t *testing.T) {
		userId := "64ac6cb9b0a3e8792efc438e"
		body := requestBody{Comment: "comment"}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request, _ = http.NewRequest("POST", "http://example.com", &buf)
		ctx.Set("userId", userId)
		tvalid := newTestInterviewValidate(t)
		got, err := tvalid.interviewValidate.ValidateAddInterviewComment(ctx)
		expected := helpers.NewCustomError(http.StatusBadRequest, "id: Missing required field")
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
	t.Run("validate add interview comment error when comment is missing", func(t *testing.T) {
		id := "6476f457e64589e868aac97b"
		userId := "64ac6cb9b0a3e8792efc438e"
		body := requestBody{}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = []gin.Param{
			{Key: "id", Value: id},
		}
		ctx.Request, _ = http.NewRequest("POST", "http://example.com", &buf)
		ctx.Set("userId", userId)
		tvalid := newTestInterviewValidate(t)
		got, err := tvalid.interviewValidate.ValidateAddInterviewComment(ctx)
		expected := helpers.NewCustomError(http.StatusBadRequest, "comment: Missing required field")
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
}

func TestValidateUpdateInterviewComment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	govalidator.SetFieldsRequiredByDefault(true)
	type requestBody struct {
		Comment string
	}
	t.Run("validate update interview comment success", func(t *testing.T) {
		id := "6476f457e64589e868aac97b"
		commentId := "64ace6bd981e163c387f494e"
		userId := "64ac6cb9b0a3e8792efc438e"
		body := requestBody{Comment: "update comment"}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = []gin.Param{
			{Key: "id", Value: id},
			{Key: "commentId", Value: commentId},
		}
		ctx.Request, _ = http.NewRequest("POST", "http://example.com", &buf)
		ctx.Set("userId", userId)
		tvalid := newTestInterviewValidate(t)
		got, err := tvalid.interviewValidate.ValidateUpdateInterviewComment(ctx)
		expected := &dto.UpdateInterviewCommentRequest{
			ID:        id,
			CommentID: commentId,
			Comment:   body.Comment,
			UserID:    userId,
		}
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})
	t.Run("validate update interview comment error when id is missing", func(t *testing.T) {
		commentId := "64ace6bd981e163c387f494e"
		userId := "64ac6cb9b0a3e8792efc438e"
		body := requestBody{Comment: "update comment"}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = []gin.Param{
			{Key: "commentId", Value: commentId},
		}
		ctx.Request, _ = http.NewRequest("POST", "http://example.com", &buf)
		ctx.Set("userId", userId)
		tvalid := newTestInterviewValidate(t)
		got, err := tvalid.interviewValidate.ValidateUpdateInterviewComment(ctx)
		expected := helpers.NewCustomError(http.StatusBadRequest, "id: Missing required field")
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
	t.Run("validate update interview comment error when comment id is missing", func(t *testing.T) {
		id := "6476f457e64589e868aac97b"
		userId := "64ac6cb9b0a3e8792efc438e"
		body := requestBody{Comment: "update comment"}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = []gin.Param{
			{Key: "id", Value: id},
		}
		ctx.Request, _ = http.NewRequest("POST", "http://example.com", &buf)
		ctx.Set("userId", userId)
		tvalid := newTestInterviewValidate(t)
		got, err := tvalid.interviewValidate.ValidateUpdateInterviewComment(ctx)
		expected := helpers.NewCustomError(http.StatusBadRequest, "commentId: Missing required field")
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
	t.Run("validate update interview comment error when comment is missing", func(t *testing.T) {
		id := "6476f457e64589e868aac97b"
		commentId := "64ace6bd981e163c387f494e"
		userId := "64ac6cb9b0a3e8792efc438e"
		body := requestBody{}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = []gin.Param{
			{Key: "id", Value: id},
			{Key: "commentId", Value: commentId},
		}
		ctx.Request, _ = http.NewRequest("POST", "http://example.com", &buf)
		ctx.Set("userId", userId)
		tvalid := newTestInterviewValidate(t)
		got, err := tvalid.interviewValidate.ValidateUpdateInterviewComment(ctx)
		expected := helpers.NewCustomError(http.StatusBadRequest, "comment: Missing required field")
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
	t.Run("validate update interview comment error when id is invalid format", func(t *testing.T) {
		id := "xxxxxxx"
		commentId := "64ace6bd981e163c387f494e"
		userId := "64ac6cb9b0a3e8792efc438e"
		body := requestBody{Comment: "update comment"}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = []gin.Param{
			{Key: "id", Value: id},
			{Key: "commentId", Value: commentId},
		}
		ctx.Request, _ = http.NewRequest("POST", "http://example.com", &buf)
		ctx.Set("userId", userId)
		tvalid := newTestInterviewValidate(t)
		got, err := tvalid.interviewValidate.ValidateUpdateInterviewComment(ctx)
		expected := helpers.NewCustomError(http.StatusBadRequest, "id in param must be of type bsonobjectid: \"xxxxxxx\"")
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
	t.Run("validate update interview comment error when comment id is invalid format", func(t *testing.T) {
		id := "6476f457e64589e868aac97b"
		commentId := "xxxxxxx"
		userId := "64ac6cb9b0a3e8792efc438e"
		body := requestBody{Comment: "update comment"}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Params = []gin.Param{
			{Key: "id", Value: id},
			{Key: "commentId", Value: commentId},
		}
		ctx.Request, _ = http.NewRequest("POST", "http://example.com", &buf)
		ctx.Set("userId", userId)
		tvalid := newTestInterviewValidate(t)
		got, err := tvalid.interviewValidate.ValidateUpdateInterviewComment(ctx)
		expected := helpers.NewCustomError(http.StatusBadRequest, "commentId in param must be of type bsonobjectid: \"xxxxxxx\"")
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
}
