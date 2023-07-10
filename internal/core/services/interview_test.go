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
	"go.mongodb.org/mongo-driver/mongo"
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
	t.Run("get interview appointments success", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		offset := uint32(0)
		limit := uint32(3)
		expected := []domains.InterviewAppointment{mockInterviewAppointment1, mockInterviewAppointment2}
		tsvc.interviewAppointmentRepo.On("GetAll", ctx, offset, limit).Return(expected, nil)
		got, err := tsvc.service.GetInterviewAppointments(ctx, offset, limit)
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})
	t.Run("get interview appointments error", func(t *testing.T) {
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
	t.Run("get interview appointment success", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		id := "64aaf0156999249a602ff55f"
		objId, _ := primitive.ObjectIDFromHex(id)
		expected := &mockInterviewAppointment1
		tsvc.interviewAppointmentRepo.On("Get", ctx, objId).Return(expected, nil)
		got, err := tsvc.service.GetInterviewAppointment(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})
	t.Run("get interview appointment error when invalid id format", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		id := "xxxxx"
		expected := helpers.InternalError
		got, err := tsvc.service.GetInterviewAppointment(ctx, id)
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
	t.Run("get interview appointment error when query error", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		id := "64aaf0156999249a602ff55f"
		objId, _ := primitive.ObjectIDFromHex(id)
		expected := helpers.InternalError
		tsvc.interviewAppointmentRepo.On("Get", ctx, objId).Return(nil, errors.New("some error"))
		got, err := tsvc.service.GetInterviewAppointment(ctx, id)
		assert.Nil(t, got)
		assert.Equal(t, expected, err)
	})
	t.Run("get interview appointment error when data not found", func(t *testing.T) {
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
	t.Run("create interview appointment success", func(t *testing.T) {
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
	t.Run("create interview appointment error when invalid user id", func(t *testing.T) {
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
	t.Run("create interview appointment error when query user fail", func(t *testing.T) {
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
	t.Run("create interview appointment error when not exist user", func(t *testing.T) {
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
	t.Run("create interview appointment error when query create fail", func(t *testing.T) {
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

func TestUpdateInterviewAppointment(t *testing.T) {
	t.Run("update interview appointment success", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		id := "64aaf0156999249a602ff55f"
		objId, _ := primitive.ObjectIDFromHex(id)
		req := &dto.UpdateInterviewAppointmentRequest{
			ID:          id,
			Title:       "Title",
			Description: "Description",
			Status:      "IN_PROGRESS",
		}
		params := &domains.UpdateInterviewAppointmentParams{
			ID:          objId,
			Title:       req.Title,
			Description: req.Description,
			Status:      req.Status,
		}
		updated := &domains.InterviewAppointment{
			ID:          params.ID,
			Title:       params.Title,
			Description: params.Description,
			Comments:    []domains.InterviewComment{},
			Status:      params.Status,
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
		tsvc.interviewAppointmentRepo.On("Update", ctx, params).Return(updated, nil)
		err := tsvc.service.UpdateInterviewAppointment(ctx, req)
		assert.NoError(t, err)
	})
	t.Run("update interview appointment error when invalid id format", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		id := "xxxxxx"
		req := &dto.UpdateInterviewAppointmentRequest{
			ID:          id,
			Title:       "Title",
			Description: "Description",
			Status:      "IN_PROGRESS",
		}
		expected := helpers.InternalError
		err := tsvc.service.UpdateInterviewAppointment(ctx, req)
		assert.Equal(t, expected, err)
	})
	t.Run("update interview appointment error when data not found", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		id := "64aaf0156999249a602ff55f"
		objId, _ := primitive.ObjectIDFromHex(id)
		req := &dto.UpdateInterviewAppointmentRequest{
			ID:          id,
			Title:       "Title",
			Description: "Description",
			Status:      "IN_PROGRESS",
		}
		params := &domains.UpdateInterviewAppointmentParams{
			ID:          objId,
			Title:       req.Title,
			Description: req.Description,
			Status:      req.Status,
		}
		expected := helpers.NewCustomError(http.StatusNotFound, "Interview appointment not found.")
		tsvc.interviewAppointmentRepo.On("Update", ctx, params).Return(nil, nil)
		err := tsvc.service.UpdateInterviewAppointment(ctx, req)
		assert.Equal(t, expected, err)
	})
	t.Run("update interview appointment error when update query fail", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		id := "64aaf0156999249a602ff55f"
		objId, _ := primitive.ObjectIDFromHex(id)
		req := &dto.UpdateInterviewAppointmentRequest{
			ID:          id,
			Title:       "Title",
			Description: "Description",
			Status:      "IN_PROGRESS",
		}
		params := &domains.UpdateInterviewAppointmentParams{
			ID:          objId,
			Title:       req.Title,
			Description: req.Description,
			Status:      req.Status,
		}
		expected := helpers.InternalError
		tsvc.interviewAppointmentRepo.On("Update", ctx, params).Return(nil, errors.New("some error"))
		err := tsvc.service.UpdateInterviewAppointment(ctx, req)
		assert.Equal(t, expected, err)
	})
}

func TestArchiveInterviewAppointment(t *testing.T) {
	t.Run("archive interview appointment success", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		id := "64aaf0156999249a602ff55f"
		objId, _ := primitive.ObjectIDFromHex(id)
		tsvc.interviewAppointmentRepo.On("ArchiveInterviewAppointment", ctx, objId).Return(nil)
		err := tsvc.service.ArchiveInterviewAppointment(ctx, id)
		assert.NoError(t, err)
	})
	t.Run("archive interview appointment error when invalid id format", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		id := "xxxxxxxxx"
		expected := helpers.InternalError
		err := tsvc.service.ArchiveInterviewAppointment(ctx, id)
		assert.Equal(t, expected, err)
	})
	t.Run("archive interview appointment error when data not found", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		id := "64aaf0156999249a602ff55f"
		objId, _ := primitive.ObjectIDFromHex(id)
		expected := helpers.NewCustomError(http.StatusNotFound, "Interview appointment not found or archived")
		tsvc.interviewAppointmentRepo.On("ArchiveInterviewAppointment", ctx, objId).Return(mongo.ErrNoDocuments)
		err := tsvc.service.ArchiveInterviewAppointment(ctx, id)
		assert.Equal(t, expected, err)
	})
	t.Run("archive interview appointment error when query fail", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		id := "64aaf0156999249a602ff55f"
		objId, _ := primitive.ObjectIDFromHex(id)
		expected := helpers.InternalError
		tsvc.interviewAppointmentRepo.On("ArchiveInterviewAppointment", ctx, objId).Return(errors.New("some error"))
		err := tsvc.service.ArchiveInterviewAppointment(ctx, id)
		assert.Equal(t, expected, err)
	})
}

func TestAddInterviewComment(t *testing.T) {
	t.Run("add interview comment success", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		id := "64aaf0156999249a602ff55f"
		objId, _ := primitive.ObjectIDFromHex(id)
		userId := "6476f457e64589e868aac97d"
		userObjId, _ := primitive.ObjectIDFromHex(userId)
		req := &dto.AddInterviewCommentRequest{
			ID:      id,
			Comment: "comment",
			UserID:  userId,
		}
		params := &domains.AddInterviewCommentParams{
			ID:      objId,
			Comment: req.Comment,
			UserID:  userObjId,
		}
		tsvc.interviewAppointmentRepo.On("AddComment", ctx, params).Return(nil)
		err := tsvc.service.AddInterviewComment(ctx, req)
		assert.NoError(t, err)
	})
	t.Run("add interview comment error when invalid interview appointment id format", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		id := "xxxxxxxx"
		userId := "6476f457e64589e868aac97d"
		req := &dto.AddInterviewCommentRequest{
			ID:      id,
			Comment: "comment",
			UserID:  userId,
		}
		expected := helpers.InternalError
		err := tsvc.service.AddInterviewComment(ctx, req)
		assert.Equal(t, expected, err)
	})
	t.Run("add interview comment error when invalid user id format", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		id := "64aaf0156999249a602ff55f"
		userId := "xxxxxxx"
		req := &dto.AddInterviewCommentRequest{
			ID:      id,
			Comment: "comment",
			UserID:  userId,
		}
		expected := helpers.InternalError
		err := tsvc.service.AddInterviewComment(ctx, req)
		assert.Equal(t, expected, err)
	})
	t.Run("add interview comment error when interview appointment not found", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		id := "64aaf0156999249a602ff55f"
		objId, _ := primitive.ObjectIDFromHex(id)
		userId := "6476f457e64589e868aac97d"
		userObjId, _ := primitive.ObjectIDFromHex(userId)
		req := &dto.AddInterviewCommentRequest{
			ID:      id,
			Comment: "comment",
			UserID:  userId,
		}
		params := &domains.AddInterviewCommentParams{
			ID:      objId,
			Comment: req.Comment,
			UserID:  userObjId,
		}
		expected := helpers.NewCustomError(http.StatusNotFound, "Interview appointment not found.")
		tsvc.interviewAppointmentRepo.On("AddComment", ctx, params).Return(mongo.ErrNoDocuments)
		err := tsvc.service.AddInterviewComment(ctx, req)
		assert.Equal(t, expected, err)
	})
	t.Run("add interview comment error when query fail", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		id := "64aaf0156999249a602ff55f"
		objId, _ := primitive.ObjectIDFromHex(id)
		userId := "6476f457e64589e868aac97d"
		userObjId, _ := primitive.ObjectIDFromHex(userId)
		req := &dto.AddInterviewCommentRequest{
			ID:      id,
			Comment: "comment",
			UserID:  userId,
		}
		params := &domains.AddInterviewCommentParams{
			ID:      objId,
			Comment: req.Comment,
			UserID:  userObjId,
		}
		expected := helpers.InternalError
		tsvc.interviewAppointmentRepo.On("AddComment", ctx, params).Return(errors.New("some error"))
		err := tsvc.service.AddInterviewComment(ctx, req)
		assert.Equal(t, expected, err)
	})
}

func TestUpdateInterviewComment(t *testing.T) {
	t.Run("update interview comment success", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		id := "6476f457e64589e868aac993"
		commentId := "6476f457e64589e868aac996"
		userId := "6476f457e64589e868aac997"
		objId, _ := primitive.ObjectIDFromHex(id)
		commentObjId, _ := primitive.ObjectIDFromHex(commentId)
		userObjId, _ := primitive.ObjectIDFromHex(userId)
		interviewAppointment := mockInterviewAppointment1
		interviewAppointment.Comments = append(interviewAppointment.Comments, []domains.InterviewComment{
			{
				ID:      commentObjId,
				Comment: "comment",
				User: domains.User{
					ID:       userObjId,
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
		req := &dto.UpdateInterviewCommentRequest{
			ID:        id,
			CommentID: commentId,
			Comment:   "Update comment",
			UserID:    userId,
		}
		params := &domains.UpdateInterviewCommentParams{
			ID:        objId,
			CommentID: commentObjId,
			Comment:   req.Comment,
		}
		tsvc.interviewAppointmentRepo.On("Get", ctx, objId).Return(&interviewAppointment, nil)
		tsvc.interviewAppointmentRepo.On("UpdateComment", ctx, params).Return(nil)
		err := tsvc.service.UpdateInterviewComment(ctx, req)
		assert.NoError(t, err)
	})
	t.Run("update interview comment error when invalid interview appointment id format", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		id := "xxxxxxx"
		commentId := "6476f457e64589e868aac996"
		userId := "6476f457e64589e868aac997"
		req := &dto.UpdateInterviewCommentRequest{
			ID:        id,
			CommentID: commentId,
			Comment:   "Update comment",
			UserID:    userId,
		}
		expected := helpers.InternalError
		err := tsvc.service.UpdateInterviewComment(ctx, req)
		assert.Equal(t, expected, err)
	})
	t.Run("update interview comment error when invalid interview comment id format", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		id := "6476f457e64589e868aac993"
		commentId := "xxxxxxx"
		userId := "6476f457e64589e868aac997"
		req := &dto.UpdateInterviewCommentRequest{
			ID:        id,
			CommentID: commentId,
			Comment:   "Update comment",
			UserID:    userId,
		}
		expected := helpers.InternalError
		err := tsvc.service.UpdateInterviewComment(ctx, req)
		assert.Equal(t, expected, err)
	})
	t.Run("update interview comment error when get interview appointment fail", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		id := "6476f457e64589e868aac993"
		commentId := "6476f457e64589e868aac996"
		userId := "6476f457e64589e868aac997"
		objId, _ := primitive.ObjectIDFromHex(id)
		req := &dto.UpdateInterviewCommentRequest{
			ID:        id,
			CommentID: commentId,
			Comment:   "Update comment",
			UserID:    userId,
		}
		expected := helpers.InternalError
		tsvc.interviewAppointmentRepo.On("Get", ctx, objId).Return(nil, errors.New("some error"))
		err := tsvc.service.UpdateInterviewComment(ctx, req)
		assert.Equal(t, expected, err)
	})
	t.Run("update interview comment error when interview appointment not found", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		id := "6476f457e64589e868aac993"
		commentId := "6476f457e64589e868aac996"
		userId := "6476f457e64589e868aac997"
		objId, _ := primitive.ObjectIDFromHex(id)
		req := &dto.UpdateInterviewCommentRequest{
			ID:        id,
			CommentID: commentId,
			Comment:   "Update comment",
			UserID:    userId,
		}
		expected := helpers.NewCustomError(http.StatusNotFound, "Interview appointment not found.")
		tsvc.interviewAppointmentRepo.On("Get", ctx, objId).Return(nil, nil)
		err := tsvc.service.UpdateInterviewComment(ctx, req)
		assert.Equal(t, expected, err)
	})
	t.Run("update interview comment error when interview comment not found", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		id := "6476f457e64589e868aac993"
		commentId := "6476f457e64589e868aac996"
		userId := "6476f457e64589e868aac997"
		objId, _ := primitive.ObjectIDFromHex(id)
		interviewAppointment := mockInterviewAppointment1
		interviewAppointment.Comments = []domains.InterviewComment{}
		req := &dto.UpdateInterviewCommentRequest{
			ID:        id,
			CommentID: commentId,
			Comment:   "Update comment",
			UserID:    userId,
		}
		expected := helpers.NewCustomError(http.StatusNotFound, "Interview comment not found.")
		tsvc.interviewAppointmentRepo.On("Get", ctx, objId).Return(&interviewAppointment, nil)
		err := tsvc.service.UpdateInterviewComment(ctx, req)
		assert.Equal(t, expected, err)
	})
	t.Run("update interview comment error when user comment not match", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		id := "6476f457e64589e868aac993"
		commentId := "6476f457e64589e868aac996"
		userId := "6476f457e64589e868aac997"
		objId, _ := primitive.ObjectIDFromHex(id)
		commentObjId, _ := primitive.ObjectIDFromHex(commentId)
		interviewAppointment := mockInterviewAppointment1
		interviewAppointment.Comments = append(interviewAppointment.Comments, []domains.InterviewComment{
			{
				ID:      commentObjId,
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
		req := &dto.UpdateInterviewCommentRequest{
			ID:        id,
			CommentID: commentId,
			Comment:   "Update comment",
			UserID:    userId,
		}
		expected := helpers.NewCustomError(http.StatusForbidden, "You don't have permission to update this comment")
		tsvc.interviewAppointmentRepo.On("Get", ctx, objId).Return(&interviewAppointment, nil)
		err := tsvc.service.UpdateInterviewComment(ctx, req)
		assert.Equal(t, expected, err)
	})
	t.Run("update interview comment error when update query fail", func(t *testing.T) {
		tsvc := newTestInterviewService(t)
		id := "6476f457e64589e868aac993"
		commentId := "6476f457e64589e868aac996"
		userId := "6476f457e64589e868aac997"
		objId, _ := primitive.ObjectIDFromHex(id)
		commentObjId, _ := primitive.ObjectIDFromHex(commentId)
		userObjId, _ := primitive.ObjectIDFromHex(userId)
		interviewAppointment := mockInterviewAppointment1
		interviewAppointment.Comments = append(interviewAppointment.Comments, []domains.InterviewComment{
			{
				ID:      commentObjId,
				Comment: "comment",
				User: domains.User{
					ID:       userObjId,
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
		req := &dto.UpdateInterviewCommentRequest{
			ID:        id,
			CommentID: commentId,
			Comment:   "Update comment",
			UserID:    userId,
		}
		params := &domains.UpdateInterviewCommentParams{
			ID:        objId,
			CommentID: commentObjId,
			Comment:   req.Comment,
		}
		expected := helpers.InternalError
		tsvc.interviewAppointmentRepo.On("Get", ctx, objId).Return(&interviewAppointment, nil)
		tsvc.interviewAppointmentRepo.On("UpdateComment", ctx, params).Return(errors.New("some error"))
		err := tsvc.service.UpdateInterviewComment(ctx, req)
		assert.Equal(t, expected, err)
	})
}
