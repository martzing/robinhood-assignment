package services_test

import (
	"errors"
	"net/http"
	"robinhood-assignment/config"
	"robinhood-assignment/helpers"
	"robinhood-assignment/internal/core/constants"
	"robinhood-assignment/internal/core/domains"
	"robinhood-assignment/internal/core/ports"
	"robinhood-assignment/internal/core/ports/mocks"
	"robinhood-assignment/internal/core/services"
	"robinhood-assignment/internal/dto"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type testAuthService struct {
	userRepo *mocks.UserRepository
	myBcrypt *mocks.MyBcrypt
	myJWT    *mocks.MyJWT
	service  ports.AuthServie
}

func newTestAuthService(t *testing.T) testAuthService {
	userRepo := mocks.NewUserRepository(t)
	myBcrypt := mocks.NewMyBcrypt(t)
	myJWT := mocks.NewMyJWT(t)

	service := services.NewAuthService(userRepo, myBcrypt, myJWT)
	return testAuthService{userRepo, myBcrypt, myJWT, service}
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

func TestCreateStaff(t *testing.T) {
	t.Setenv("BCRYPT_COST", "8")
	config.New()
	t.Run("create staff success", func(t *testing.T) {
		tsvc := newTestAuthService(t)
		passHash := "mockhashpassword"
		req := &dto.CreateStaffRequest{
			Name:     name,
			Email:    email,
			Username: username,
			Password: password,
			ImageUrl: imageUrl,
			Role:     constants.STAFF_ROLE,
		}
		params := &domains.CreateUserParams{
			Name:     req.Name,
			Email:    req.Email,
			Username: req.Username,
			Password: passHash,
			ImageUrl: req.ImageUrl,
			Role:     req.Role,
		}
		tsvc.userRepo.On("GetByUsername", ctx, username).Return(nil, nil)
		tsvc.myBcrypt.On("GenerateFromPassword", password, bcryptCost).Return(&passHash, nil)
		tsvc.userRepo.On("Create", ctx, params).Return(&user, nil)
		err := tsvc.service.CreateStaff(ctx, req)
		assert.NoError(t, err)
	})
	t.Run("create staff error when input duplicate username", func(t *testing.T) {
		tsvc := newTestAuthService(t)
		req := &dto.CreateStaffRequest{
			Name:     name,
			Email:    email,
			Username: username,
			Password: password,
			ImageUrl: imageUrl,
			Role:     constants.STAFF_ROLE,
		}
		expected := helpers.NewCustomError(http.StatusConflict, "Duplicate username")
		tsvc.userRepo.On("GetByUsername", ctx, username).Return(&user, nil)
		err := tsvc.service.CreateStaff(ctx, req)
		assert.Equal(t, expected, err)
	})
	t.Run("create staff error when get user fail", func(t *testing.T) {
		tsvc := newTestAuthService(t)
		req := &dto.CreateStaffRequest{
			Name:     name,
			Email:    email,
			Username: username,
			Password: password,
			ImageUrl: imageUrl,
			Role:     constants.STAFF_ROLE,
		}
		expected := helpers.InternalError
		tsvc.userRepo.On("GetByUsername", ctx, username).Return(nil, errors.New("some error"))
		err := tsvc.service.CreateStaff(ctx, req)
		assert.Equal(t, expected, err)
	})
	t.Run("create staff error when generate password hash fail", func(t *testing.T) {
		tsvc := newTestAuthService(t)
		req := &dto.CreateStaffRequest{
			Name:     name,
			Email:    email,
			Username: username,
			Password: password,
			ImageUrl: imageUrl,
			Role:     constants.STAFF_ROLE,
		}
		expected := helpers.InternalError
		tsvc.userRepo.On("GetByUsername", ctx, username).Return(nil, nil)
		tsvc.myBcrypt.On("GenerateFromPassword", password, bcryptCost).Return(nil, errors.New("some error"))
		err := tsvc.service.CreateStaff(ctx, req)
		assert.Equal(t, expected, err)
	})
	t.Run("create staff error when query fail", func(t *testing.T) {
		tsvc := newTestAuthService(t)
		passHash := "mockhashpassword"
		req := &dto.CreateStaffRequest{
			Name:     name,
			Email:    email,
			Username: username,
			Password: password,
			ImageUrl: imageUrl,
			Role:     constants.STAFF_ROLE,
		}
		params := &domains.CreateUserParams{
			Name:     req.Name,
			Email:    req.Email,
			Username: req.Username,
			Password: passHash,
			ImageUrl: req.ImageUrl,
			Role:     req.Role,
		}
		expected := helpers.NewCustomError(http.StatusConflict, "Create staff fail")
		tsvc.userRepo.On("GetByUsername", ctx, username).Return(nil, nil)
		tsvc.myBcrypt.On("GenerateFromPassword", password, bcryptCost).Return(&passHash, nil)
		tsvc.userRepo.On("Create", ctx, params).Return(nil, errors.New("some error"))
		err := tsvc.service.CreateStaff(ctx, req)
		assert.Equal(t, expected, err)
	})
}

func TestLogin(t *testing.T) {
	t.Setenv("BCRYPT_COST", "8")
	t.Setenv("JWT_SECRET", "mock-jwt-secret")
	config.New()
	t.Run("login success", func(t *testing.T) {
		tsvc := newTestAuthService(t)
		req := &dto.LoginRequest{
			Username: username,
			Password: password,
		}
		claims := domains.Claims{
			UserID: user.ID.Hex(),
			Role:   user.Role,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute)),
			},
		}
		token := jwt.New(jwt.SigningMethodHS256)
		expected, err := token.SignedString([]byte(config.Get().Auth.JwtSecret))

		tsvc.userRepo.On("GetByUsername", ctx, username).Return(&user, nil)
		tsvc.myBcrypt.On("CompareHashAndPassword", user.Password, password).Return(nil)
		tsvc.myJWT.On("NewWithClaims", jwt.SigningMethodHS256, claims).Return(token, nil)
		got, err := tsvc.service.Login(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})
	t.Run("login error when username not found", func(t *testing.T) {
		tsvc := newTestAuthService(t)
		req := &dto.LoginRequest{
			Username: username,
			Password: password,
		}
		expectedRes := ""
		expectedErr := helpers.NewCustomError(http.StatusNotFound, "Username not found")

		tsvc.userRepo.On("GetByUsername", ctx, username).Return(nil, nil)
		res, err := tsvc.service.Login(ctx, req)
		assert.Equal(t, expectedRes, res)
		assert.Equal(t, expectedErr, err)
	})
	t.Run("login error when get user fail", func(t *testing.T) {
		tsvc := newTestAuthService(t)
		req := &dto.LoginRequest{
			Username: username,
			Password: password,
		}
		expectedRes := ""
		expectedErr := helpers.InternalError

		tsvc.userRepo.On("GetByUsername", ctx, username).Return(nil, errors.New("some error"))
		res, err := tsvc.service.Login(ctx, req)
		assert.Equal(t, expectedRes, res)
		assert.Equal(t, expectedErr, err)
	})
	t.Run("login error when compare password fail", func(t *testing.T) {
		tsvc := newTestAuthService(t)
		req := &dto.LoginRequest{
			Username: username,
			Password: password,
		}
		expectedRes := ""
		expectedErr := helpers.NewCustomError(http.StatusUnauthorized, "Password is incorrect")

		tsvc.userRepo.On("GetByUsername", ctx, username).Return(&user, nil)
		tsvc.myBcrypt.On("CompareHashAndPassword", user.Password, password).Return(errors.New("some error"))
		res, err := tsvc.service.Login(ctx, req)
		assert.Equal(t, expectedRes, res)
		assert.Equal(t, expectedErr, err)
	})
	t.Run("call NewWithClaims func with correct parameter", func(t *testing.T) {
		tsvc := newTestAuthService(t)
		req := &dto.LoginRequest{
			Username: username,
			Password: password,
		}
		claims := domains.Claims{
			UserID: user.ID.Hex(),
			Role:   user.Role,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute)),
			},
		}
		token := jwt.New(jwt.SigningMethodHS256)

		tsvc.userRepo.On("GetByUsername", ctx, username).Return(&user, nil)
		tsvc.myBcrypt.On("CompareHashAndPassword", user.Password, password).Return(nil)
		tsvc.myJWT.On("NewWithClaims", jwt.SigningMethodHS256, claims).Return(token, nil)
		tsvc.service.Login(ctx, req)
		tsvc.userRepo.AssertCalled(t, "GetByUsername", ctx, username)
		tsvc.myBcrypt.AssertCalled(t, "CompareHashAndPassword", user.Password, password)
		tsvc.myJWT.AssertCalled(t, "NewWithClaims", jwt.SigningMethodHS256, claims)
	})
}
