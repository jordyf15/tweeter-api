package main

import (
	"github.com/jordyf15/tweeter-api/controllers"
	tr "github.com/jordyf15/tweeter-api/token/repository"
	tu "github.com/jordyf15/tweeter-api/token/usecase"
)

func initializeRoutes() {
	tokenRepo := tr.NewTokenRepository(db, redisClient)

	tokenUsecase := tu.NewTokenUsecase(tokenRepo)

	tokenController := controllers.NewTokenController(tokenUsecase)

	router.POST("tokens/refresh", tokenController.RefreshAccessToken)
	router.DELETE("tokens/remove", tokenController.DeleteRefreshToken)
}
