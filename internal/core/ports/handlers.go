package ports

import (
	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	RegisterAdmin(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type InterviewHandler interface {
	GetInterviewAppointments(ctx *gin.Context)
	GetInterviewAppointment(ctx *gin.Context)
	CreateInterviewAppointment(ctx *gin.Context)
	UpdateInterviewAppointment(ctx *gin.Context)
	ArchiveInterviewAppointment(ctx *gin.Context)
	AddInterviewComment(ctx *gin.Context)
	UpdateInterviewComment(ctx *gin.Context)
}
