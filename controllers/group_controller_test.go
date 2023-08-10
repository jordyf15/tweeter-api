package controllers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jordyf15/tweeter-api/controllers"
	"github.com/jordyf15/tweeter-api/custom_errors"
	"github.com/jordyf15/tweeter-api/group"
	groupMocks "github.com/jordyf15/tweeter-api/group/mocks"
	"github.com/jordyf15/tweeter-api/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestGroupController(t *testing.T) {
	suite.Run(t, new(groupControllerSuite))
}

type groupControllerSuite struct {
	suite.Suite
	router     *gin.Engine
	response   *httptest.ResponseRecorder
	controller controllers.GroupsController
	context    *gin.Context
}

var (
	thumbnailImg = &models.Image{
		Width:  group.ThumbnailPictureRes,
		Height: group.ThumbnailPictureRes,
		URL:    "url",
	}
	bannerImg = &models.Image{
		Width:  group.BannerPictureWidth,
		Height: group.BannerPictureHeight,
		URL:    "url",
	}

	ugtGroup = &models.Group{
		ID:          "groupID",
		Name:        "name",
		Description: "description",
		Images:      []*models.Image{thumbnailImg, bannerImg},
		MemberCount: 1,
		Creator: &models.User{
			Username: "username",
		},
		IsOpen:    true,
		CreatedAt: time.Now(),
	}
)

func (s *groupControllerSuite) SetupTest() {
	groupUsecase := new(groupMocks.Usecase)

	s.controller = controllers.NewGroupsController(groupUsecase)
	s.response = httptest.NewRecorder()
	s.context, s.router = gin.CreateTestContext(s.response)

	groupUsecase.On("Create", mock.AnythingOfType("*models.Group"), mock.Anything).Return(ugtGroup, nil)

	s.router.POST("/groups", func(c *gin.Context) {
		c.Set("current_user_id", "userID")
		c.Next()
	}, s.controller.CreateGroup)
}

func (s *groupControllerSuite) TestCreateGroupMissingImage() {
	var receivedResponse map[string]interface{}

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	name, _ := writer.CreateFormField("name")
	name.Write([]byte(ugtGroup.Name))
	description, _ := writer.CreateFormField("description")
	description.Write([]byte(ugtGroup.Description))
	isOpen, _ := writer.CreateFormField("is_open")
	isOpen.Write([]byte("true"))
	writer.Close()

	s.context.Request, _ = http.NewRequest("POST", "/groups", buf)
	s.context.Request.Header.Set("Content-Type", writer.FormDataContentType())
	s.router.ServeHTTP(s.response, s.context.Request)

	assert.Equal(s.T(), http.StatusBadRequest, s.response.Code)

	json.NewDecoder(s.response.Body).Decode(&receivedResponse)

	errors, isExist := receivedResponse["errors"].([]interface{})
	assert.True(s.T(), isExist)

	error1 := errors[0].(map[string]interface{})
	msg, isExist := error1["message"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), custom_errors.ErrGroupImageMissing.Message, msg)

	code, isExist := error1["code"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), float64(custom_errors.ErrGroupImageMissing.Code), code)
}

func (s *groupControllerSuite) TestCreateGroupSuccessful() {
	var receivedResponse map[string]interface{}

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	name, _ := writer.CreateFormField("name")
	name.Write([]byte(ugtGroup.Name))
	description, _ := writer.CreateFormField("description")
	description.Write([]byte(ugtGroup.Description))
	isOpen, _ := writer.CreateFormField("is_open")
	isOpen.Write([]byte("true"))

	imgHeader, _ := writer.CreateFormFile("image", "../assets/images/default-banner.jpg")
	file, _ := os.Open("../assets/images/default-banner.jpg")
	defer file.Close()

	_, _ = io.Copy(imgHeader, file)
	writer.Close()

	s.context.Request, _ = http.NewRequest("POST", "/groups", buf)
	s.context.Request.Header.Set("Content-Type", writer.FormDataContentType())
	s.router.ServeHTTP(s.response, s.context.Request)

	assert.Equal(s.T(), http.StatusOK, s.response.Code)

	json.NewDecoder(s.response.Body).Decode(&receivedResponse)

	assert.Equal(s.T(), ugtGroup.ID, receivedResponse["id"])
	assert.Equal(s.T(), ugtGroup.Name, receivedResponse["name"])
	assert.Equal(s.T(), ugtGroup.Description, receivedResponse["description"])
	assert.Equal(s.T(), ugtGroup.Creator.Username, receivedResponse["creator"])
	assert.Equal(s.T(), float64(ugtGroup.MemberCount), receivedResponse["member_count"])
	assert.Equal(s.T(), ugtGroup.IsOpen, receivedResponse["is_open"])

	images, isExist := receivedResponse["images"].([]interface{})
	assert.True(s.T(), isExist)
	image1 := images[0].(map[string]interface{})
	width, isExist := image1["width"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), float64(ugtGroup.Images[0].Width), width)
	height, isExist := image1["height"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), float64(ugtGroup.Images[0].Height), height)
	url, isExist := image1["url"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), ugtGroup.Images[0].URL, url)

	image2 := images[1].(map[string]interface{})
	width, isExist = image2["width"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), float64(ugtGroup.Images[1].Width), width)
	height, isExist = image2["height"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), float64(ugtGroup.Images[1].Height), height)
	url, isExist = image2["url"]
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), ugtGroup.Images[1].URL, url)
}
