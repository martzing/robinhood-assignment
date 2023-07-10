package services

import (
	"context"
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

func (a *authService) CreateStaff(ctx context.Context, req *dto.CreateStaffRequest) error {
	user, err := a.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return helpers.InternalError
	}
	if user != nil {
		return helpers.NewCustomError(http.StatusConflict, "Duplicate username")
	}
	passHash, err := a.myBcrypt.GenerateFromPassword(req.Password, config.Get().Auth.BcryptCost)
	if err != nil {
		return helpers.InternalError
	}

	params := &domains.CreateUserParams{
		Name:     req.Name,
		Email:    req.Email,
		Username: req.Username,
		Password: *passHash,
		ImageUrl: req.ImageUrl,
		Role:     req.Role,
	}
	if _, err := a.userRepo.Create(ctx, params); err != nil {
		return helpers.NewCustomError(http.StatusConflict, "Create staff fail")
	}
	return nil
}

func (a *authService) Login(ctx context.Context, req *dto.LoginRequest) (string, error) {
	user, err := a.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return "", helpers.InternalError
	}
	if user == nil {
		return "", helpers.NewCustomError(http.StatusNotFound, "Username not found")
	}
	if err := a.myBcrypt.CompareHashAndPassword(user.Password, req.Password); err != nil {
		return "", helpers.NewCustomError(http.StatusNotFound, "Password is incorrect")
	}
	token := a.myJWT.NewWithClaims(jwt.SigningMethodHS256, domains.Claims{
		UserID: user.ID.Hex(),
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute)),
		},
	})
	tokenString, err := token.SignedString([]byte(config.Get().Auth.JwtSecret))
	if err != nil {
		return "", helpers.InternalError
	}
	return tokenString, nil
}
