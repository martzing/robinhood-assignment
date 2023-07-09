package middlewares

import (
	"net/http"
	"robinhood-assignment/config"
	"robinhood-assignment/internal/core/constants"
	"robinhood-assignment/internal/core/domains"
	"robinhood-assignment/internal/core/ports"
	"robinhood-assignment/internal/dto"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type middlewares struct {
	myJWT ports.MyJWT
}

func NewMidlewares(myJWT ports.MyJWT) ports.Middlewares {
	return &middlewares{myJWT}
}

func (m middlewares) AdminMiddleware(ctx *gin.Context) {
	authorization := ctx.GetHeader("Authorization")
	if authorization == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
			StatusCode: http.StatusUnauthorized,
			Error:      "Authorization is missing",
		})
		return
	}
	jwtToken := strings.Split(authorization, "Bearer ")
	if len(jwtToken) != 2 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
			StatusCode: http.StatusUnauthorized,
			Error:      "Invalid token format",
		})
		return
	}
	tokenString := jwtToken[1]
	claims := &domains.Claims{}
	if _, err := m.myJWT.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Get().Auth.JwtSecret), nil
	}); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
			StatusCode: http.StatusUnauthorized,
			Error:      err.Error(),
		})
		return
	}
	if claims.Role != constants.ADMIN_ROLE {
		ctx.AbortWithStatusJSON(http.StatusForbidden, dto.ErrorResponse{
			StatusCode: http.StatusForbidden,
			Error:      "You don't have permission for this API",
		})
		return
	}
	ctx.Set("userId", claims.UserID)
	ctx.Set("role", claims.Role)
	ctx.Next()
}

func (m middlewares) StaffMiddleware(ctx *gin.Context) {
	authorization := ctx.GetHeader("Authorization")
	if authorization == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
			StatusCode: http.StatusUnauthorized,
			Error:      "Authorization is missing",
		})
		return
	}
	jwtToken := strings.Split(authorization, "Bearer ")
	if len(jwtToken) != 2 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
			StatusCode: http.StatusUnauthorized,
			Error:      "Invalid token format",
		})
		return
	}
	tokenString := jwtToken[1]
	claims := &domains.Claims{}
	if _, err := m.myJWT.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Get().Auth.JwtSecret), nil
	}); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
			StatusCode: http.StatusUnauthorized,
			Error:      err.Error(),
		})
		return
	}
	if claims.Role != constants.STAFF_ROLE && claims.Role != constants.ADMIN_ROLE {
		ctx.AbortWithStatusJSON(http.StatusForbidden, dto.ErrorResponse{
			StatusCode: http.StatusForbidden,
			Error:      "You don't have permission for this API",
		})
		return
	}
	ctx.Set("userId", claims.UserID)
	ctx.Set("role", claims.Role)
	ctx.Next()
}
