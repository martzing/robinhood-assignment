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
	userRepo                 ports.UserRepository
}

func NewInterviewService(interviewAppointmentRepo ports.InterviewAppointmentRepository, userRepo ports.UserRepository) ports.InterviewService {
	return &interviewService{
		interviewAppointmentRepo: interviewAppointmentRepo,
		userRepo:                 userRepo,
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

func (s *interviewService) CreateInterviewAppointment(ctx context.Context, req *dto.CreateInterviewAppointmentRequest) (*domains.InterviewAppointment, error) {
	userId, err := primitive.ObjectIDFromHex(req.CreatedBy)
	if err != nil {
		return nil, helpers.InternalError
	}
	user, err := s.userRepo.Get(ctx, userId)
	if err != nil {
		return nil, helpers.InternalError
	}
	if user == nil {
		return nil, helpers.NewCustomError(http.StatusUnauthorized, "Invalid user token")
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
	return &domains.InterviewAppointment{
		ID:          data.ID,
		Title:       data.Title,
		Description: data.Description,
		Comments:    data.Comments,
		Status:      data.Status,
		IsArchived:  data.IsArchived,
		CreateUser: domains.User{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			ImageUrl: user.ImageUrl,
		},
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}, nil
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

func (s *interviewService) ArchiveInterviewAppointment(ctx context.Context, id string) error {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return helpers.InternalError
	}
	if err := s.interviewAppointmentRepo.ArchiveInterviewAppointment(ctx, objId); err != nil {
		if err == mongo.ErrNoDocuments {
			return helpers.NewCustomError(http.StatusNotFound, "Interview appointment not found or archived")
		}
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

func (s *interviewService) UpdateInterviewComment(ctx context.Context, req *dto.UpdateInterviewCommentRequest) error {
	id, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return helpers.InternalError
	}
	commentId, err := primitive.ObjectIDFromHex(req.CommentID)
	if err != nil {
		return helpers.InternalError
	}
	data, err := s.interviewAppointmentRepo.Get(ctx, id)
	if err != nil {
		return helpers.InternalError
	}
	if data == nil {
		return helpers.NewCustomError(http.StatusNotFound, "Interview appointment not found.")
	}
	comment := &domains.InterviewComment{}
	for i := 0; i < len(data.Comments); i++ {
		if data.Comments[i].ID == commentId {
			comment = &data.Comments[i]
		}
	}
	if comment == nil {
		return helpers.NewCustomError(http.StatusNotFound, "Interview comment not found.")
	}
	if comment.User.ID.Hex() != req.UserID {
		return helpers.NewCustomError(http.StatusForbidden, "You don't have permission to update this comment")
	}
	params := domains.UpdateInterviewCommentParams{
		ID:        id,
		CommentID: commentId,
		Comment:   req.Comment,
	}
	if err := s.interviewAppointmentRepo.UpdateComment(ctx, &params); err != nil {
		if err == mongo.ErrNoDocuments {
			return helpers.NewCustomError(http.StatusNotFound, "Interview comment not found.")
		}
		return helpers.InternalError
	}
	return nil
}
