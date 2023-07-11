package repositories_test

import (
	"context"
	"fmt"
	"robinhood-assignment/internal/core/domains"
	"robinhood-assignment/internal/core/ports"
	"robinhood-assignment/internal/repositories"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

type testInterviewAppointmentRepository struct {
	interviewRepo ports.InterviewAppointmentRepository
}

func newTestInterviewAppointmentRepository(mc *mongo.Client, db string) testInterviewAppointmentRepository {
	interviewRepo := repositories.NewInterviewAppointmentRepository(mc, db)
	return testInterviewAppointmentRepository{interviewRepo}
}

var (
	ctx                       = context.Background()
	dbName                    = "interview"
	collectionName            = "interviewAppointment"
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
	}
)

func TestGetAll(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("get all success", func(mt *mtest.T) {
		trepo := newTestInterviewAppointmentRepository(mt.Client, dbName)
		first := mtest.CreateCursorResponse(1, fmt.Sprintf("%s.%s", dbName, collectionName), mtest.FirstBatch, bson.D{
			{Key: "_id", Value: mockInterviewAppointment1.ID},
			{Key: "title", Value: mockInterviewAppointment1.Title},
			{Key: "description", Value: mockInterviewAppointment1.Description},
			{Key: "comments", Value: bson.A{}},
			{Key: "status", Value: mockInterviewAppointment1.Status},
			{Key: "isArchived", Value: mockInterviewAppointment1.IsArchived},
			{Key: "createUser", Value: mockInterviewAppointment1.CreateUser},
		})
		second := mtest.CreateCursorResponse(1, fmt.Sprintf("%s.%s", dbName, collectionName), mtest.NextBatch, bson.D{
			{Key: "_id", Value: mockInterviewAppointment2.ID},
			{Key: "title", Value: mockInterviewAppointment2.Title},
			{Key: "description", Value: mockInterviewAppointment2.Description},
			{Key: "comments", Value: bson.A{}},
			{Key: "status", Value: mockInterviewAppointment2.Status},
			{Key: "isArchived", Value: mockInterviewAppointment2.IsArchived},
			{Key: "createUser", Value: mockInterviewAppointment2.CreateUser},
		})
		killCursors := mtest.CreateCursorResponse(0, fmt.Sprintf("%s.%s", dbName, collectionName), mtest.NextBatch)
		mt.AddMockResponses(first, second, killCursors)
		data, err := trepo.interviewRepo.GetAll(ctx, 0, 20)
		assert.Nil(t, err)
		assert.Equal(t, []domains.InterviewAppointment{
			mockInterviewAppointment1,
			mockInterviewAppointment2,
		}, data)
	})
	mt.Run("get all error", func(mt *mtest.T) {
		trepo := newTestInterviewAppointmentRepository(mt.Client, dbName)
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "duplicate key error",
		}))
		data, err := trepo.interviewRepo.GetAll(ctx, 0, 20)
		assert.Error(t, err)
		assert.Equal(t, []domains.InterviewAppointment{}, data)
	})
}
func TestGet(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("get success", func(mt *mtest.T) {
		trepo := newTestInterviewAppointmentRepository(mt.Client, dbName)
		expected := &mockInterviewAppointment1
		first := mtest.CreateCursorResponse(1, fmt.Sprintf("%s.%s", dbName, collectionName), mtest.FirstBatch, bson.D{
			{Key: "_id", Value: mockInterviewAppointment1.ID},
			{Key: "title", Value: mockInterviewAppointment1.Title},
			{Key: "description", Value: mockInterviewAppointment1.Description},
			{Key: "comments", Value: bson.A{}},
			{Key: "status", Value: mockInterviewAppointment1.Status},
			{Key: "isArchived", Value: mockInterviewAppointment1.IsArchived},
			{Key: "createUser", Value: mockInterviewAppointment1.CreateUser},
		})
		killCursors := mtest.CreateCursorResponse(0, fmt.Sprintf("%s.%s", dbName, collectionName), mtest.NextBatch)
		mt.AddMockResponses(first, killCursors)
		data, err := trepo.interviewRepo.Get(ctx, mockInterviewAppointment1.ID)
		assert.Nil(t, err)
		assert.Equal(t, expected, data)
	})
	mt.Run("get error", func(mt *mtest.T) {
		trepo := newTestInterviewAppointmentRepository(mt.Client, dbName)
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "duplicate key error",
		}))
		data, err := trepo.interviewRepo.Get(ctx, mockInterviewAppointment1.ID)
		assert.Error(t, err)
		assert.Nil(t, data)
	})
}

func TestCreate(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("create success", func(mt *mtest.T) {
		trepo := newTestInterviewAppointmentRepository(mt.Client, dbName)
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		params := &domains.CreateInterviewAppointmentParams{
			Title:       mockInterviewAppointment1.Title,
			Description: mockInterviewAppointment1.Description,
			UserID:      mockInterviewAppointment1.CreateUser.ID,
		}
		data, err := trepo.interviewRepo.Create(ctx, params)
		assert.Nil(t, err)
		assert.Equal(t, params.Title, data.Title)
		assert.Equal(t, params.Description, data.Description)
		assert.Equal(t, params.UserID, data.CreateUserId)
	})
	mt.Run("create error", func(mt *mtest.T) {
		trepo := newTestInterviewAppointmentRepository(mt.Client, dbName)
		params := &domains.CreateInterviewAppointmentParams{
			Title:       mockInterviewAppointment1.Title,
			Description: mockInterviewAppointment1.Description,
			UserID:      mockInterviewAppointment1.CreateUser.ID,
		}
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "duplicate key error",
		}))
		data, err := trepo.interviewRepo.Create(ctx, params)
		assert.Error(t, err)
		assert.Nil(t, data)
	})
}

