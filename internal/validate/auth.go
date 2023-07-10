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
	req := &dto.LoginRequest{}
	if err := ctx.BindJSON(req); err != nil {
		return nil, helpers.NewCustomError(http.StatusBadRequest, "Invalid input parameter")
	}
	if _, err := govalidator.ValidateStruct(req); err != nil {
		return nil, helpers.NewCustomError(http.StatusBadRequest, err.Error())
	}
	return req, nil
}

func (v authValidate) ValidateCreateStaff(ctx *gin.Context) (*dto.CreateStaffRequest, error) {
	req := &dto.CreateStaffRequest{}
	if err := ctx.BindJSON(req); err != nil {
		return nil, helpers.NewCustomError(http.StatusBadRequest, "Invalid input parameter")
	}
	if _, err := govalidator.ValidateStruct(req); err != nil {
		return nil, helpers.NewCustomError(http.StatusBadRequest, err.Error())
	}
	return req, nil
}
