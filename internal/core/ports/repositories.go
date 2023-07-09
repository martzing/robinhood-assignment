package ports

import (
	"context"
	"robinhood-assignment/internal/core/domains"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	Get(ctx context.Context, id string) (*domains.User, error)
	GetByUsername(ctx context.Context, username string) (*domains.User, error)
	Create(ctx context.Context, params *domains.CreateUserParams) (*domains.User, error)
}

type InterviewAppointmentRepository interface {
	GetAll(ctx context.Context, offset int64, limit int64) ([]domains.InterviewAppointment, error)
	Get(ctx context.Context, id primitive.ObjectID) (*domains.InterviewAppointment, error)
	Create(ctx context.Context, params *domains.CreateInterviewAppointmentParams) (*domains.CreateInterviewAppointment, error)
	// Update(ctx context.Context, params *domains.UpdateFortuneTellingParams) (*domains.FortuneTelling, error)
	// Delete(ctx context.Context, id string) error
}
