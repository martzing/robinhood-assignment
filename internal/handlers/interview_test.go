package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"robinhood-assignment/helpers"
	"robinhood-assignment/internal/core/domains"
	"robinhood-assignment/internal/core/ports"
	"robinhood-assignment/internal/core/ports/mocks"
	"robinhood-assignment/internal/dto"
	"robinhood-assignment/internal/handlers"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type testInterviewHandler struct {
	interviewService  *mocks.InterviewService
	interviewValidate *mocks.InterviewValidate
	handler           ports.InterviewHandler
}

func newTestInterviewHandler(t *testing.T) testInterviewHandler {
	interviewService := mocks.NewInterviewService(t)
	interviewValidate := mocks.NewInterviewValidate(t)
	handler := handlers.NewInterviewHandler(interviewService, interviewValidate)
	return testInterviewHandler{interviewService, interviewValidate, handler}
}

var (
	ctx                       = context.Background()
	now                       = time.Now()
	mockInterviewAppointment1 = domains.InterviewAppointment{
		ID:          primitive.NewObjectID(),
		Title:       "Title 1",
		Description: "Description 1",
		Comments:    []domains.InterviewComment{},
		Status:      "TODO",
		IsArchived:  false,
		CreateUser: domains.User{
			ID:       primitive.NewObjectID(),
			Name:     "User name 1",
			Email:    "User email 1",
			Username: "Username 1",
			Password: "Password 1",
			ImageUrl: "https://image-url.com",
			Role:     "ADMIN",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	mockInterviewAppointment2 = domains.InterviewAppointment{
		ID:          primitive.NewObjectID(),
		Title:       "Title 2",
		Description: "Description 2",
		Comments:    []domains.InterviewComment{},
		Status:      "TODO",
		IsArchived:  false,
		CreateUser: domains.User{
			ID:       primitive.NewObjectID(),
			Name:     "User name 2",
			Email:    "User email 2",
			Username: "Username 2",
			Password: "Password 2",
			ImageUrl: "https://image-url.com",
			Role:     "ADMIN",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
)

func TestGetInterviewAppointments(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("get interview appointments success", func(t *testing.T) {
		req := dto.GetInterviewAppointmentsRequest{
			Page:  1,
			Limit: 20,
		}
		if req.Page < 1 {
			req.Page = 1
		}
		if req.Limit < 1 {
			req.Limit = 20
		}
		offset := (req.Page - 1) * req.Limit
		limit := req.Limit + 1
		data := []domains.InterviewAppointment{mockInterviewAppointment1, mockInterviewAppointment2}
		interviews := make([]dto.InterviewAppointment, len(data))
		for i := 0; i < len(data); i++ {
			interviews[i] = dto.InterviewAppointment{
				ID:          data[i].ID.Hex(),
				Title:       data[i].Title,
				Description: data[i].Description,
				Status:      data[i].Status,
				CreateUser: dto.User{
					Name:     data[i].CreateUser.Name,
					Email:    data[i].CreateUser.Email,
					ImageUrl: data[i].CreateUser.ImageUrl,
				},
				CreatedAt: data[i].CreatedAt,
			}
		}
		size, hasNext := helpers.Paginate(&interviews, int64(req.Limit))
		res := dto.GetInterviewAppointmentsResponse{
			StatusCode: http.StatusOK,
			Data:       interviews,
			Pagination: dto.Pagination{
				Page:    uint32(req.Page),
				Size:    uint32(size),
				HasNext: hasNext,
			},
		}

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		thld := newTestInterviewHandler(t)
		thld.interviewValidate.On("ValidateGetInterviewAppointments", ctx).Return(&req, nil)
		thld.interviewService.On("GetInterviewAppointments", ctx, offset, limit).Return(data, nil)
		thld.handler.GetInterviewAppointments(ctx)
		expected, _ := json.Marshal(res)
		got := w.Body.Bytes()
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expected, got)
	})
	t.Run("get interview appointments error when validate fail", func(t *testing.T) {
		req := dto.GetInterviewAppointmentsRequest{
			Page:  1,
			Limit: 20,
		}
		errMsg := "Invalid page query parameter"
		res := &dto.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      errMsg,
		}
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		thld := newTestInterviewHandler(t)
		thld.interviewValidate.On("ValidateGetInterviewAppointments", ctx).Return(&req, helpers.NewCustomError(http.StatusBadRequest, errMsg))
		thld.handler.GetInterviewAppointments(ctx)
		expected, _ := json.Marshal(res)
		got := w.Body.Bytes()
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expected, got)
	})
	t.Run("get interview appointments error when call service fail", func(t *testing.T) {
		req := dto.GetInterviewAppointmentsRequest{
			Page:  1,
			Limit: 20,
		}
		if req.Page < 1 {
			req.Page = 1
		}
		if req.Limit < 1 {
			req.Limit = 20
		}
		offset := (req.Page - 1) * req.Limit
		limit := req.Limit + 1

		errMsg := "Cannot get interview appointment"
		res := &dto.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      errMsg,
		}
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		thld := newTestInterviewHandler(t)
		thld.interviewValidate.On("ValidateGetInterviewAppointments", ctx).Return(&req, nil)
		thld.interviewService.On("GetInterviewAppointments", ctx, offset, limit).Return(nil, helpers.NewCustomError(http.StatusInternalServerError, errMsg))
		thld.handler.GetInterviewAppointments(ctx)
		expected, _ := json.Marshal(res)
		got := w.Body.Bytes()
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, expected, got)
	})
}

