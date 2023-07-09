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

	iarepo := repositories.NewInterviewAppointmentRepository(mc, config.Get().Mongo.Database)
	userRepo := repositories.NewUserRepository(mc, config.Get().Mongo.Database)

	insvc := services.NewInterviewService(iarepo)
	authsvc := services.NewAuthService(userRepo, myBcrypt, myJWT)

	invalidate := validate.NewInterviewValidate()
	authvalidate := validate.NewAuthValidate()

	inhdl := handlers.NewInterviewHandler(insvc, invalidate)
	authhdl := handlers.NewAuthHandler(authsvc, authvalidate)

	r := gin.Default()
	conf := cors.DefaultConfig()
	conf.AllowAllOrigins = true
	conf.AddAllowHeaders("Authorization")
	r.Use(cors.New(conf))
	r.Use(helmet.Default())
	r.GET("/healthz", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"message": "OK"}) })

	r.GET("/api/interviews", inhdl.GetInterviewAppointments)
	r.GET("/api/interviews/:id", inhdl.GetInterviewAppointment)
	r.POST("/api/interviews", inhdl.CreateInterviewAppointment)
	r.PATCH("/api/interviews", inhdl.UpdateInterviewAppointment)
	r.POST("/api/interviews/comment", inhdl.AddInterviewComment)
	r.POST("/api/user", authhdl.RegisterAdmin)

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
