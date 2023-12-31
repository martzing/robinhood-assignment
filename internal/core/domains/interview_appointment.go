package domains

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateInterviewAppointment struct {
	ID           primitive.ObjectID `bson:"_id"`
	Title        string             `bson:"title"`
	Description  string             `bson:"description"`
	Comments     []InterviewComment `bson:"comments"`
	Status       string             `bson:"status"`
	IsArchived   bool               `bson:"isArchived"`
	CreateUserId primitive.ObjectID `bson:"createUserId"`
	CreatedAt    time.Time          `bson:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt"`
}

type InterviewAppointment struct {
	ID          primitive.ObjectID `bson:"_id"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	Comments    []InterviewComment `bson:"comments"`
	Status      string             `bson:"status"`
	IsArchived  bool               `bson:"isArchived"`
	CreateUser  User               `bson:"createUser"`
	CreatedAt   time.Time          `bson:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt"`
}

type CreateInterviewAppointmentParams struct {
	Title       string
	Description string
	UserID      primitive.ObjectID
}

type UpdateInterviewAppointmentParams struct {
	ID          primitive.ObjectID
	Title       string
	Description string
	Status      string
}