func TestGetInterviewAppointment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("get interview appointment success", func(t *testing.T) {
		data := mockInterviewAppointment1
		data.Comments = append(data.Comments, []domains.InterviewComment{
			{
				ID:      primitive.NewObjectID(),
				Comment: "comment",
				User: domains.User{
					ID:       primitive.NewObjectID(),
					Name:     "User name 1",
					Email:    "User email 1",
					Username: "Username 1",
					Password: "Password 1",
					ImageUrl: "https://image-url.com",
					Role:     "ADMIN",
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:      primitive.NewObjectID(),
				Comment: "comment 2",
				User: domains.User{
					ID:       primitive.NewObjectID(),
					Name:     "User name 2",
					Email:    "User email 2",
					Username: "Username 2",
					Password: "Password 2",
					ImageUrl: "https://image-url.com",
					Role:     "STAFF",
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
		}...)

		comments := []dto.InterviewComment{}
		for i := 0; i < len(data.Comments); i++ {
			if !data.Comments[i].ID.IsZero() {
				comments = append(comments, dto.InterviewComment{
					ID:      data.Comments[i].ID.Hex(),
					Comment: data.Comments[i].Comment,
					User: dto.User{
						Name:     data.Comments[i].User.Name,
						Email:    data.Comments[i].User.Email,
						ImageUrl: data.Comments[i].User.ImageUrl,
					},
					CreatedAt: data.CreatedAt,
				})
			}
		}
		res := dto.GetInterviewAppointmentResponse{
			StatusCode: http.StatusOK,
			Data: dto.InterviewAppointmentDetail{
				ID:          data.ID.Hex(),
				Title:       data.Title,
				Description: data.Description,
				Status:      data.Status,
				CreateUser: dto.User{
					Name:     data.CreateUser.Name,
					Email:    data.CreateUser.Email,
					ImageUrl: data.CreateUser.ImageUrl,
				},
				CreatedAt: data.CreatedAt,
				Comments:  comments,
			},
		}

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		thld := newTestInterviewHandler(t)
		thld.interviewValidate.On("ValidateGetInterviewAppointment", ctx).Return(data.ID.Hex(), nil)
		thld.interviewService.On("GetInterviewAppointment", ctx, data.ID.Hex()).Return(&data, nil)
		thld.handler.GetInterviewAppointment(ctx)
		expected, _ := json.Marshal(res)
		got := w.Body.Bytes()
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expected, got)
	})
	t.Run("get interview appointment error when validate fail", func(t *testing.T) {
		errMsg := "id: Missing required field"
		res := &dto.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      errMsg,
		}

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		thld := newTestInterviewHandler(t)
		thld.interviewValidate.On("ValidateGetInterviewAppointment", ctx).Return("", helpers.NewCustomError(http.StatusBadRequest, errMsg))

		thld.handler.GetInterviewAppointment(ctx)
		expected, _ := json.Marshal(res)
		got := w.Body.Bytes()
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expected, got)
	})
	t.Run("get interview appointment error when call service fail", func(t *testing.T) {
		id := "6476f457e64589e868aac97e"
		errMsg := "Interview appointment not found."
		res := &dto.ErrorResponse{
			StatusCode: http.StatusNotFound,
			Error:      errMsg,
		}

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		thld := newTestInterviewHandler(t)
		thld.interviewValidate.On("ValidateGetInterviewAppointment", ctx).Return(id, nil)
		thld.interviewService.On("GetInterviewAppointment", ctx, id).Return(nil, helpers.NewCustomError(http.StatusNotFound, errMsg))
		thld.handler.GetInterviewAppointment(ctx)
		expected, _ := json.Marshal(res)
		got := w.Body.Bytes()
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, expected, got)
	})
}

