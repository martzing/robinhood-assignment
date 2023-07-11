package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"robinhood-assignment/config"
	"robinhood-assignment/helpers"
	"robinhood-assignment/infrastructures"
	"robinhood-assignment/internal/core/services"
	"robinhood-assignment/internal/handlers"
	"robinhood-assignment/internal/middlewares"
	"robinhood-assignment/internal/repositories"
	"robinhood-assignment/internal/validate"
	"syscall"
	"time"

	"github.com/asaskevich/govalidator"
	helmet "github.com/danielkov/gin-helmet"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	config.New()
	govalidator.SetFieldsRequiredByDefault(true)
}

func main() {

	mc := infrastructures.NewMongoDB()

	// helper layer
	myBcrypt := helpers.NewMyBcrypt()
	myJWT := helpers.NewMyJWT()

	interviewRepo := repositories.NewInterviewAppointmentRepository(mc, config.Get().Mongo.Database)
	userRepo := repositories.NewUserRepository(mc, config.Get().Mongo.Database)

	interviewService := services.NewInterviewService(interviewRepo, userRepo)
	authService := services.NewAuthService(userRepo, myBcrypt, myJWT)

	interviewValidate := validate.NewInterviewValidate()
	authValidate := validate.NewAuthValidate()

	interviewHandler := handlers.NewInterviewHandler(interviewService, interviewValidate)
	authHandler := handlers.NewAuthHandler(authService, authValidate)

	middleware := middlewares.NewMidlewares(myJWT)

	r := gin.Default()
	conf := cors.DefaultConfig()
	conf.AllowAllOrigins = true
	conf.AddAllowHeaders("Authorization")
	r.Use(cors.New(conf))
	r.Use(helmet.Default())
	r.GET("/healthz", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"message": "OK"}) })

	interviewGroup := r.Group("/api/interviews")
	interviewGroup.GET("", middleware.StaffMiddleware, interviewHandler.GetInterviewAppointments)
	interviewGroup.GET("/:id", middleware.StaffMiddleware, interviewHandler.GetInterviewAppointment)
	interviewGroup.POST("", middleware.StaffMiddleware, interviewHandler.CreateInterviewAppointment)
	interviewGroup.PATCH("/:id", middleware.StaffMiddleware, interviewHandler.UpdateInterviewAppointment)
	interviewGroup.PATCH("/:id/archive", middleware.StaffMiddleware, interviewHandler.ArchiveInterviewAppointment)
	interviewGroup.POST("/:id/comment", middleware.StaffMiddleware, interviewHandler.AddInterviewComment)
	interviewGroup.PATCH("/:id/comment/:commentId", middleware.StaffMiddleware, interviewHandler.UpdateInterviewComment)

	authGroup := r.Group("/api/auth")
	authGroup.POST("/login", authHandler.Login)
	authGroup.POST("/staff", middleware.AdminMiddleware, authHandler.CreateStaff)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Get().HTTPServer.Port),
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	log.Println("Server exiting")
}
