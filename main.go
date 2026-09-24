package main

import (
	"booking-app/handler"
	"booking-app/initializers"
	"booking-app/service"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariable()
	initializers.Database()
	initializers.SyncDb()
}

func main() {

	r := gin.Default()

	limiter := service.NewRateLimiter(60, time.Minute)
	r.Use(limiter.LimitMiddleWare())
	loginLimiter := service.NewRateLimiter(3, time.Minute)

	r.POST("/user", handler.Register)
	r.POST("/login", loginLimiter.LimitMiddleWareWithMessage("Too many login attempts. Please wait a minute and try again."), handler.Login)
	r.POST("/event", handler.CreateEvent)
	r.POST("/book", handler.RequireAuth, handler.BookTransaction)
	r.GET("/event/:id", handler.GetEvent)
	r.GET("/validate", handler.RequireAuth, handler.Validate)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Println("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen error: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Forced shutdown: %s", err)
	}

	log.Println("Server exited gracefully")
}
