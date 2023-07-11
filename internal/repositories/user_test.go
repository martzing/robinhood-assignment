package repositories_test

import (
	"fmt"
	"robinhood-assignment/internal/core/constants"
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

type testUserRepository struct {
	userRepo ports.UserRepository
}

func newTestUserRepository(mc *mongo.Client, db string) testUserRepository {
	userRepo := repositories.NewUserRepository(mc, db)
	return testUserRepository{userRepo}
}

var (
	name       = "Samart"
	email      = "samart.ph.work@gmail.com"
	username   = "samart"
	password   = "1234567890"
	imageUrl   = "https://image-url.com"
	bcryptCost = 8
	userId     = primitive.NewObjectID()
	user       = domains.User{
		ID:       userId,
		Name:     name,
		Email:    email,
		Username: username,
		Password: password,
		ImageUrl: imageUrl,
		Role:     constants.STAFF_ROLE,
	}
)

func TestGetUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("get user success", func(mt *mtest.T) {
		trepo := newTestUserRepository(mt.Client, dbName)
		expectedUser := user
		mt.AddMockResponses(mtest.CreateCursorResponse(1, fmt.Sprintf("%s.%s", dbName, collectionName), mtest.FirstBatch, bson.D{
			{Key: "_id", Value: user.ID},
			{Key: "name", Value: user.Name},
			{Key: "email", Value: user.Email},
			{Key: "username", Value: user.Username},
			{Key: "password", Value: user.Password},
			{Key: "imageUrl", Value: user.ImageUrl},
			{Key: "role", Value: user.Role},
		}))
		data, err := trepo.userRepo.Get(ctx, user.ID)
		assert.Nil(t, err)
		assert.Equal(t, &expectedUser, data)
	})
	mt.Run("get user error", func(mt *mtest.T) {
		trepo := newTestUserRepository(mt.Client, dbName)
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "not found document",
		}))
		data, err := trepo.userRepo.Get(ctx, user.ID)
		assert.Error(t, err)
		assert.Nil(t, data)
	})
}

func TestGetByUserUsername(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("get user by username success", func(mt *mtest.T) {
		trepo := newTestUserRepository(mt.Client, dbName)
		expectedUser := user
		mt.AddMockResponses(mtest.CreateCursorResponse(1, fmt.Sprintf("%s.%s", dbName, collectionName), mtest.FirstBatch, bson.D{
			{Key: "_id", Value: user.ID},
			{Key: "name", Value: user.Name},
			{Key: "email", Value: user.Email},
			{Key: "username", Value: user.Username},
			{Key: "password", Value: user.Password},
			{Key: "imageUrl", Value: user.ImageUrl},
			{Key: "role", Value: user.Role},
		}))
		data, err := trepo.userRepo.GetByUsername(ctx, user.Username)
		assert.Nil(t, err)
		assert.Equal(t, &expectedUser, data)
	})
	mt.Run("get user by username error", func(mt *mtest.T) {
		trepo := newTestUserRepository(mt.Client, dbName)
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "not found document",
		}))
		data, err := trepo.userRepo.GetByUsername(ctx, user.Username)
		assert.Error(t, err)
		assert.Nil(t, data)
	})
}

func TestCreateUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("create user success", func(mt *mtest.T) {
		trepo := newTestUserRepository(mt.Client, dbName)
		params := &domains.CreateUserParams{
			Name:     user.Name,
			Email:    user.Email,
			Username: user.Username,
			Password: user.Password,
			ImageUrl: user.ImageUrl,
			Role:     user.Role,
		}
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		data, err := trepo.userRepo.Create(ctx, params)
		assert.NoError(t, err)
		assert.Equal(t, params.Name, data.Name)
		assert.Equal(t, params.Email, data.Email)
		assert.Equal(t, params.Username, data.Username)
		assert.Equal(t, params.Password, data.Password)
		assert.Equal(t, params.ImageUrl, data.ImageUrl)
		assert.Equal(t, params.Role, data.Role)
	})
	mt.Run("create user error", func(mt *mtest.T) {
		trepo := newTestUserRepository(mt.Client, dbName)
		params := &domains.CreateUserParams{
			Name:     user.Name,
			Email:    user.Email,
			Username: user.Username,
			Password: user.Password,
			ImageUrl: user.ImageUrl,
			Role:     user.Role,
		}
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "duplicate key error",
		}))
		data, err := trepo.userRepo.Create(ctx, params)
		assert.Nil(t, data)
		assert.NotNil(t, err)
		assert.True(t, mongo.IsDuplicateKeyError(err))
	})
}
