package handlers

import (
	"net/http"
	"robinhood-assignment/helpers"
	"robinhood-assignment/internal/core/domains"
	"robinhood-assignment/internal/core/ports"
	"robinhood-assignment/internal/dto"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type interviewHandler struct {
	interviewService  ports.InterviewService
	interviewValidate ports.InterviewValidate
}

func NewInterviewHandler(interviewService ports.InterviewService, interviewValidate ports.InterviewValidate) ports.InterviewHandler {
	return &interviewHandler{
		interviewService:  interviewService,
		interviewValidate: interviewValidate,
	}
}

func (h *interviewHandler) GetInterviewAppointments(ctx *gin.Context) {
	req, err := h.interviewValidate.ValidateGetInterviewAppointments(ctx)
	if err != nil {
		errRes := helpers.ErrorHandler(err)
		ctx.AbortWithStatusJSON(errRes.StatusCode, errRes)
		return
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 20
	}
	offset := (req.Page - 1) * req.Limit
	limit := req.Limit + 1

	data, err := h.interviewService.GetInterviewAppointments(ctx, int64(offset), int64(limit))
	if err != nil {
		errRes := helpers.ErrorHandler(err)
		ctx.AbortWithStatusJSON(errRes.StatusCode, errRes)
		return
	}
	interviews := make([]dto.InterviewAppointment, len(data))
	for i := 0; i < len(data); i++ {
		interviews[i] = dto.InterviewAppointment{
			ID:          data[i].ID.Hex(),
			Title:       data[i].Title,
			Description: data[i].Description,
			Status:      data[i].Status,
			CreateUser: dto.User{
				Name:     data[i].CreateUser.Name,
				Email:    data[i].CreateUser.Email,
				ImageUrl: data[i].CreateUser.ImageUrl,
			},
			CreatedAt: data[i].CreatedAt,
		}
	}
	size, hasNext := helpers.Paginate(&interviews, int64(req.Limit))
	response := dto.GetInterviewAppointmentsResponse{
		StatusCode: http.StatusOK,
		Data:       interviews,
		Pagination: dto.Pagination{
			Page:    uint64(req.Page),
			Size:    uint64(size),
			HasNext: hasNext,
		},
	}
	ctx.JSON(http.StatusOK, response)
}
func (h *interviewHandler) GetInterviewAppointment(ctx *gin.Context) {
	id, err := h.interviewValidate.ValidateGetInterviewAppointment(ctx)
	if err != nil {
		errRes := helpers.ErrorHandler(err)
		ctx.AbortWithStatusJSON(errRes.StatusCode, errRes)
		return
	}
	data, err := h.interviewService.GetInterviewAppointment(ctx, id)
	if err != nil {
		errRes := helpers.ErrorHandler(err)
		ctx.AbortWithStatusJSON(errRes.StatusCode, errRes)
		return
	}

	comments := []dto.InterviewComment{}
	for i := 0; i < len(data.Comments); i++ {
		if !data.Comments[i].ID.IsZero() {
			comments = append(comments, dto.InterviewComment{
				ID:      data.Comments[i].ID.Hex(),
				Comment: data.Comments[i].Comment,
				User: dto.User{
					Name:     data.Comments[i].User.Name,
					Email:    data.Comments[i].User.Email,
					ImageUrl: data.Comments[i].User.ImageUrl,
				},
				CreatedAt: data.CreatedAt,
			})
		}
	}
	response := dto.GetInterviewAppointmentResponse{
		StatusCode: http.StatusOK,
		Data: dto.InterviewAppointmentDetail{
			ID:          data.ID.Hex(),
			Title:       data.Title,
			Description: data.Description,
			Status:      data.Status,
			CreateUser: dto.User{
				Name:     data.CreateUser.Name,
				Email:    data.CreateUser.Email,
				ImageUrl: data.CreateUser.ImageUrl,
			},
			CreatedAt: data.CreatedAt,
			Comments:  comments,
		},
	}

	ctx.JSON(http.StatusOK, response)
}
func (h *interviewHandler) CreateInterviewAppointment(ctx *gin.Context) {
	req, err := h.interviewValidate.ValidateCreateInterviewAppointment(ctx)
	if err != nil {
		errRes := helpers.ErrorHandler(err)
		ctx.AbortWithStatusJSON(errRes.StatusCode, errRes)
		return
	}
	data, err := h.interviewService.CreateInterviewAppointment(ctx, req)
	if err != nil {
		errRes := helpers.ErrorHandler(err)
		ctx.AbortWithStatusJSON(errRes.StatusCode, errRes)
		return
	}
	response := dto.CreateInterviewAppointmentResponse{
		StatusCode: http.StatusCreated,
		Data: dto.InterviewAppointmentDetail{
			ID:          data.ID.Hex(),
			Title:       data.Title,
			Description: data.Description,
			Status:      data.Status,
			CreateUser:  dto.User{},
			CreatedAt:   data.CreatedAt,
			Comments:    []dto.InterviewComment{},
		},
	}
	ctx.JSON(http.StatusCreated, response)
}

func (h *interviewHandler) UpdateInterviewAppointment(ctx *gin.Context) {
	req, err := h.interviewValidate.ValidateUpdateInterviewAppointment(ctx)
	if err != nil {
		errRes := helpers.ErrorHandler(err)
		ctx.AbortWithStatusJSON(errRes.StatusCode, errRes)
		return
	}
	if err := h.interviewService.UpdateInterviewAppointment(ctx, req); err != nil {
		errRes := helpers.ErrorHandler(err)
		ctx.AbortWithStatusJSON(errRes.StatusCode, errRes)
		return
	}
	response := dto.BaseResponse{
		StatusCode: http.StatusOK,
		Message:    "success",
	}
	ctx.JSON(http.StatusOK, response)
}

func (h *interviewHandler) ArchiveInterviewAppointment(ctx *gin.Context) {
	id, err := h.interviewValidate.ValidateGetInterviewAppointment(ctx)
	if err != nil {
		errRes := helpers.ErrorHandler(err)
		ctx.AbortWithStatusJSON(errRes.StatusCode, errRes)
		return
	}
	if err := h.interviewService.ArchiveInterviewAppointment(ctx, id); err != nil {
		errRes := helpers.ErrorHandler(err)
		ctx.AbortWithStatusJSON(errRes.StatusCode, errRes)
		return
	}
	response := dto.BaseResponse{
		StatusCode: http.StatusOK,
		Message:    "success",
	}
	ctx.JSON(http.StatusOK, response)
}

func (h *interviewHandler) AddInterviewComment(ctx *gin.Context) {
	id, _ := primitive.ObjectIDFromHex("64a9d15a033183d31aded893")
	ctx.Set("user", domains.User{
		ID:   id,
		Name: "samart",
	})
	req, err := h.interviewValidate.ValidateAddInterviewComment(ctx)
	if err != nil {
		errRes := helpers.ErrorHandler(err)
		ctx.AbortWithStatusJSON(errRes.StatusCode, errRes)
		return
	}
	if err := h.interviewService.AddInterviewComment(ctx, req); err != nil {
		errRes := helpers.ErrorHandler(err)
		ctx.AbortWithStatusJSON(errRes.StatusCode, errRes)
		return
	}

	ctx.JSON(http.StatusOK, dto.BaseResponse{
		StatusCode: http.StatusOK,
		Message:    "success",
	})
}

func (h *interviewHandler) UpdateInterviewComment(ctx *gin.Context) {
	req, err := h.interviewValidate.ValidateUpdateInterviewComment(ctx)
	if err != nil {
		errRes := helpers.ErrorHandler(err)
		ctx.AbortWithStatusJSON(errRes.StatusCode, errRes)
		return
	}
	if err := h.interviewService.UpdateInterviewComment(ctx, req); err != nil {
		errRes := helpers.ErrorHandler(err)
		ctx.AbortWithStatusJSON(errRes.StatusCode, errRes)
		return
	}
	response := dto.BaseResponse{
		StatusCode: http.StatusOK,
		Message:    "success",
	}
	ctx.JSON(http.StatusOK, response)
}
