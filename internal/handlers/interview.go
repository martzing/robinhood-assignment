package handlers

import (
	"net/http"
	"robinhood-assignment/helpers"
	"robinhood-assignment/internal/core/ports"
	"robinhood-assignment/internal/dto"

	"github.com/gin-gonic/gin"
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
			CreateUser:  dto.User{},
			CreatedAt:   data[i].CreatedAt,
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

	comments := make([]dto.InterviewComment, len(data.Comments))
	for i := 0; i < len(data.Comments); i++ {
		comments[i] = dto.InterviewComment{
			ID:        data.Comments[i].ID.Hex(),
			Comment:   data.Comments[i].Comment,
			User:      dto.User{},
			CreatedAt: data.CreatedAt,
		}
	}
	response := dto.GetInterviewAppointmentResponse{
		StatusCode: http.StatusOK,
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
