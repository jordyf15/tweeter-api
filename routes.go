package main

import (
	"github.com/jordyf15/tweeter-api/controllers"
	fr "github.com/jordyf15/tweeter-api/follow/repository"
	fu "github.com/jordyf15/tweeter-api/follow/usecase"
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
	followRepo := fr.NewFollowRepo(db)

	tokenUsecase := tu.NewTokenUsecase(tokenRepo)
	userUsecase := uu.NewUserUsecase(userRepo, tokenRepo, _storage)
	followUsecase := fu.NewFollowUsecase(followRepo, userRepo)

	tokenController := controllers.NewTokenController(tokenUsecase)
	userController := controllers.NewUsersController(userUsecase)
	followController := controllers.NewFollowsController(followUsecase)

	router.POST("register", userController.Register)
	router.POST("login", userController.Login)

	router.POST("users/:user_id/password/change", middlewares.EnsureCurrentUserIDMatchesPath, userController.ChangeUserPassword)
	router.PATCH("users/:user_id", middlewares.EnsureCurrentUserIDMatchesPath, userController.EditUserProfile)
	router.POST("users/:user_id/follow", followController.FollowUser)

	router.POST("tokens/refresh", tokenController.RefreshAccessToken)
	router.DELETE("tokens/remove", tokenController.DeleteRefreshToken)
}
