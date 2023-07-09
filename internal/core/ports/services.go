package ports

import (
	"context"
	"robinhood-assignment/internal/core/domains"
	"robinhood-assignment/internal/dto"
)

type AuthServie interface {
	CreateStaff(ctx context.Context, params *dto.CreateStaffRequest) error
	Login(ctx context.Context, params *dto.LoginRequest) (string, error)
}

type InterviewService interface {
	GetInterviewAppointments(ctx context.Context, offset int64, limit int64) ([]domains.InterviewAppointment, error)
	GetInterviewAppointment(ctx context.Context, id string) (*domains.InterviewAppointment, error)
	CreateInterviewAppointment(ctx context.Context, req *dto.CreateInterviewAppointmentRequest) (*domains.InterviewAppointment, error)
	UpdateInterviewAppointment(ctx context.Context, req *dto.UpdateInterviewAppointmentRequest) error
	ArchiveInterviewAppointment(ctx context.Context, id string) error
	AddInterviewComment(ctx context.Context, req *dto.AddInterviewCommentRequest) error
	UpdateInterviewComment(ctx context.Context, req *dto.UpdateInterviewCommentRequest) error
}
