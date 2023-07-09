package handlers

import (
	"net/http"
	"robinhood-assignment/helpers"
	"robinhood-assignment/internal/core/ports"
	"robinhood-assignment/internal/dto"

	"github.com/gin-gonic/gin"
)

type authHandler struct {
	authSvc  ports.AuthServie
	validate ports.AuthValidate
}

func NewAuthHandler(authSvc ports.AuthServie, validate ports.AuthValidate) ports.AuthHandler {
	return &authHandler{authSvc, validate}
}

func (a *authHandler) CreateStaff(ctx *gin.Context) {
	params, err := a.validate.ValidateCreateStaff(ctx)
	if err != nil {
		errRes := helpers.ErrorHandler(err)
		ctx.AbortWithStatusJSON(errRes.StatusCode, errRes)
		return
	}

	if err := a.authSvc.CreateStaff(ctx, params); err != nil {
		errRes := helpers.ErrorHandler(err)
		ctx.AbortWithStatusJSON(errRes.StatusCode, errRes)
		return
	}
	ctx.JSON(http.StatusCreated, dto.BaseResponse{
		StatusCode: http.StatusOK,
		Message:    "success",
	})
}

func (a *authHandler) Login(ctx *gin.Context) {
	params, err := a.validate.ValidateLogin(ctx)
	if err != nil {
		errRes := helpers.ErrorHandler(err)
		ctx.AbortWithStatusJSON(errRes.StatusCode, errRes)
		return
	}

	token, err := a.authSvc.Login(ctx, params)
	if err != nil {
		errRes := helpers.ErrorHandler(err)
		ctx.AbortWithStatusJSON(errRes.StatusCode, errRes)
		return
	}
	response := dto.LoginResponse{
		StatusCode: http.StatusOK,
		Token:      token,
	}
	ctx.JSON(http.StatusOK, response)
}