func TestCreateInterviewAppointment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("create interview appointment success", func(t *testing.T) {
		req := dto.CreateInterviewAppointmentRequest{
			Title:       mockInterviewAppointment1.Title,
			Description: mockInterviewAppointment1.Description,
			CreatedBy:   mockInterviewAppointment1.CreateUser.ID.Hex(),
		}
		data := mockInterviewAppointment1
		res := dto.CreateInterviewAppointmentResponse{
			StatusCode: http.StatusCreated,
			Data: dto.InterviewAppointmentDetail{
				ID:          data.ID.Hex(),
				Title:       data.Title,
				Description: data.Description,
				Status:      data.Status,
				CreateUser: dto.User{
					Name:     data.CreateUser.Name,
					Email:    data.CreateUser.Email,
					ImageUrl: data.CreateUser.ImageUrl,
				},
				CreatedAt: data.CreatedAt,
				Comments:  []dto.InterviewComment{},
			},
		}

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		thld := newTestInterviewHandler(t)
		thld.interviewValidate.On("ValidateCreateInterviewAppointment", ctx).Return(&req, nil)
		thld.interviewService.On("CreateInterviewAppointment", ctx, &req).Return(&data, nil)
		thld.handler.CreateInterviewAppointment(ctx)
		expected, _ := json.Marshal(res)
		got := w.Body.Bytes()
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, expected, got)
	})
	t.Run("create interview appointment error when validate fail", func(t *testing.T) {
		errMsg := "Invalid input parameter"
		res := &dto.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      errMsg,
		}
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		thld := newTestInterviewHandler(t)
		thld.interviewValidate.On("ValidateCreateInterviewAppointment", ctx).Return(nil, helpers.NewCustomError(http.StatusBadRequest, errMsg))
		thld.handler.CreateInterviewAppointment(ctx)
		expected, _ := json.Marshal(res)
		got := w.Body.Bytes()
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expected, got)
	})
	t.Run("create interview appointment error when call service fail", func(t *testing.T) {
		req := dto.CreateInterviewAppointmentRequest{
			Title:       mockInterviewAppointment1.Title,
			Description: mockInterviewAppointment1.Description,
			CreatedBy:   mockInterviewAppointment1.CreateUser.ID.Hex(),
		}
		errMsg := "Something went wrong please contact developer."
		res := &dto.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Error:      errMsg,
		}

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		thld := newTestInterviewHandler(t)
		thld.interviewValidate.On("ValidateCreateInterviewAppointment", ctx).Return(&req, nil)
		thld.interviewService.On("CreateInterviewAppointment", ctx, &req).Return(nil, helpers.InternalError)
		thld.handler.CreateInterviewAppointment(ctx)
		expected, _ := json.Marshal(res)
		got := w.Body.Bytes()
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, expected, got)
	})
}

