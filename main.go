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

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariable()
	initializers.Database()
	initializers.SyncDb()
}

func main() {
	initializers.SeedAdmin()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
			"http://localhost:3001",
			"http://127.0.0.1:3001",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	limiter := service.NewRateLimiter(60, time.Minute)
	r.Use(limiter.LimitMiddleWare())
	loginLimiter := service.NewRateLimiter(3, time.Minute)

	r.POST("/register", handler.Register)
	r.POST("/login", loginLimiter.LimitMiddleWareWithMessage("Too many login attempts. Please wait a minute and try again."), handler.Login)
	r.GET("/validate", handler.RequireAuth, handler.Validate)
	r.GET("/events", handler.GetAllEvents)

	userGroup := r.Group("/")
	userGroup.Use(handler.RequireAuth)
	{
		userGroup.POST("/book", handler.BookTransaction)
		userGroup.GET("/my-tickets", handler.GetMyTickets)
		userGroup.GET("/my-tickets/:id", handler.GetTicket)
	}

	adminGroup := r.Group("/admin")
	adminGroup.Use(handler.RequireAuth, handler.RequireAdmin)
	{
		adminGroup.GET("/users", handler.ListUsers)
		adminGroup.POST("/users", handler.CreateUser)
		adminGroup.PUT("/users/:id", handler.UpdateUser)
		adminGroup.DELETE("/users/:id", handler.DeleteUser)
		adminGroup.POST("/events", handler.CreateEvent)
		adminGroup.PUT("/events/:id", handler.UpdateEvent)
		adminGroup.DELETE("/events/:id", handler.DeleteEvent)
	}

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
