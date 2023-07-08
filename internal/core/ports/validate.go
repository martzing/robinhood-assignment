package ports

import (
	"robinhood-assignment/internal/dto"

	"github.com/gin-gonic/gin"
)

type AuthValidate interface {
	ValidateLogin(ctx *gin.Context) (*dto.LoginRequest, error)
	ValidateRegisterAdmin(ctx *gin.Context) (*dto.RegisterAdminRequest, error)
}

type InterviewValidate interface {
	ValidateGetInterviewAppointments(ctx *gin.Context) (*dto.GetInterviewAppointmentsRequest, error)
	ValidateGetInterviewAppointment(ctx *gin.Context) (string, error)
	ValidateCreateInterviewAppointment(ctx *gin.Context) (*dto.CreateInterviewAppointmentRequest, error)
}
