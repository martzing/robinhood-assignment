package validate

import (
	"net/http"
	"robinhood-assignment/helpers"
	"robinhood-assignment/internal/core/domains"
	"robinhood-assignment/internal/core/ports"
	"robinhood-assignment/internal/dto"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

type interviewValidate struct {
}

func NewInterviewValidate() ports.InterviewValidate {
	return &interviewValidate{}
}

func (v interviewValidate) ValidateGetInterviewAppointments(ctx *gin.Context) (*dto.GetInterviewAppointmentsRequest, error) {
	req := dto.GetInterviewAppointmentsRequest{}
	if page, ok := ctx.GetQuery("page"); ok {
		v, err := strconv.Atoi(page)
		if err != nil {
			return nil, helpers.NewCustomError(http.StatusBadRequest, "Invalid page query parameter")
		}
		i := uint64(v)
		req.Page = i
	}
	if limit, ok := ctx.GetQuery("limit"); ok {
		v, err := strconv.Atoi(limit)
		if err != nil {
			return nil, helpers.NewCustomError(http.StatusBadRequest, "Invalid limit query parameter")
		}
		i := uint64(v)
		req.Limit = i
	}
	if _, err := govalidator.ValidateStruct(req); err != nil {
		return nil, helpers.NewCustomError(http.StatusBadRequest, err.Error())
	}
	return &req, nil
}

func (v interviewValidate) ValidateGetInterviewAppointment(ctx *gin.Context) (string, error) {
	id := ctx.Param("id")
	if id == "" {
		return "", helpers.NewCustomError(http.StatusBadRequest, "id: Missing required field")
	}
	formats := strfmt.Default
	if err := validate.FormatOf("id", "param", "bsonobjectid", id, formats); err != nil {
		return "", helpers.NewCustomError(http.StatusBadRequest, err.Error())
	}
	return id, nil
}

func (v interviewValidate) ValidateCreateInterviewAppointment(ctx *gin.Context) (*dto.CreateInterviewAppointmentRequest, error) {
	req := dto.CreateInterviewAppointmentRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		return nil, helpers.NewCustomError(http.StatusBadRequest, "Invalid input parameter")
	}
	if _, err := govalidator.ValidateStruct(req); err != nil {
		return nil, helpers.NewCustomError(http.StatusBadRequest, err.Error())
	}
	return &req, nil
}

func (v interviewValidate) ValidateUpdateInterviewAppointment(ctx *gin.Context) (*dto.UpdateInterviewAppointmentRequest, error) {
	req := dto.UpdateInterviewAppointmentRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		return nil, helpers.NewCustomError(http.StatusBadRequest, "Invalid input parameter")
	}
	id := ctx.Param("id")
	if id == "" {
		return nil, helpers.NewCustomError(http.StatusBadRequest, "id: Missing required field")
	}
	req.ID = id
	if req.Title == "" && req.Description == "" && req.Status == "" {
		return nil, helpers.NewCustomError(http.StatusBadRequest, "at least one field required")
	}
	if _, err := govalidator.ValidateStruct(req); err != nil {
		return nil, helpers.NewCustomError(http.StatusBadRequest, err.Error())
	}
	formats := strfmt.Default
	if err := validate.FormatOf("id", "body", "bsonobjectid", req.ID, formats); err != nil {
		return nil, helpers.NewCustomError(http.StatusBadRequest, err.Error())
	}
	return &req, nil
}

func (v interviewValidate) ValidateArchiveInterviewAppointment(ctx *gin.Context) (string, error) {
	id := ctx.Param("id")
	if id == "" {
		return "", helpers.NewCustomError(http.StatusBadRequest, "id: Missing required field")
	}
	formats := strfmt.Default
	if err := validate.FormatOf("id", "param", "bsonobjectid", id, formats); err != nil {
		return "", helpers.NewCustomError(http.StatusBadRequest, err.Error())
	}
	return id, nil
}

func (v interviewValidate) ValidateAddInterviewComment(ctx *gin.Context) (*dto.AddInterviewCommentRequest, error) {
	req := dto.AddInterviewCommentRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		return nil, helpers.NewCustomError(http.StatusBadRequest, "Invalid input parameter")
	}
	id := ctx.Param("id")
	if id == "" {
		return nil, helpers.NewCustomError(http.StatusBadRequest, "id: Missing required field")
	}
	req.ID = id
	value, exists := ctx.Get("user")
	if !exists {
		return nil, helpers.InternalError
	}
	user := value.(domains.User)
	req.UserID = user.ID.Hex()

	if _, err := govalidator.ValidateStruct(req); err != nil {
		return nil, helpers.NewCustomError(http.StatusBadRequest, err.Error())
	}
	return &req, nil
}

func (v interviewValidate) ValidateUpdateInterviewComment(ctx *gin.Context) (*dto.UpdateInterviewCommentRequest, error) {
	req := dto.UpdateInterviewCommentRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		return nil, helpers.NewCustomError(http.StatusBadRequest, "Invalid input parameter")
	}
	id := ctx.Param("id")
	if id == "" {
		return nil, helpers.NewCustomError(http.StatusBadRequest, "id: Missing required field")
	}
	req.ID = id
	commentId := ctx.Param("commentId")
	if commentId == "" {
		return nil, helpers.NewCustomError(http.StatusBadRequest, "commentId: Missing required field")
	}
	req.CommentID = commentId

	if _, err := govalidator.ValidateStruct(req); err != nil {
		return nil, helpers.NewCustomError(http.StatusBadRequest, err.Error())
	}
	formats := strfmt.Default
	if err := validate.FormatOf("id", "param", "bsonobjectid", id, formats); err != nil {
		return nil, helpers.NewCustomError(http.StatusBadRequest, err.Error())
	}
	if err := validate.FormatOf("commentId", "param", "bsonobjectid", commentId, formats); err != nil {
		return nil, helpers.NewCustomError(http.StatusBadRequest, err.Error())
	}
	return &req, nil
}
