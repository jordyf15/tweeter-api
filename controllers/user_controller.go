package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jordyf15/tweeter-api/models"
	"github.com/jordyf15/tweeter-api/user"
)

type UsersController interface {
	Register(c *gin.Context)
}

type usersController struct {
	userUsecase user.Usecase
}

func NewUsersController(userUsecase user.Usecase) UsersController {
	return &usersController{userUsecase: userUsecase}
}

func (controller *usersController) Register(c *gin.Context) {
	user := &models.User{}
	user.Fullname = c.PostForm("fullname")
	user.Username = c.PostForm("username")
	user.Email = c.PostForm("email")
	user.Password = c.PostForm("password")

	resp, err := controller.userUsecase.Create(user)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
