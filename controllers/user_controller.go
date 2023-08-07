package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jordyf15/tweeter-api/custom_errors"
	"github.com/jordyf15/tweeter-api/models"
	"github.com/jordyf15/tweeter-api/user"
	"github.com/jordyf15/tweeter-api/utils"
)

type UsersController interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	ChangeUserPassword(c *gin.Context)
	EditUserProfile(c *gin.Context)
}

type usersController struct {
	userUsecase user.Usecase
}

const (
	pictureSizesInMb = 1024 * 1024
)

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

func (controller *usersController) EditUserProfile(c *gin.Context) {
	userId := c.Param("user_id")

	updates := map[string]string{}
	if newFullname, isExist := c.GetPostForm("fullname"); isExist {
		updates["fullname"] = newFullname
	}

	if newUsername, isExist := c.GetPostForm("username"); isExist {
		updates["username"] = newUsername
	}

	if newDescription, isExist := c.GetPostForm("description"); isExist {
		updates["description"] = newDescription
	}

	if newEmail, isExist := c.GetPostForm("email"); isExist {
		updates["email"] = newEmail
	}

	willRemoveProfileImage := false
	if removeProfileImage := c.PostForm("is_remove_profile_image"); removeProfileImage == "true" {
		willRemoveProfileImage = true
	}

	willRemoveBackgroundImage := false
	if removeBackgroundImage := c.PostForm("is_remove_background_image"); removeBackgroundImage == "true" {
		willRemoveBackgroundImage = true
	}

	profileImageFileHeader, err := c.FormFile("profile_image")
	if err != nil {
		fmt.Println(err.Error())
	}

	backgroundImageFileHeader, err := c.FormFile("background_image")
	if err != nil {
		fmt.Println(err.Error())
	}

	var profileImageFile utils.NamedFileReader
	if profileImageFileHeader != nil {
		if profileImageFileHeader.Size > pictureSizesInMb*2 {
			respondBasedOnError(c, custom_errors.ErrProfileImageTooLarge)
			return
		}

		file, err := profileImageFileHeader.Open()
		if err != nil {
			fmt.Println(err.Error())
		}
		profileImageFile = utils.NewNamedFileReader(file, profileImageFileHeader.Filename)
	}

	var backgroundImageFile utils.NamedFileReader
	if backgroundImageFileHeader != nil {
		if backgroundImageFileHeader.Size > pictureSizesInMb*5 {
			respondBasedOnError(c, custom_errors.ErrBackgroundImageTooLarge)
			return
		}

		file, err := backgroundImageFileHeader.Open()
		if err != nil {
			fmt.Println(err.Error())
		}

		backgroundImageFile = utils.NewNamedFileReader(file, backgroundImageFileHeader.Filename)
	}

	user, err := controller.userUsecase.EditUserProfile(userId, updates, profileImageFile, backgroundImageFile, willRemoveProfileImage, willRemoveBackgroundImage)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}
