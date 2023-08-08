package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jordyf15/tweeter-api/follow"
)

type FollowsController interface {
	FollowUser(c *gin.Context)
}

type followsController struct {
	usecase follow.Usecase
}

func NewFollowsController(usecase follow.Usecase) FollowsController {
	return &followsController{usecase: usecase}
}

func (controller *followsController) FollowUser(c *gin.Context) {
	followerID := c.MustGet("current_user_id").(string)
	followingID := c.Param("user_id")

	err := controller.usecase.FollowUser(followerID, followingID)
	if err != nil {
		respondBasedOnError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
