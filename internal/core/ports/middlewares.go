package ports

import (
	"github.com/gin-gonic/gin"
)

type Middlewares interface {
	AdminMiddleware(ctx *gin.Context)
	StaffMiddleware(ctx *gin.Context)
}
