package services

import (
	"context"
	"fmt"
	"net/http"
	"robinhood-assignment/config"
	"robinhood-assignment/helpers"
	"robinhood-assignment/internal/core/domains"
	"robinhood-assignment/internal/core/ports"
	"robinhood-assignment/internal/dto"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type authService struct {
	userRepo ports.UserRepository
	myBcrypt ports.MyBcrypt
	myJWT    ports.MyJWT
}

func NewAuthService(userRepo ports.UserRepository, myBcrypt ports.MyBcrypt, myJWT ports.MyJWT) ports.AuthServie {
	return &authService{userRepo, myBcrypt, myJWT}
}

func (a *authService) RegisterAdmin(ctx context.Context, req *dto.RegisterAdminRequest) (*dto.RegisterAdminResponse, error) {
	user, err := a.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, helpers.InternalError
	}
	if user != nil {
		return nil, helpers.CustomError{
			StatusCode: http.StatusConflict,
			Message:    "Duplicate username",
		}
	}
	passHash, err := a.myBcrypt.GenerateFromPassword(req.Password, config.Get().Auth.BcryptCost)
	if err != nil {
		return nil, helpers.InternalError
	}

	createUserPatams := &domains.CreateUserParams{
		Name:     req.Name,
		Email:    req.Email,
		Username: req.Username,
		Password: *passHash,
		ImageUrl: req.ImageUrl,
		Role:     req.Role,
	}
	if _, err := a.userRepo.Create(ctx, createUserPatams); err != nil {
		fmt.Printf("err: %#v\n", err)
		return nil, helpers.CustomError{
			StatusCode: http.StatusConflict,
			Message:    "Create user fail.",
		}
	}
	return &dto.RegisterAdminResponse{
		StatusCode: http.StatusCreated,
		Message:    "Register admin success.",
	}, nil
}

func (a *authService) Login(ctx context.Context, params *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := a.userRepo.GetByUsername(ctx, params.Username)
	if err != nil {
		return nil, helpers.InternalError
	}
	if user == nil {
		return nil, helpers.CustomError{
			StatusCode: http.StatusNotFound,
			Message:    "Username not found",
		}
	}
	if err := a.myBcrypt.CompareHashAndPassword(user.Password, params.Password); err != nil {
		return nil, helpers.CustomError{
			StatusCode: http.StatusUnauthorized,
			Message:    "Password is incorrect",
		}
	}
	token := a.myJWT.NewWithClaims(jwt.SigningMethodHS256, domains.Claims{
		UserID: user.ID.Hex(),
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
		},
	})
	tokenString, err := token.SignedString([]byte(config.Get().Auth.JwtSecret))
	if err != nil {
		return nil, helpers.InternalError
	}
	return &dto.LoginResponse{
		StatusCode: http.StatusOK,
		Token:      tokenString,
	}, nil
}
