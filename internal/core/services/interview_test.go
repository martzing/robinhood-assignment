package services_test

import (
	"context"
	"errors"
	"net/http"
	"robinhood-assignment/helpers"
	"robinhood-assignment/internal/core/domains"
	"robinhood-assignment/internal/core/ports"
	"robinhood-assignment/internal/core/ports/mocks"
	"robinhood-assignment/internal/core/services"
	"robinhood-assignment/internal/dto"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type testInterviewService struct {
	interviewAppointmentRepo *mocks.InterviewAppointmentRepository
	userRepo                 *mocks.UserRepository
	service                  ports.InterviewService
}

func newTestInterviewService(t *testing.T) testInterviewService {
	interviewAppointmentRepo := mocks.NewInterviewAppointmentRepository(t)
	userRepo := mocks.NewUserRepository(t)

	service := services.NewInterviewService(interviewAppointmentRepo, userRepo)
	return testInterviewService{interviewAppointmentRepo, userRepo, service}
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
	t.Run("Get interview appointments success", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		offset := uint32(0)
		limit := uint32(3)
		expected := []domains.InterviewAppointment{mockInterviewAppointment1, mockInterviewAppointment2}
		tsvc.interviewAppointmentRepo.On("GetAll", ctx, offset, limit).Return(expected, nil)
		got, err := tsvc.service.GetInterviewAppointments(ctx, offset, limit)
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})
	t.Run("Get interview appointments error", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		offset := uint32(0)
		limit := uint32(3)
		expected := helpers.NewCustomError(http.StatusInternalServerError, "Cannot get interview appointment.")
		tsvc.interviewAppointmentRepo.On("GetAll", ctx, offset, limit).Return(nil, expected)
		got, err := tsvc.service.GetInterviewAppointments(ctx, offset, limit)
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
}

func TestGetInterviewAppointment(t *testing.T) {
	ctx := context.Background()
	t.Run("Get interview appointment success", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		id := "64aaf0156999249a602ff55f"
		objId, _ := primitive.ObjectIDFromHex(id)
		expected := &mockInterviewAppointment1
		tsvc.interviewAppointmentRepo.On("Get", ctx, objId).Return(expected, nil)
		got, err := tsvc.service.GetInterviewAppointment(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})
	t.Run("Get interview appointment error when invalid id format", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		id := "xxxxx"
		expected := helpers.InternalError
		got, err := tsvc.service.GetInterviewAppointment(ctx, id)
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
	t.Run("Get interview appointment error when query error", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		id := "64aaf0156999249a602ff55f"
		objId, _ := primitive.ObjectIDFromHex(id)
		expected := helpers.InternalError
		tsvc.interviewAppointmentRepo.On("Get", ctx, objId).Return(nil, errors.New("some error"))
		got, err := tsvc.service.GetInterviewAppointment(ctx, id)
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
	t.Run("Get interview appointment error when data not found", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		id := "64aaf0156999249a602ff55f"
		objId, _ := primitive.ObjectIDFromHex(id)
		expected := helpers.NewCustomError(http.StatusNotFound, "Interview appointment not found.")
		tsvc.interviewAppointmentRepo.On("Get", ctx, objId).Return(nil, nil)
		got, err := tsvc.service.GetInterviewAppointment(ctx, id)
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
}

func TestCreateInterviewAppointment(t *testing.T) {
	t.Run("Create interview appointment success", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		userId := "6476f457e64589e868aac977"
		req := &dto.CreateInterviewAppointmentRequest{
			Title:       "Title",
			Description: "Description",
			CreatedBy:   userId,
		}
		userObjId, _ := primitive.ObjectIDFromHex(userId)
		params := &domains.CreateInterviewAppointmentParams{
			Title:       req.Title,
			Description: req.Description,
			UserID:      userObjId,
		}
		user := &domains.User{
			ID:       userObjId,
			Name:     "User name 1",
			Email:    "User email 1",
			Username: "Username 1",
			Password: "Password 1",
			ImageUrl: "https://image-url.com",
			Role:     "ADMIN",
		}
		created := &domains.CreateInterviewAppointment{
			ID:           primitive.NewObjectID(),
			Title:        params.Title,
			Description:  params.Description,
			Comments:     []domains.InterviewComment{},
			Status:       "TODO",
			IsArchived:   false,
			CreateUserId: userObjId,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		expected := &domains.InterviewAppointment{
			ID:          created.ID,
			Title:       created.Title,
			Description: created.Description,
			Comments:    created.Comments,
			Status:      created.Status,
			IsArchived:  created.IsArchived,
			CreateUser: domains.User{
				ID:       user.ID,
				Name:     user.Name,
				Email:    user.Email,
				ImageUrl: user.ImageUrl,
			},
			CreatedAt: created.CreatedAt,
			UpdatedAt: created.UpdatedAt,
		}
		tsvc.userRepo.On("Get", ctx, userObjId).Return(user, nil)
		tsvc.interviewAppointmentRepo.On("Create", ctx, params).Return(created, nil)
		got, err := tsvc.service.CreateInterviewAppointment(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})
	t.Run("Create interview appointment error when invalid user id", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		userId := "xxxxx"
		req := &dto.CreateInterviewAppointmentRequest{
			Title:       "Title",
			Description: "Description",
			CreatedBy:   userId,
		}
		expected := helpers.InternalError
		got, err := tsvc.service.CreateInterviewAppointment(ctx, req)
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
	t.Run("Create interview appointment error when query user fail", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		userId := "6476f457e64589e868aac977"
		req := &dto.CreateInterviewAppointmentRequest{
			Title:       "Title",
			Description: "Description",
			CreatedBy:   userId,
		}
		userObjId, _ := primitive.ObjectIDFromHex(userId)
		expected := helpers.InternalError
		tsvc.userRepo.On("Get", ctx, userObjId).Return(nil, errors.New("some error"))
		got, err := tsvc.service.CreateInterviewAppointment(ctx, req)
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
	t.Run("Create interview appointment error when not exist user", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		userId := "6476f457e64589e868aac977"
		req := &dto.CreateInterviewAppointmentRequest{
			Title:       "Title",
			Description: "Description",
			CreatedBy:   userId,
		}
		userObjId, _ := primitive.ObjectIDFromHex(userId)
		expected := helpers.NewCustomError(http.StatusUnauthorized, "Invalid user token")
		tsvc.userRepo.On("Get", ctx, userObjId).Return(nil, nil)
		got, err := tsvc.service.CreateInterviewAppointment(ctx, req)
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
	t.Run("Create interview appointment error when query create fail", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		userId := "6476f457e64589e868aac977"
		req := &dto.CreateInterviewAppointmentRequest{
			Title:       "Title",
			Description: "Description",
			CreatedBy:   userId,
		}
		userObjId, _ := primitive.ObjectIDFromHex(userId)
		params := &domains.CreateInterviewAppointmentParams{
			Title:       req.Title,
			Description: req.Description,
			UserID:      userObjId,
		}
		user := &domains.User{
			ID:       userObjId,
			Name:     "User name 1",
			Email:    "User email 1",
			Username: "Username 1",
			Password: "Password 1",
			ImageUrl: "https://image-url.com",
			Role:     "ADMIN",
		}
		expected := helpers.InternalError
		tsvc.userRepo.On("Get", ctx, userObjId).Return(user, nil)
		tsvc.interviewAppointmentRepo.On("Create", ctx, params).Return(nil, errors.New("some error"))
		got, err := tsvc.service.CreateInterviewAppointment(ctx, req)
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
}
