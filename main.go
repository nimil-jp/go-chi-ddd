package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-chi-ddd/infrastructure/email"
	"go-chi-ddd/infrastructure/persistence"
	"go-chi-ddd/interface/handler"
	"go-chi-ddd/usecase"

	"go-chi-ddd/infrastructure/log"
	// "go-chi-ddd/interface/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	logger := log.Logger()

	// err := jwt.SetUp(
	//	jwt.Option{
	//		Realm:            constant.DefaultRealm,
	//		SigningAlgorithm: jwt.HS256,
	//		SecretKey:        []byte(config.Env.App.Secret),
	//	},
	// )
	// if err != nil {
	//	panic(err)
	// }
	// logger.Info("Succeeded in setting up JWT.")
	//
	//
	// engine.Use(middleware.Log(logger, time.RFC3339, false))
	// engine.Use(middleware.RecoveryWithLog(logger, true))
	//
	// engine.GET("health", func(c *gin.Context) { c.Status(http.StatusOK) })

	// cors

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors)

	// dependencies injection
	// ----- infrastructure -----
	emailDriver := email.New()

	// persistence
	userPersistence := persistence.NewUser()

	// ----- use case -----
	userUseCase := usecase.NewUser(emailDriver, userPersistence)

	// ----- handler -----
	userHandler := handler.NewUser(userUseCase)

	handler.Get(
		r, "/health", func(w http.ResponseWriter, r *http.Request) error {
			w.WriteHeader(http.StatusOK)
			return nil
		},
	)
	r.Route(
		"/user", func(r chi.Router) {
			handler.Post(r, "/", userHandler.Create)
			handler.Post(r, "/login", userHandler.Login)
			handler.Get(r, "/refresh-token", userHandler.RefreshToken)
			handler.Patch(r, "/reset-password-request", userHandler.ResetPasswordRequest)
			handler.Patch(r, "/reset-password", userHandler.ResetPassword)
		},
	)

	logger.Info("Succeeded in setting up routes.")

	// serve
	var port = ":8080"
	if portEnv := os.Getenv("PORT"); portEnv != "" {
		port = portEnv
	}

	srv := &http.Server{
		Addr:    port,
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

func cors(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	}

	return http.HandlerFunc(fn)
}
