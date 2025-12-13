package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/crutchm/elite/internal/auth"
	"github.com/crutchm/elite/internal/config"
	"github.com/crutchm/elite/internal/database"
	"github.com/crutchm/elite/internal/handler"
	"github.com/crutchm/elite/internal/middleware"
	"github.com/crutchm/elite/internal/repository"
	"github.com/crutchm/elite/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx := context.Background()

	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	authService := auth.NewTelegramAuth(cfg.TelegramBotToken, cfg.JWTSecret)

	userRepo := repository.NewUserRepository(pool)
	voteRepo := repository.NewVoteRepository(pool)

	userService := service.NewUserService(userRepo)
	voteService := service.NewVoteService(voteRepo)

	authHandler := handler.NewAuthHandler(authService, userService)
	voteHandler := handler.NewVoteHandler(voteService)

	authMiddleware := middleware.AuthMiddleware(authService)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/api/auth", authHandler.Authenticate)
	mux.HandleFunc("/api/vote", authMiddleware(http.HandlerFunc(voteHandler.Vote)).ServeHTTP)

	// Обертываем все маршруты в CORS middleware
	handler := middleware.CORSMiddleware(mux)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
