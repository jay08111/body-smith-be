package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"body-smith-be/internal/config"
	"body-smith-be/internal/db"
	"body-smith-be/internal/handler"
	"body-smith-be/internal/middleware"
	"body-smith-be/internal/repository"
	"body-smith-be/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if err := db.EnsureDatabase(cfg); err != nil {
		log.Fatalf("ensure database: %v", err)
	}

	database, err := db.New(cfg)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer database.Close()

	if err := db.RunMigrations(database, cfg.DBName); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	userRepo := repository.NewUserRepository(database)
	postRepo := repository.NewPostRepository(database)

	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiration)
	postService := service.NewPostService(postRepo)

	authHandler := handler.NewAuthHandler(authService)
	postHandler := handler.NewPostHandler(postService)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestLogger())

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.JWTAuth(cfg.JWTSecret), authHandler.Me)
		}

		publicPosts := api.Group("/posts")
		{
			publicPosts.GET("", postHandler.ListPublicPosts)
			publicPosts.GET("/:slug", postHandler.GetPublicPost)
		}

		admin := api.Group("/admin")
		admin.Use(middleware.JWTAuth(cfg.JWTSecret))
		{
			adminPosts := admin.Group("/posts")
			{
				adminPosts.POST("", postHandler.CreatePost)
				adminPosts.GET("", postHandler.ListAdminPosts)
				adminPosts.PUT("/:id", postHandler.UpdatePost)
				adminPosts.DELETE("/:id", postHandler.DeletePost)
			}
		}
	}

	server := &http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("server listening on :%s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown server: %v", err)
	}
}
