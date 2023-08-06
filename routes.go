package main

import (
	"github.com/jordyf15/tweeter-api/controllers"
	"github.com/jordyf15/tweeter-api/middlewares"
	"github.com/jordyf15/tweeter-api/storage"
	tr "github.com/jordyf15/tweeter-api/token/repository"
	tu "github.com/jordyf15/tweeter-api/token/usecase"
	ur "github.com/jordyf15/tweeter-api/user/repository"
	uu "github.com/jordyf15/tweeter-api/user/usecase"
)

func initializeRoutes() {
	_storage := storage.NewCloudStorage()

	tokenRepo := tr.NewTokenRepository(db, redisClient)
	userRepo := ur.NewUserRepository(db)

	tokenUsecase := tu.NewTokenUsecase(tokenRepo)
	userUsecase := uu.NewUserUsecase(userRepo, tokenRepo, _storage)

	tokenController := controllers.NewTokenController(tokenUsecase)
	userController := controllers.NewUsersController(userUsecase)

	router.POST("register", userController.Register)
	router.POST("login", userController.Login)

	router.POST("users/:user_id/password/change", middlewares.EnsureCurrentUserIDMatchesPath, userController.ChangeUserPassword)

	router.POST("tokens/refresh", tokenController.RefreshAccessToken)
	router.DELETE("tokens/remove", tokenController.DeleteRefreshToken)
}
