package services

import (
	"context"
	"net/http"
	"robinhood-assignment/helpers"
	"robinhood-assignment/internal/core/domains"
	"robinhood-assignment/internal/core/ports"
	"robinhood-assignment/internal/dto"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type interviewService struct {
	interviewAppointmentRepo ports.InterviewAppointmentRepository
}

func NewInterviewService(interviewAppointmentRepo ports.InterviewAppointmentRepository) ports.InterviewService {
	return &interviewService{
		interviewAppointmentRepo: interviewAppointmentRepo,
	}
}

func (s *interviewService) GetInterviewAppointments(ctx context.Context, offset int64, limit int64) ([]domains.InterviewAppointment, error) {
	data, err := s.interviewAppointmentRepo.GetAll(ctx, offset, limit)
	if err != nil {
		return nil, helpers.NewCustomError(http.StatusInternalServerError, "Cannot get interview appointment.")
	}
	return data, nil
}
func (s *interviewService) GetInterviewAppointment(ctx context.Context, id string) (*domains.InterviewAppointment, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, helpers.NewCustomError(http.StatusInternalServerError, "Something went wrong please contact developer")
	}
	data, err := s.interviewAppointmentRepo.Get(ctx, objID)
	return data, nil
}

func (s *interviewService) CreateInterviewAppointment(ctx context.Context, req *dto.CreateInterviewAppointmentRequest) (*domains.InterviewAppointment, error) {
	userId, err := primitive.ObjectIDFromHex(req.CreatedBy)
	if err != nil {
		return nil, helpers.NewCustomError(http.StatusInternalServerError, "Something went wrong please contact developer")
	}
	params := &domains.CreateInterviewAppointmentParams{
		Title:       req.Title,
		Description: req.Description,
		UserID:      userId,
	}
	data, err := s.interviewAppointmentRepo.Create(ctx, params)
	return data, nil
}
