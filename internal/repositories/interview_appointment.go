package repositories

import (
	"context"
	"robinhood-assignment/internal/core/domains"
	"robinhood-assignment/internal/core/ports"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	pipeline := []bson.D{
		{
			{
				Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "user"},
					{Key: "localField", Value: "createUserId"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "createUser"},
				},
			},
		},
		{
			{
				Key: "$unwind",
				Value: bson.D{
					{Key: "path", Value: "$createUser"},
					{Key: "preserveNullAndEmptyArrays", Value: false},
				},
			},
		},
		{
			{Key: "$skip", Value: offset},
		},
		{
			{Key: "$limit", Value: limit},
		},
	}

	res := []domains.InterviewAppointment{}
	cur, err := r.col.Aggregate(ctx, pipeline)
	if err != nil {
		return res, err
	}
	if err := cur.All(ctx, &res); err != nil {
		return res, err
	}
	return res, nil
}

func (r *interviewAppointmentRepository) Get(ctx context.Context, id primitive.ObjectID) (*domains.InterviewAppointment, error) {
	pipeline := []bson.D{
		{
			{
				Key:   "$match",
				Value: bson.D{{Key: "_id", Value: id}},
			},
		},
		{
			{
				Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "user"},
					{Key: "localField", Value: "createUserId"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "createUser"},
				},
			},
		},
		{
			{
				Key: "$unwind",
				Value: bson.D{
					{Key: "path", Value: "$createUser"},
					{Key: "preserveNullAndEmptyArrays", Value: false},
				},
			},
		},
		{
			{Key: "$limit", Value: 1},
		},
	}
	res := []domains.InterviewAppointment{}
	cur, err := r.col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	if err := cur.All(ctx, &res); err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return &res[0], nil
}

func (r *interviewAppointmentRepository) Create(ctx context.Context, params *domains.CreateInterviewAppointmentParams) (*domains.CreateInterviewAppointment, error) {
	now := time.Now()
	interviewAppointment := domains.CreateInterviewAppointment{
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
