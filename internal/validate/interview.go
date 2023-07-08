package validate

import (
	"net/http"
	"robinhood-assignment/helpers"
	"robinhood-assignment/internal/core/ports"
	"robinhood-assignment/internal/dto"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

type interviewValidate struct {
}

func NewInterviewValidate() ports.InterviewValidate {
	return &interviewValidate{}
}

func (v interviewValidate) ValidateGetInterviewAppointments(ctx *gin.Context) (*dto.GetInterviewAppointmentsRequest, error) {
	req := dto.GetInterviewAppointmentsRequest{}
	if page, ok := ctx.GetQuery("page"); ok {
		v, err := strconv.Atoi(page)
		if err != nil {
			return nil, helpers.NewCustomError(http.StatusBadRequest, "Invalid page query parameter")
		}
		i := uint64(v)
		req.Page = i
	}
	if limit, ok := ctx.GetQuery("limit"); ok {
		v, err := strconv.Atoi(limit)
		if err != nil {
			return nil, helpers.NewCustomError(http.StatusBadRequest, "Invalid limit query parameter")
		}
		i := uint64(v)
		req.Limit = i
	}
	if _, err := govalidator.ValidateStruct(req); err != nil {
		return nil, helpers.NewCustomError(http.StatusBadRequest, err.Error())
	}
	return &req, nil
}

func (v interviewValidate) ValidateGetInterviewAppointment(ctx *gin.Context) (string, error) {
	id := ctx.Param("id")
	if id == "" {
		return "", helpers.NewCustomError(http.StatusBadRequest, "id: Missing required field")
	}
	formats := strfmt.Default
	if err := validate.FormatOf("id", "param", "bsonobjectid", id, formats); err != nil {
		return "", helpers.NewCustomError(http.StatusBadRequest, err.Error())
	}
	return id, nil
}

func (v interviewValidate) ValidateCreateInterviewAppointment(ctx *gin.Context) (*dto.CreateInterviewAppointmentRequest, error) {
	req := dto.CreateInterviewAppointmentRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		return nil, helpers.NewCustomError(http.StatusBadRequest, "Invalid input parameter")
	}
	if _, err := govalidator.ValidateStruct(req); err != nil {
		return nil, helpers.NewCustomError(http.StatusBadRequest, err.Error())
	}
	return &req, nil
}
