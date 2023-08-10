package usecase_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/jordyf15/tweeter-api/custom_errors"
	"github.com/jordyf15/tweeter-api/group"
	groupMocks "github.com/jordyf15/tweeter-api/group/mocks"
	"github.com/jordyf15/tweeter-api/group/usecase"
	groupMemberMocks "github.com/jordyf15/tweeter-api/group_member/mocks"
	"github.com/jordyf15/tweeter-api/models"
	storageMocks "github.com/jordyf15/tweeter-api/storage/mocks"
	userMocks "github.com/jordyf15/tweeter-api/user/mocks"
	"github.com/jordyf15/tweeter-api/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestGroupUsecase(t *testing.T) {
	suite.Run(t, new(groupUsecaseSuite))
}

type groupUsecaseSuite struct {
	suite.Suite
	usecase         group.Usecase
	userRepo        *userMocks.Repository
	groupRepo       *groupMocks.Repository
	groupMemberRepo *groupMemberMocks.Repository
	storageMock     *storageMocks.Storage
}

var (
	utGroup1 = &models.Group{
		Name:        "hololive",
		Description: "vtuber comedian idol group",
		IsOpen:      true,
	}
)

func (s *groupUsecaseSuite) SetupTest() {
	s.userRepo = new(userMocks.Repository)
	s.groupRepo = new(groupMocks.Repository)
	s.groupMemberRepo = new(groupMemberMocks.Repository)
	s.storageMock = new(storageMocks.Storage)

	s.groupRepo.On("CreateTransaction", mock.Anything).Return(nil)
	s.storageMock.On("AssignImageURLToGroup", mock.AnythingOfType("*models.Group"))

	s.usecase = usecase.NewGroupUsecase(s.groupRepo, s.groupMemberRepo, s.userRepo, s.storageMock)
}

func (s *groupUsecaseSuite) TestCreateGroupNameTooShort() {
	group := &models.Group{
		Name: "a",
	}

	imgFile, _ := os.Open("../../assets/images/default-profile.png")
	defer imgFile.Close()
	imgFileReader := utils.NewNamedFileReader(imgFile, fmt.Sprintf("%s.%s", utils.RandString(8), utils.GetFileExtension(imgFile.Name())))

	expectedErrors := &custom_errors.MultipleErrors{Errors: []error{custom_errors.ErrGroupNameTooShort}}
	result, err := s.usecase.Create(group, imgFileReader)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), expectedErrors.Error(), err.Error())
	s.groupRepo.AssertNumberOfCalls(s.T(), "CreateTransaction", 0)
	s.storageMock.AssertNumberOfCalls(s.T(), "AssignImageURLToGroup", 0)
}

func (s *groupUsecaseSuite) TestCreateGroupNameTooLong() {
	group := &models.Group{
		Name: "asdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdq",
	}

	imgFile, _ := os.Open("../../assets/images/default-profile.png")
	defer imgFile.Close()
	imgFileReader := utils.NewNamedFileReader(imgFile, fmt.Sprintf("%s.%s", utils.RandString(8), utils.GetFileExtension(imgFile.Name())))

	expectedErrors := &custom_errors.MultipleErrors{Errors: []error{custom_errors.ErrGroupNameTooLong}}
	result, err := s.usecase.Create(group, imgFileReader)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), expectedErrors.Error(), err.Error())
	s.groupRepo.AssertNumberOfCalls(s.T(), "CreateTransaction", 0)
	s.storageMock.AssertNumberOfCalls(s.T(), "AssignImageURLToGroup", 0)
}

func (s *groupUsecaseSuite) TestCreateGroupImageInvalidFormat() {
	imgFile, _ := os.Open("../../assets/images/test_pic.gif")
	defer imgFile.Close()
	imgFileReader := utils.NewNamedFileReader(imgFile, fmt.Sprintf("%s.%s", utils.RandString(8), utils.GetFileExtension(imgFile.Name())))

	expectedErrors := &custom_errors.MultipleErrors{Errors: []error{custom_errors.ErrGroupImageInvalidFormat}}
	result, err := s.usecase.Create(utGroup1, imgFileReader)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), expectedErrors.Error(), err.Error())
	s.groupRepo.AssertNumberOfCalls(s.T(), "CreateTransaction", 0)
	s.storageMock.AssertNumberOfCalls(s.T(), "AssignImageURLToGroup", 0)
}

func (s *groupUsecaseSuite) TestCreateGroupSuccessful() {
	imgFile, _ := os.Open("../../assets/images/default-profile.png")
	defer imgFile.Close()
	imgFileReader := utils.NewNamedFileReader(imgFile, fmt.Sprintf("%s.%s", utils.RandString(8), utils.GetFileExtension(imgFile.Name())))

	result, err := s.usecase.Create(utGroup1, imgFileReader)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	s.groupRepo.AssertNumberOfCalls(s.T(), "CreateTransaction", 1)
	s.storageMock.AssertNumberOfCalls(s.T(), "AssignImageURLToGroup", 1)
}