func TestUpdateInterviewAppointment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("update interview appointment success", func(t *testing.T) {
		req := dto.UpdateInterviewAppointmentRequest{
			ID:          mockInterviewAppointment1.ID.Hex(),
			Title:       "Update title",
			Description: "Update description",
			Status:      "IN_PROGRESS",
		}
		res := dto.BaseResponse{
			StatusCode: http.StatusOK,
			Message:    "success",
		}
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		thld := newTestInterviewHandler(t)
		thld.interviewValidate.On("ValidateUpdateInterviewAppointment", ctx).Return(&req, nil)
		thld.interviewService.On("UpdateInterviewAppointment", ctx, &req).Return(nil)
		thld.handler.UpdateInterviewAppointment(ctx)
		expected, _ := json.Marshal(res)
		got := w.Body.Bytes()
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expected, got)
	})
	t.Run("update interview appointment error when validate fail", func(t *testing.T) {
		errMsg := "at least one field required"
		res := &dto.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      errMsg,
		}
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		thld := newTestInterviewHandler(t)
		thld.interviewValidate.On("ValidateUpdateInterviewAppointment", ctx).Return(nil, helpers.NewCustomError(http.StatusBadRequest, errMsg))
		thld.handler.UpdateInterviewAppointment(ctx)
		expected, _ := json.Marshal(res)
		got := w.Body.Bytes()
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expected, got)
	})
	t.Run("update interview appointment error when call service fail", func(t *testing.T) {
		req := dto.UpdateInterviewAppointmentRequest{
			ID:          mockInterviewAppointment1.ID.Hex(),
			Title:       "Update title",
			Description: "Update description",
			Status:      "IN_PROGRESS",
		}
		errMsg := "Interview appointment not found."
		res := &dto.ErrorResponse{
			StatusCode: http.StatusNotFound,
			Error:      errMsg,
		}
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		thld := newTestInterviewHandler(t)
		thld.interviewValidate.On("ValidateUpdateInterviewAppointment", ctx).Return(&req, nil)
		thld.interviewService.On("UpdateInterviewAppointment", ctx, &req).Return(helpers.NewCustomError(http.StatusNotFound, errMsg))
		thld.handler.UpdateInterviewAppointment(ctx)
		expected, _ := json.Marshal(res)
		got := w.Body.Bytes()
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, expected, got)
	})
}

func TestArchiveInterviewAppointment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("archive interview appointment success", func(t *testing.T) {
		id := "6476f457e64589e868aac981"
		res := dto.BaseResponse{
			StatusCode: http.StatusOK,
			Message:    "success",
		}
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		thld := newTestInterviewHandler(t)
		thld.interviewValidate.On("ValidateGetInterviewAppointment", ctx).Return(id, nil)
		thld.interviewService.On("ArchiveInterviewAppointment", ctx, id).Return(nil)
		thld.handler.ArchiveInterviewAppointment(ctx)
		expected, _ := json.Marshal(res)
		got := w.Body.Bytes()
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expected, got)
	})
	t.Run("archive interview appointment error when validate fail", func(t *testing.T) {
		errMsg := "id: Missing required field"
		res := &dto.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      errMsg,
		}
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		thld := newTestInterviewHandler(t)
		thld.interviewValidate.On("ValidateGetInterviewAppointment", ctx).Return("", helpers.NewCustomError(http.StatusBadRequest, errMsg))
		thld.handler.ArchiveInterviewAppointment(ctx)
		expected, _ := json.Marshal(res)
		got := w.Body.Bytes()
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expected, got)
	})
	t.Run("archive interview appointment error when call service fail", func(t *testing.T) {
		id := "6476f457e64589e868aac981"
		errMsg := "Interview appointment not found or archived"
		res := &dto.ErrorResponse{
			StatusCode: http.StatusNotFound,
			Error:      errMsg,
		}
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		thld := newTestInterviewHandler(t)
		thld.interviewValidate.On("ValidateGetInterviewAppointment", ctx).Return(id, nil)
		thld.interviewService.On("ArchiveInterviewAppointment", ctx, id).Return(helpers.NewCustomError(http.StatusNotFound, errMsg))
		thld.handler.ArchiveInterviewAppointment(ctx)
		expected, _ := json.Marshal(res)
		got := w.Body.Bytes()
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, expected, got)
	})
}

