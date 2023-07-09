package dto

import (
	"time"
)

type BaseResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

type Pagination struct {
	Page    uint64 `json:"page"`
	Size    uint64 `json:"size"`
	HasNext bool   `json:"hasNext"`
}

type GetInterviewAppointmentsRequest struct {
	Page  uint64 `query:"page" valid:"type(uint64),optional"`
	Limit uint64 `query:"limit" valid:"type(uint64),optional"`
}

type GetInterviewAppointmentsResponse struct {
	StatusCode int                    `json:"statusCode"`
	Data       []InterviewAppointment `json:"data"`
	Pagination Pagination             `json:"pagination"`
}

type InterviewAppointment struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreateUser  User      `json:"createUser"`
	CreatedAt   time.Time `json:"createdAt"`
}

type GetInterviewAppointmentResponse struct {
	StatusCode int                        `json:"statusCode"`
	Data       InterviewAppointmentDetail `json:"data"`
}

type CreateInterviewAppointmentRequest struct {
	Title       string `json:"title" from:"title" valid:"type(string)"`
	Description string `json:"description" from:"description" valid:"type(string)"`
	CreatedBy   string `json:"createdBy" from:"createdBy" valid:"type(string)"`
}

type CreateInterviewAppointmentResponse struct {
	StatusCode int                        `json:"statusCode"`
	Data       InterviewAppointmentDetail `json:"data"`
}

type AddInterviewCommentRequest struct {
	ID      string `json:"id" from:"id" valid:"type(string)"`
	Comment string `json:"comment" from:"comment" valid:"type(string)"`
	UserID  string `json:"userId" from:"userId" valid:"type(string)"`
}

type UpdateInterviewCommentRequest struct {
	ID        string `json:"id" from:"id" valid:"type(string)"`
	CommentID string `json:"commentId" from:"commentId" valid:"type(string)"`
	Comment   string `json:"comment" from:"comment" valid:"type(string)"`
	UserID    string `json:"userId" from:"userId" valid:"type(string)"`
}

type UpdateInterviewAppointmentRequest struct {
	ID          string `json:"id" from:"id" valid:"type(string)"`
	Title       string `json:"title" from:"title" valid:"type(string),optional"`
	Description string `json:"description" from:"description" valid:"type(string),optional"`
	Status      string `json:"status" from:"status" valid:"type(string),in(TODO|IN_PROGRESS|DONE),optional"`
}

type InterviewAppointmentDetail struct {
	ID          string             `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Status      string             `json:"status"`
	CreateUser  User               `json:"createUser"`
	CreatedAt   time.Time          `json:"createdAt"`
	Comments    []InterviewComment `json:"comments"`
}

type InterviewComment struct {
	ID        string    `json:"id"`
	Comment   string    `json:"comment"`
	User      User      `json:"user"`
	CreatedAt time.Time `json:"CreatedAt"`
}

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	ImageUrl string `json:"imageUrl"`
}

type ErrorResponse struct {
	StatusCode int    `json:"statusCode" from:"statusCode"`
	Error      string `json:"error" from:"error"`
}
