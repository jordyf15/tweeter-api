package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jordyf15/tweeter-api/controllers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	followMocks "github.com/jordyf15/tweeter-api/follow/mocks"
)

func TestFollowController(t *testing.T) {
	suite.Run(t, new(followControllerSuite))
}

type followControllerSuite struct {
	suite.Suite
	router     *gin.Engine
	response   *httptest.ResponseRecorder
	controller controllers.FollowsController
	context    *gin.Context
}

func (s *followControllerSuite) SetupTest() {
	followUsecase := new(followMocks.Usecase)

	followUsecase.On("FollowUser", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
	followUsecase.On("UnfollowUser", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

	s.controller = controllers.NewFollowsController(followUsecase)
	s.response = httptest.NewRecorder()
	s.context, s.router = gin.CreateTestContext(s.response)

	s.router.POST("/users/:user_id/follow", func(c *gin.Context) {
		c.Set("current_user_id", "userID")
		c.Next()
	}, s.controller.FollowUser)
	s.router.DELETE("/users/:user_id/follow", func(c *gin.Context) {
		c.Set("current_user_id", "userID")
		c.Next()
	}, s.controller.UnfollowUser)
}

func (s *followControllerSuite) TestFollowUserSuccessful() {
	s.context.Request, _ = http.NewRequest("POST", "/users/userID/follow", nil)
	s.router.ServeHTTP(s.response, s.context.Request)

	assert.Equal(s.T(), http.StatusNoContent, s.response.Code)
}

func (s *followControllerSuite) TestUnfollowUserSuccessful() {
	s.context.Request, _ = http.NewRequest("DELETE", "/users/userID/follow", nil)
	s.router.ServeHTTP(s.response, s.context.Request)

	assert.Equal(s.T(), http.StatusNoContent, s.response.Code)
}
