package main

import (
	"github.com/gin-gonic/gin"
	"go-chi-ddd/infrastructure/email"
	"go-chi-ddd/infrastructure/persistence"
	"go-chi-ddd/interface/handler"
	"go-chi-ddd/usecase"
)

func inject(engine *gin.Engine) {
	// dependencies injection
	// ----- infrastructure -----
	emailDriver := email.New()

	// persistence
	userPersistence := persistence.NewUser()

	// ----- use case -----
	userUseCase := usecase.NewUser(emailDriver, userPersistence)

	// ----- handler -----
	handler.NewUser(engine.Group("user"), userUseCase)
}
