package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jordyf15/tweeter-api/custom_errors"
	"github.com/jordyf15/tweeter-api/group"
	"github.com/jordyf15/tweeter-api/models"
	"github.com/jordyf15/tweeter-api/utils"
)

type GroupsController interface {
	CreateGroup(c *gin.Context)
}

type groupsController struct {
	usecase group.Usecase
}

func NewGroupsController(usecase group.Usecase) GroupsController {
	return &groupsController{usecase: usecase}
}

func (controller *groupsController) CreateGroup(c *gin.Context) {
	creatorID := c.MustGet("current_user_id").(string)

	_group := &models.Group{}
	_group.CreatorID = creatorID
	_group.Name = c.PostForm("name")
	_group.Description = c.PostForm("description")
	_group.IsOpen = c.PostForm("is_open") == "true"

	groupImageHeader, err := c.FormFile("image")
	if err != nil {
		fmt.Println(err.Error())
	}

	var groupImageFile utils.NamedFileReader
	if groupImageHeader == nil {
		respondBasedOnError(c, custom_errors.ErrGroupImageMissing)
		return
	} else {
		if groupImageHeader.Size > pictureSizesInMb*5 {
			respondBasedOnError(c, custom_errors.ErrGroupImageTooLarge)
			return
		}

		file, err := groupImageHeader.Open()
		if err != nil {
			fmt.Println(err.Error())
		}
		groupImageFile = utils.NewNamedFileReader(file, groupImageHeader.Filename)
	}

	createdGroup, err := controller.usecase.Create(_group, groupImageFile)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	c.JSON(http.StatusOK, createdGroup)
}
