package repositories

import (
	"context"
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

func (r *interviewAppointmentRepository) GetAll(ctx context.Context, offset uint32, limit uint32) ([]domains.InterviewAppointment, error) {
	pipeline := []bson.D{
		{{Key: "$match", Value: bson.D{{Key: "isArchived", Value: false}}}},
		{{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "user"},
				{Key: "localField", Value: "createUserId"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "createUser"},
			},
		}},
		{{
			Key: "$unwind",
			Value: bson.D{
				{Key: "path", Value: "$createUser"},
				{Key: "preserveNullAndEmptyArrays", Value: false},
			},
		}},
		{{Key: "$skip", Value: offset}},
		{{Key: "$limit", Value: limit}},
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
		{{Key: "$match", Value: bson.D{{Key: "_id", Value: id}, {Key: "isArchived", Value: false}}}},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$comments"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		{{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "user"},
				{Key: "localField", Value: "comments.userId"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "comments.user"},
			},
		}},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$comments.userId"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$comments.user"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}},
		{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$_id"}, {Key: "comments", Value: bson.D{{Key: "$push", Value: "$comments"}}}}}},
		{{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "interviewAppointment"},
				{Key: "localField", Value: "_id"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "interviewDetail"},
			},
		}},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$interviewDetail"}}}},
		{{Key: "$addFields", Value: bson.D{{Key: "interviewDetail.comments", Value: "$comments"}}}},
		{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: "$interviewDetail"}}}},
		{{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "user"},
				{Key: "localField", Value: "createUserId"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "createUser"},
			},
		}},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$createUser"}, {Key: "preserveNullAndEmptyArrays", Value: false}}}},
		{{Key: "$limit", Value: 1}},
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
		Status:       "TODO",
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

func (r *interviewAppointmentRepository) Update(ctx context.Context, params *domains.UpdateInterviewAppointmentParams) (*domains.InterviewAppointment, error) {
	filter := bson.D{{Key: "_id", Value: params.ID}, {Key: "isArchived", Value: false}}
	updateValue := bson.D{{Key: "updatedAt", Value: time.Now()}}
	if params.Title != "" {
		updateValue = append(updateValue, bson.E{Key: "title", Value: params.Title})
	}
	if params.Description != "" {
		updateValue = append(updateValue, bson.E{Key: "description", Value: params.Description})
	}
	if params.Status != "" {
		updateValue = append(updateValue, bson.E{Key: "status", Value: params.Status})
	}
	update := bson.D{{Key: "$set", Value: updateValue}}
	opts := &options.FindOneAndUpdateOptions{}
	opts.SetReturnDocument(options.After).SetUpsert(false)
	res := domains.InterviewAppointment{}
	updated := r.col.FindOneAndUpdate(ctx, filter, update, opts)
	if err := updated.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	if err := updated.Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *interviewAppointmentRepository) ArchiveInterviewAppointment(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.D{{Key: "_id", Value: id}, {Key: "isArchived", Value: false}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "isArchived", Value: true}}}}
	opts := &options.FindOneAndUpdateOptions{}
	opts.SetUpsert(false)
	if err := r.col.FindOneAndUpdate(ctx, filter, update, opts).Err(); err != nil {
		return err
	}
	return nil
}

func (r *interviewAppointmentRepository) AddComment(ctx context.Context, params *domains.AddInterviewCommentParams) error {
	now := time.Now()
	filter := bson.D{{Key: "_id", Value: params.ID}, {Key: "isArchived", Value: false}}
	comment := domains.AddInterviewComment{
		ID:        primitive.NewObjectID(),
		Comment:   params.Comment,
		UserID:    params.UserID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "comments", Value: comment}}}}
	opts := &options.FindOneAndUpdateOptions{}
	opts.SetUpsert(false)
	if err := r.col.FindOneAndUpdate(ctx, filter, update, opts).Err(); err != nil {
		return err
	}
	return nil
}

func (r *interviewAppointmentRepository) UpdateComment(ctx context.Context, params *domains.UpdateInterviewCommentParams) error {
	now := time.Now()
	filter := bson.D{
		{Key: "_id", Value: params.ID},
		{Key: "isArchived", Value: false},
		{Key: "comments._id", Value: params.CommentID},
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "comments.$.comment", Value: params.Comment},
			{Key: "comments.$.updatedAt", Value: now},
		}},
	}
	opts := &options.FindOneAndUpdateOptions{}
	opts.SetUpsert(false)
	if err := r.col.FindOneAndUpdate(ctx, filter, update, opts).Err(); err != nil {
		return err
	}
	return nil
}
