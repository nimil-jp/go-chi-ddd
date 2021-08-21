package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"go-chi-ddd/config"
	"go-chi-ddd/constant"
	"go-chi-ddd/infrastructure/email"
	"go-chi-ddd/infrastructure/log"
	"go-chi-ddd/infrastructure/persistence"
	"go-chi-ddd/interface/handler"
	"go-chi-ddd/interface/middleware"
	"go-chi-ddd/usecase"
)

func Execute() {
	config.SetMode(constant.DebugMode)
	logger := log.Logger()

	// dependencies injection
	// ----- infrastructure -----
	emailDriver := email.New()

	// persistence
	userPersistence := persistence.NewUser()

	// ----- use case -----
	userUseCase := usecase.NewUser(emailDriver, userPersistence)

	// ----- handler -----
	userHandler := handler.NewUser(userUseCase)

	// api

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.Cors)

	get(
		r, "/health", func(w http.ResponseWriter, r *http.Request) error {
			w.WriteHeader(http.StatusOK)
			return nil
		},
	)

	r.Route(
		"/user", func(r chi.Router) {
			post(r, "/", userHandler.Create)
			post(r, "/login", userHandler.Login)
			get(r, "/refresh-token", userHandler.RefreshToken)
			patch(r, "/reset-password-request", userHandler.ResetPasswordRequest)
			patch(r, "/reset-password", userHandler.ResetPassword)
		},
	)

	logger.Info("Succeeded in setting up routes.")

	// serve
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Env.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	logger.Info("Succeeded in listen and serve.")

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %+v", err)
	}

	logger.Info("Server exiting")
}
