package repositories

import (
	"context"
	"fmt"
	"robinhood-assignment/internal/core/domains"
	"robinhood-assignment/internal/core/ports"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type interviewAppointmentRepository struct {
	mc  *mongo.Client
	db  string
	cn  string
	col *mongo.Collection
}

func NewInterviewAppointmentRepository(mc *mongo.Client, db string) ports.InterviewAppointmentRepository {
	cn := "interviewAppointment"
	return &interviewAppointmentRepository{
		mc:  mc,
		db:  db,
		cn:  cn,
		col: mc.Database(db).Collection(cn),
	}
}

func (r *interviewAppointmentRepository) GetAll(ctx context.Context, offset int64, limit int64) ([]domains.InterviewAppointment, error) {
	fmt.Printf("offset %#v\n", offset)
	fmt.Printf("limit %#v\n", limit)
	res := []domains.InterviewAppointment{}
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetLimit(limit).SetSkip(offset)

	cur, err := r.col.Find(ctx, bson.D{}, opts)
	if err != nil {
		return res, err
	}
	if err := cur.All(ctx, &res); err != nil {
		return res, err
	}
	return res, nil
}

func (r *interviewAppointmentRepository) Get(ctx context.Context, id primitive.ObjectID) (*domains.InterviewAppointment, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	res := domains.InterviewAppointment{}
	if err := r.col.FindOne(ctx, filter).Decode(&res); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (r *interviewAppointmentRepository) Create(ctx context.Context, params *domains.CreateInterviewAppointmentParams) (*domains.InterviewAppointment, error) {
	now := time.Now()
	interviewAppointment := domains.InterviewAppointment{
		ID:           primitive.NewObjectID(),
		Title:        params.Title,
		Description:  params.Description,
		Status:       "to_do",
		Comments:     []domains.InterviewComment{},
		CreateUserId: params.UserID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if _, err := r.col.InsertOne(ctx, interviewAppointment); err != nil {
		return nil, err
	}
	return &interviewAppointment, nil
}