func TestUpdate(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("update success", func(mt *mtest.T) {
		trepo := newTestInterviewAppointmentRepository(mt.Client, dbName)
		params := &domains.UpdateInterviewAppointmentParams{
			ID:          mockInterviewAppointment1.ID,
			Title:       mockInterviewAppointment1.Title,
			Description: mockInterviewAppointment1.Description,
			Status:      mockInterviewAppointment1.Status,
		}
		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "value", Value: bson.D{
				{Key: "_id", Value: params.ID},
				{Key: "title", Value: params.Title},
				{Key: "description", Value: params.Description},
				{Key: "status", Value: params.Status},
			}},
		})
		data, err := trepo.interviewRepo.Update(ctx, params)
		assert.Nil(t, err)
		assert.Equal(t, params.Title, data.Title)
		assert.Equal(t, params.Description, data.Description)
		assert.Equal(t, params.Status, data.Status)
	})
	mt.Run("update error", func(mt *mtest.T) {
		trepo := newTestInterviewAppointmentRepository(mt.Client, dbName)
		params := &domains.UpdateInterviewAppointmentParams{
			ID:          mockInterviewAppointment1.ID,
			Title:       mockInterviewAppointment1.Title,
			Description: mockInterviewAppointment1.Description,
			Status:      mockInterviewAppointment1.Status,
		}
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "update fail",
		}))
		data, err := trepo.interviewRepo.Update(ctx, params)
		assert.Error(t, err)
		assert.Nil(t, data)
	})
}

func TestArchiveInterviewAppointment(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("archive interview appointment success", func(mt *mtest.T) {
		trepo := newTestInterviewAppointmentRepository(mt.Client, dbName)
		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "value", Value: bson.D{
				{Key: "_id", Value: mockInterviewAppointment1.ID},
				{Key: "title", Value: mockInterviewAppointment1.Title},
				{Key: "description", Value: mockInterviewAppointment1.Description},
				{Key: "status", Value: mockInterviewAppointment1.Status},
				{Key: "isArchived", Value: true},
			}},
		})
		err := trepo.interviewRepo.ArchiveInterviewAppointment(ctx, mockInterviewAppointment1.ID)
		assert.NoError(t, err)
	})
	mt.Run("archive interview appointment success", func(mt *mtest.T) {
		trepo := newTestInterviewAppointmentRepository(mt.Client, dbName)
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "update fail",
		}))
		err := trepo.interviewRepo.ArchiveInterviewAppointment(ctx, mockInterviewAppointment1.ID)
		assert.Error(t, err)
	})
}

func TestAddComment(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("add comment success", func(mt *mtest.T) {
		trepo := newTestInterviewAppointmentRepository(mt.Client, dbName)
		params := &domains.AddInterviewCommentParams{
			ID:      mockInterviewAppointment1.ID,
			Comment: "add comment",
			UserID:  mockInterviewAppointment1.CreateUser.ID,
		}
		commentId := primitive.NewObjectID()
		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "value", Value: bson.D{
				{Key: "_id", Value: mockInterviewAppointment1.ID},
				{Key: "title", Value: mockInterviewAppointment1.Title},
				{Key: "description", Value: mockInterviewAppointment1.Description},
				{Key: "comments", Value: bson.A{
					bson.D{
						{Key: "_id", Value: commentId},
						{Key: "comment", Value: params.Comment},
						{Key: "userId", Value: params.UserID},
					},
				}},
				{Key: "status", Value: mockInterviewAppointment1.Status},
				{Key: "isArchived", Value: false},
			}},
		})
		err := trepo.interviewRepo.AddComment(ctx, params)
		assert.NoError(t, err)
	})
	mt.Run("add comment error", func(mt *mtest.T) {
		trepo := newTestInterviewAppointmentRepository(mt.Client, dbName)
		params := &domains.AddInterviewCommentParams{
			ID:      mockInterviewAppointment1.ID,
			Comment: "add comment",
			UserID:  mockInterviewAppointment1.CreateUser.ID,
		}
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "update fail",
		}))
		err := trepo.interviewRepo.AddComment(ctx, params)
		assert.Error(t, err)
	})
}

func TestUpdateComment(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("update comment success", func(mt *mtest.T) {
		trepo := newTestInterviewAppointmentRepository(mt.Client, dbName)
		params := &domains.UpdateInterviewCommentParams{
			ID:        mockInterviewAppointment1.ID,
			Comment:   "update comment",
			CommentID: primitive.NewObjectID(),
		}
		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "value", Value: bson.D{
				{Key: "_id", Value: mockInterviewAppointment1.ID},
				{Key: "title", Value: mockInterviewAppointment1.Title},
				{Key: "description", Value: mockInterviewAppointment1.Description},
				{Key: "comments", Value: bson.A{
					bson.D{
						{Key: "_id", Value: params.CommentID},
						{Key: "comment", Value: params.Comment},
					},
				}},
				{Key: "status", Value: mockInterviewAppointment1.Status},
				{Key: "isArchived", Value: false},
			}},
		})
		err := trepo.interviewRepo.UpdateComment(ctx, params)
		assert.NoError(t, err)
	})
	mt.Run("update comment error", func(mt *mtest.T) {
		trepo := newTestInterviewAppointmentRepository(mt.Client, dbName)
		params := &domains.UpdateInterviewCommentParams{
			ID:        mockInterviewAppointment1.ID,
			Comment:   "update comment",
			CommentID: primitive.NewObjectID(),
		}
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "update fail",
		}))
		err := trepo.interviewRepo.UpdateComment(ctx, params)
		assert.Error(t, err)
	})
}