func TestAddInterviewComment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("add interview comment success", func(t *testing.T) {
		req := dto.AddInterviewCommentRequest{
			ID:      "6476f457e64589e868aac981",
			Comment: "comment",
			UserID:  "6476f457e64589e868aac984",
		}
		res := dto.BaseResponse{
			StatusCode: http.StatusOK,
			Message:    "success",
		}
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		thld := newTestInterviewHandler(t)
		thld.interviewValidate.On("ValidateAddInterviewComment", ctx).Return(&req, nil)
		thld.interviewService.On("AddInterviewComment", ctx, &req).Return(nil)
		thld.handler.AddInterviewComment(ctx)
		expected, _ := json.Marshal(res)
		got := w.Body.Bytes()
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expected, got)
	})
	t.Run("add interview comment error when validate fail", func(t *testing.T) {
		errMsg := "id: Missing required field"
		res := &dto.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      errMsg,
		}
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		thld := newTestInterviewHandler(t)
		thld.interviewValidate.On("ValidateAddInterviewComment", ctx).Return(nil, helpers.NewCustomError(http.StatusBadRequest, errMsg))
		thld.handler.AddInterviewComment(ctx)
		expected, _ := json.Marshal(res)
		got := w.Body.Bytes()
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expected, got)
	})
	t.Run("add interview comment error when call service fail", func(t *testing.T) {
		req := dto.AddInterviewCommentRequest{
			ID:      "6476f457e64589e868aac981",
			Comment: "comment",
			UserID:  "6476f457e64589e868aac984",
		}
		errMsg := "Interview appointment not found."
		res := &dto.ErrorResponse{
			StatusCode: http.StatusNotFound,
			Error:      errMsg,
		}
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		thld := newTestInterviewHandler(t)
		thld.interviewValidate.On("ValidateAddInterviewComment", ctx).Return(&req, nil)
		thld.interviewService.On("AddInterviewComment", ctx, &req).Return(helpers.NewCustomError(http.StatusNotFound, errMsg))
		thld.handler.AddInterviewComment(ctx)
		expected, _ := json.Marshal(res)
		got := w.Body.Bytes()
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, expected, got)
	})
}

func TestUpdateInterviewComment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("update interview comment success", func(t *testing.T) {
		req := dto.UpdateInterviewCommentRequest{
			ID:        "6476f457e64589e868aac981",
			CommentID: "6476f457e64589e868aac987",
			Comment:   "comment",
			UserID:    "6476f457e64589e868aac984",
		}
		res := dto.BaseResponse{
			StatusCode: http.StatusOK,
			Message:    "success",
		}
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		thld := newTestInterviewHandler(t)
		thld.interviewValidate.On("ValidateUpdateInterviewComment", ctx).Return(&req, nil)
		thld.interviewService.On("UpdateInterviewComment", ctx, &req).Return(nil)
		thld.handler.UpdateInterviewComment(ctx)
		expected, _ := json.Marshal(res)
		got := w.Body.Bytes()
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expected, got)
	})
	t.Run("update interview comment error when validate fail", func(t *testing.T) {
		errMsg := "id: Missing required field"
		res := &dto.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Error:      errMsg,
		}
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		thld := newTestInterviewHandler(t)
		thld.interviewValidate.On("ValidateUpdateInterviewComment", ctx).Return(nil, helpers.NewCustomError(http.StatusBadRequest, errMsg))
		thld.handler.UpdateInterviewComment(ctx)
		expected, _ := json.Marshal(res)
		got := w.Body.Bytes()
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expected, got)
	})
	t.Run("update interview comment error when call service fail", func(t *testing.T) {
		req := dto.UpdateInterviewCommentRequest{
			ID:        "6476f457e64589e868aac981",
			CommentID: "6476f457e64589e868aac987",
			Comment:   "comment",
			UserID:    "6476f457e64589e868aac984",
		}
		errMsg := "You don't have permission to update this comment"
		res := &dto.ErrorResponse{
			StatusCode: http.StatusForbidden,
			Error:      errMsg,
		}
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		thld := newTestInterviewHandler(t)
		thld.interviewValidate.On("ValidateUpdateInterviewComment", ctx).Return(&req, nil)
		thld.interviewService.On("UpdateInterviewComment", ctx, &req).Return(helpers.NewCustomError(http.StatusForbidden, errMsg))
		thld.handler.UpdateInterviewComment(ctx)
		expected, _ := json.Marshal(res)
		got := w.Body.Bytes()
		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Equal(t, expected, got)
	})
}
