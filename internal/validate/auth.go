package validate

import (
	"net/http"
	"robinhood-assignment/helpers"
	"robinhood-assignment/internal/core/ports"
	"robinhood-assignment/internal/dto"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type authValidate struct {
}

func NewAuthValidate() ports.AuthValidate {
	return &authValidate{}
}

func (v authValidate) ValidateLogin(ctx *gin.Context) (*dto.LoginRequest, error) {
	params := &dto.LoginRequest{}
	if err := ctx.BindJSON(params); err != nil {
		return nil, helpers.NewCustomError(http.StatusBadRequest, "Invalid input parameter")
	}
	if _, err := govalidator.ValidateStruct(params); err != nil {
		return nil, helpers.NewCustomError(http.StatusBadRequest, err.Error())
	}
	return params, nil
}

func (v authValidate) ValidateRegisterAdmin(ctx *gin.Context) (*dto.RegisterAdminRequest, error) {
	params := &dto.RegisterAdminRequest{}
	if err := ctx.BindJSON(params); err != nil {
		return nil, helpers.NewCustomError(http.StatusBadRequest, "Invalid input parameter")
	}
	if _, err := govalidator.ValidateStruct(params); err != nil {
		return nil, helpers.NewCustomError(http.StatusBadRequest, err.Error())
	}
	return params, nil
}
