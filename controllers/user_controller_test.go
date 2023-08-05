package controllers_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jordyf15/tweeter-api/controllers"
	"github.com/jordyf15/tweeter-api/models"
	userMocks "github.com/jordyf15/tweeter-api/user/mocks"
	"github.com/jordyf15/tweeter-api/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestUserController(t *testing.T) {
	suite.Run(t, new(userControllerSuite))
}

type userControllerSuite struct {
	suite.Suite
	router     *gin.Engine
	response   *httptest.ResponseRecorder
	controller controllers.UsersController
	context    *gin.Context
}

var (
	profileImage1 = &models.Image{
		Width:  100,
		Height: 100,
		URL:    "url",
	}
	profileImage2 = &models.Image{
		Width:  200,
		Height: 200,
		URL:    "url",
	}
	backgroundImage = &models.Image{
		Width:  300,
		Height: 300,
		URL:    "url",
	}

	uctUser = &models.User{
		ID:              "userId",
		Fullname:        "fullname",
		Username:        "username",
		Email:           "email",
		Description:     "",
		FollowerCount:   0,
		FollowingCount:  0,
		BackgroundImage: *backgroundImage,
		ProfileImages: models.Images{
			profileImage1, profileImage2,
		},
	}
)

func (s *userControllerSuite) SetupTest() {
	userUsecase := new(userMocks.Usecase)

	response := utils.DataResponse(uctUser, map[string]interface{}{
		"access_token":  "accessToken",
		"refresh_token": "refreshToken",
		"expires_at":    1,
	})

	userUsecase.On("Create", mock.AnythingOfType("*models.User")).Return(response, nil)

	s.controller = controllers.NewUsersController(userUsecase)
	s.response = httptest.NewRecorder()
	s.context, s.router = gin.CreateTestContext(s.response)

	s.router.POST("/register", s.controller.Register)
}

func (s *userControllerSuite) TestCreateUser() {
	var receivedResponse map[string]interface{}

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	fullname, _ := writer.CreateFormField("fullname")
	fullname.Write([]byte(uctUser.Fullname))
	username, _ := writer.CreateFormField("username")
	username.Write([]byte(uctUser.Email))
	email, _ := writer.CreateFormField("email")
	email.Write([]byte(uctUser.Email))
	password, _ := writer.CreateFormField("password")
	password.Write([]byte("Password123!"))
	writer.Close()

	s.context.Request, _ = http.NewRequest("POST", "/register", buf)
	s.context.Request.Header.Set("Content-Type", writer.FormDataContentType())
	s.router.ServeHTTP(s.response, s.context.Request)
	json.NewDecoder(s.response.Body).Decode(&receivedResponse)

	assert.Equal(s.T(), http.StatusOK, s.response.Code)

	data, isExist := receivedResponse["data"].(map[string]interface{})
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), uctUser.ID, data["id"])
	assert.Equal(s.T(), uctUser.Username, data["username"])
	assert.Equal(s.T(), uctUser.Fullname, data["fullname"])
	assert.Equal(s.T(), uctUser.Email, data["email"])
	assert.Equal(s.T(), uctUser.Description, data["description"])
	assert.Equal(s.T(), float64(uctUser.FollowerCount), data["follower_count"])
	assert.Equal(s.T(), float64(uctUser.FollowingCount), data["following_count"])

	profileImages, isExist := data["profile_images"].([]interface{})
	assert.True(s.T(), isExist)
	profileImage1 := profileImages[0].(map[string]interface{})
	width, isExist := profileImage1["width"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), float64(uctUser.ProfileImages[0].Width), width)
	height, isExist := profileImage1["height"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), float64(uctUser.ProfileImages[0].Height), height)
	url, isExist := profileImage1["url"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), uctUser.ProfileImages[0].URL, url)

	profileImage2 := profileImages[1].(map[string]interface{})
	width, isExist = profileImage2["width"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), float64(uctUser.ProfileImages[1].Width), width)
	height, isExist = profileImage2["height"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), float64(uctUser.ProfileImages[1].Height), height)
	url, isExist = profileImage2["url"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), uctUser.ProfileImages[1].URL, url)

	backgroundImage, isExist := data["background_image"].(map[string]interface{})
	assert.True(s.T(), isExist)
	width, isExist = backgroundImage["width"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), float64(uctUser.BackgroundImage.Width), width)
	height, isExist = backgroundImage["height"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), float64(uctUser.BackgroundImage.Height), height)
	url, isExist = backgroundImage["url"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), uctUser.BackgroundImage.URL, url)

	meta, isExist := receivedResponse["meta"].(map[string]interface{})
	assert.True(s.T(), isExist)
	accessToken, isExist := meta["access_token"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), "accessToken", accessToken)
	refreshToken, isExist := meta["refresh_token"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), refreshToken, "refreshToken")
}
