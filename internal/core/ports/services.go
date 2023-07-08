package ports

import (
	"context"
	"robinhood-assignment/internal/core/domains"
	"robinhood-assignment/internal/dto"
)

type AuthServie interface {
	RegisterAdmin(ctx context.Context, params *dto.RegisterAdminRequest) (*dto.RegisterAdminResponse, error)
	Login(ctx context.Context, params *dto.LoginRequest) (*dto.LoginResponse, error)
}

type InterviewService interface {
	GetInterviewAppointments(ctx context.Context, offset int64, limit int64) ([]domains.InterviewAppointment, error)
	GetInterviewAppointment(ctx context.Context, id string) (*domains.InterviewAppointment, error)
	CreateInterviewAppointment(ctx context.Context, req *dto.CreateInterviewAppointmentRequest) (*domains.InterviewAppointment, error)
	// UpdateInterviewAppointment(ctx *gin.Context)
	// ArchiveInterviewAppointment(ctx *gin.Context)
}
