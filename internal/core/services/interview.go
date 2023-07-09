package services

import (
	"context"
	"net/http"
	"robinhood-assignment/helpers"
	"robinhood-assignment/internal/core/domains"
	"robinhood-assignment/internal/core/ports"
	"robinhood-assignment/internal/dto"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
		return nil, helpers.InternalError
	}
	data, err := s.interviewAppointmentRepo.Get(ctx, objID)
	if err != nil {
		return nil, helpers.InternalError
	}
	if data == nil {
		return nil, helpers.NewCustomError(http.StatusNotFound, "Interview appointment not found.")
	}
	return data, nil
}

func (s *interviewService) CreateInterviewAppointment(ctx context.Context, req *dto.CreateInterviewAppointmentRequest) (*domains.CreateInterviewAppointment, error) {
	userId, err := primitive.ObjectIDFromHex(req.CreatedBy)
	if err != nil {
		return nil, helpers.InternalError
	}
	params := &domains.CreateInterviewAppointmentParams{
		Title:       req.Title,
		Description: req.Description,
		UserID:      userId,
	}
	data, err := s.interviewAppointmentRepo.Create(ctx, params)
	if err != nil {
		return nil, helpers.InternalError
	}
	return data, nil
}

func (s *interviewService) UpdateInterviewAppointment(ctx context.Context, req *dto.UpdateInterviewAppointmentRequest) error {
	id, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return helpers.InternalError
	}
	params := &domains.UpdateInterviewAppointmentParams{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	}
	data, err := s.interviewAppointmentRepo.Update(ctx, params)
	if data == nil {
		return helpers.NewCustomError(http.StatusNotFound, "Interview appointment not found.")
	}
	if err != nil {
		return helpers.InternalError
	}
	return nil
}

func (s *interviewService) AddInterviewComment(ctx context.Context, req *dto.AddInterviewCommentRequest) error {
	id, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return helpers.InternalError
	}
	userId, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return helpers.InternalError
	}
	params := &domains.AddInterviewCommentParams{
		ID:      id,
		Comment: req.Comment,
		UserID:  userId,
	}
	if err := s.interviewAppointmentRepo.AddComment(ctx, params); err != nil {
		if err == mongo.ErrNoDocuments {
			return helpers.NewCustomError(http.StatusNotFound, "Interview appointment not found.")
		}
		return helpers.InternalError
	}
	return nil
}
