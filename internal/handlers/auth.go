package handlers

import (
	"net/http"
	"robinhood-assignment/helpers"
	"robinhood-assignment/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type authHandler struct {
	authSvc  ports.AuthServie
	validate ports.AuthValidate
}

func NewAuthHandler(authSvc ports.AuthServie, validate ports.AuthValidate) ports.AuthHandler {
	return &authHandler{authSvc, validate}
}

func (a *authHandler) RegisterAdmin(ctx *gin.Context) {
	params, err := a.validate.ValidateRegisterAdmin(ctx)
	if err != nil {
		errRes := helpers.ErrorHandler(err)
		ctx.AbortWithStatusJSON(errRes.StatusCode, errRes)
		return
	}

	res, err := a.authSvc.RegisterAdmin(ctx, params)
	if err != nil {
		errRes := helpers.ErrorHandler(err)
		ctx.AbortWithStatusJSON(errRes.StatusCode, errRes)
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (a *authHandler) Login(ctx *gin.Context) {
	params, err := a.validate.ValidateLogin(ctx)
	if err != nil {
		errRes := helpers.ErrorHandler(err)
		ctx.AbortWithStatusJSON(errRes.StatusCode, errRes)
		return
	}

	res, err := a.authSvc.Login(ctx, params)
	if err != nil {
		errRes := helpers.ErrorHandler(err)
		ctx.AbortWithStatusJSON(errRes.StatusCode, errRes)
		return
	}
	ctx.JSON(http.StatusOK, res)
}
