package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jordyf15/tweeter-api/custom_errors"
	"github.com/jordyf15/tweeter-api/models"
	"github.com/jordyf15/tweeter-api/user"
)

type UsersController interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	ChangeUserPassword(c *gin.Context)
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

func (controller *usersController) Login(c *gin.Context) {
	errors := make([]error, 0)

	login := c.PostForm("login")
	password := c.PostForm("password")
	if login == "" {
		errors = append(errors, custom_errors.ErrEmptyLogin)
	}
	if password == "" {
		errors = append(errors, custom_errors.ErrEmptyPassword)
	}

	if len(errors) > 0 {
		respondBasedOnError(c, &custom_errors.MultipleErrors{Errors: errors})
		return
	}

	response, err := controller.userUsecase.Login(login, password)
	if err != nil {
		respondBasedOnError(c, err)
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func (controller *usersController) ChangeUserPassword(c *gin.Context) {
	errors := make([]error, 0)
	userID := c.Param("user_id")

	oldPassword := c.PostForm("old_password")
	newPassword := c.PostForm("new_password")
	if oldPassword == "" {
		errors = append(errors, custom_errors.ErrEmptyOldPassword)
	}
	if newPassword == "" {
		errors = append(errors, custom_errors.ErrEmptyNewPassword)
	}

	if len(errors) > 0 {
		respondBasedOnError(c, &custom_errors.MultipleErrors{Errors: errors})
		return
	}

	err := controller.userUsecase.ChangeUserPassword(userID, oldPassword, newPassword)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
