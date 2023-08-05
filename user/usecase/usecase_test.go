package usecase_test

import (
	"sync"
	"testing"

	"github.com/jordyf15/tweeter-api/custom_errors"
	"github.com/jordyf15/tweeter-api/models"
	storageMocks "github.com/jordyf15/tweeter-api/storage/mocks"
	tokenMocks "github.com/jordyf15/tweeter-api/token/mocks"
	"github.com/jordyf15/tweeter-api/user"
	userMocks "github.com/jordyf15/tweeter-api/user/mocks"
	"github.com/jordyf15/tweeter-api/user/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

func TestUserUsecase(t *testing.T) {
	suite.Run(t, new(userUsecaseSuite))
}

type userUsecaseSuite struct {
	suite.Suite
	usecase     user.Usecase
	userRepo    *userMocks.Repository
	tokenRepo   *tokenMocks.Repository
	storageMock *storageMocks.Storage
}

var (
	utUser = &models.User{
		ID:                "id",
		Email:             "gura@gmail.com",
		Fullname:          "gawr gura",
		Username:          "gura",
		Description:       "",
		EncryptedPassword: bcryptHash("Password123!"),
		FollowerCount:     0,
		FollowingCount:    0,
	}
)

func bcryptHash(str string) string {
	hashedStr, _ := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	return string(hashedStr)
}

func (s *userUsecaseSuite) SetupTest() {
	s.userRepo = new(userMocks.Repository)
	s.tokenRepo = new(tokenMocks.Repository)
	s.storageMock = new(storageMocks.Storage)

	s.storageMock.On("GetFileLink", mock.AnythingOfType("string")).Return("string", nil)
	s.storageMock.On("AssignImageURLToUser", mock.AnythingOfType("*models.User"))
	s.storageMock.On("UploadFile", mock.AnythingOfType("chan<- error"), mock.AnythingOfType("*sync.WaitGroup"), mock.AnythingOfType("*os.File"), mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).Run(func(args mock.Arguments) {
		arg1 := args[0].(chan<- error)
		arg1 <- nil
		arg2 := args[1].(*sync.WaitGroup)
		arg2.Done()
	})
	s.tokenRepo.On("Create", mock.AnythingOfType("*models.TokenSet")).Return(nil)
	s.userRepo.On("CreateTransaction", mock.Anything).Return(nil)
	s.userRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)
	s.userRepo.On("GetByEmailOrUsername", mock.AnythingOfType("string")).Return(utUser, nil)

	s.usecase = usecase.NewUserUsecase(s.userRepo, s.tokenRepo, s.storageMock)
}

func (s *userUsecaseSuite) TestCreateUsernameTooShort() {
	user := &models.User{
		Username: "as",
		Fullname: "gawr gura",
		Email:    "gura@gmail.com",
		Password: "Password123!",
	}

	expectedErrors := &custom_errors.MultipleErrors{Errors: []error{custom_errors.ErrUsernameTooShort}}
	result, err := s.usecase.Create(user)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), expectedErrors.Error(), err.Error())
	s.userRepo.AssertNumberOfCalls(s.T(), "CreateTransaction", 0)
	s.storageMock.AssertNumberOfCalls(s.T(), "AssignImageURLToUser", 0)
	s.tokenRepo.AssertNumberOfCalls(s.T(), "Create", 0)
}

func (s *userUsecaseSuite) TestCreateUsernameTooLong() {
	user := &models.User{
		Username: "asdasdasdqasdasdasdqasdasdasdqasdasdasdq",
		Fullname: "gawr gura",
		Email:    "gura@gmail.com",
		Password: "Password123!",
	}

	expectedErrors := &custom_errors.MultipleErrors{Errors: []error{custom_errors.ErrUsernameTooLong}}
	result, err := s.usecase.Create(user)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), expectedErrors.Error(), err.Error())
	s.userRepo.AssertNumberOfCalls(s.T(), "CreateTransaction", 0)
	s.storageMock.AssertNumberOfCalls(s.T(), "AssignImageURLToUser", 0)
	s.tokenRepo.AssertNumberOfCalls(s.T(), "Create", 0)
}

func (s *userUsecaseSuite) TestCreateUsernameInvalid() {
	user := &models.User{
		Username: "!!!",
		Fullname: "gawr gura",
		Email:    "gura@gmail.com",
		Password: "Password123!",
	}

	expectedErrors := &custom_errors.MultipleErrors{Errors: []error{custom_errors.ErrUsernameInvalid}}
	result, err := s.usecase.Create(user)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), expectedErrors.Error(), err.Error())
	s.userRepo.AssertNumberOfCalls(s.T(), "CreateTransaction", 0)
	s.storageMock.AssertNumberOfCalls(s.T(), "AssignImageURLToUser", 0)
	s.tokenRepo.AssertNumberOfCalls(s.T(), "Create", 0)
}

func (s *userUsecaseSuite) TestCreateFullnameTooShort() {
	user := &models.User{
		Username: "gura",
		Fullname: "",
		Email:    "gura@gmail.com",
		Password: "Password123!",
	}

	expectedErrors := &custom_errors.MultipleErrors{Errors: []error{custom_errors.ErrFullnameTooShort}}
	result, err := s.usecase.Create(user)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), expectedErrors.Error(), err.Error())
	s.userRepo.AssertNumberOfCalls(s.T(), "CreateTransaction", 0)
	s.storageMock.AssertNumberOfCalls(s.T(), "AssignImageURLToUser", 0)
	s.tokenRepo.AssertNumberOfCalls(s.T(), "Create", 0)
}

func (s *userUsecaseSuite) TestCreateFullnameTooLong() {
	user := &models.User{
		Username: "gura",
		Fullname: "asdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdq",
		Email:    "gura@gmail.com",
		Password: "Password123!",
	}

	expectedErrors := &custom_errors.MultipleErrors{Errors: []error{custom_errors.ErrFullnameTooLong}}
	result, err := s.usecase.Create(user)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), expectedErrors.Error(), err.Error())
	s.userRepo.AssertNumberOfCalls(s.T(), "CreateTransaction", 0)
	s.storageMock.AssertNumberOfCalls(s.T(), "AssignImageURLToUser", 0)
	s.tokenRepo.AssertNumberOfCalls(s.T(), "Create", 0)
}

func (s *userUsecaseSuite) TestCreateEmailInvalid() {
	user := &models.User{
		Username: "gura",
		Fullname: "gawr gura",
		Email:    "gura.com",
		Password: "Password123!",
	}

	expectedErrors := &custom_errors.MultipleErrors{Errors: []error{custom_errors.ErrEmailAddressInvalid}}
	result, err := s.usecase.Create(user)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), expectedErrors.Error(), err.Error())
	s.userRepo.AssertNumberOfCalls(s.T(), "CreateTransaction", 0)
	s.storageMock.AssertNumberOfCalls(s.T(), "AssignImageURLToUser", 0)
	s.tokenRepo.AssertNumberOfCalls(s.T(), "Create", 0)
}

func (s *userUsecaseSuite) TestCreatePasswordTooShort() {
	user := &models.User{
		Username: "gura",
		Fullname: "gawr gura",
		Email:    "gura@gmail.com",
		Password: "a",
	}

	expectedErrors := &custom_errors.MultipleErrors{Errors: []error{custom_errors.ErrPasswordTooShort}}
	result, err := s.usecase.Create(user)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), expectedErrors.Error(), err.Error())
	s.userRepo.AssertNumberOfCalls(s.T(), "CreateTransaction", 0)
	s.storageMock.AssertNumberOfCalls(s.T(), "AssignImageURLToUser", 0)
	s.tokenRepo.AssertNumberOfCalls(s.T(), "Create", 0)
}

func (s *userUsecaseSuite) TestCreatePasswordTooLong() {
	user := &models.User{
		Username: "gura",
		Fullname: "gawr gura",
		Email:    "gura@gmail.com",
		Password: "asdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdqasdasdasdq",
	}

	expectedErrors := &custom_errors.MultipleErrors{Errors: []error{custom_errors.ErrPasswordTooLong}}
	result, err := s.usecase.Create(user)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), expectedErrors.Error(), err.Error())
	s.userRepo.AssertNumberOfCalls(s.T(), "CreateTransaction", 0)
	s.storageMock.AssertNumberOfCalls(s.T(), "AssignImageURLToUser", 0)
	s.tokenRepo.AssertNumberOfCalls(s.T(), "Create", 0)
}

func (s *userUsecaseSuite) TestCreatePasswordInvalid() {
	user := &models.User{
		Username: "gura",
		Fullname: "gawr gura",
		Email:    "gura@gmail.com",
		Password: "password",
	}

	expectedErrors := &custom_errors.MultipleErrors{Errors: []error{custom_errors.ErrPasswordInvalid}}
	result, err := s.usecase.Create(user)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), expectedErrors.Error(), err.Error())
	s.userRepo.AssertNumberOfCalls(s.T(), "CreateTransaction", 0)
	s.storageMock.AssertNumberOfCalls(s.T(), "AssignImageURLToUser", 0)
	s.tokenRepo.AssertNumberOfCalls(s.T(), "Create", 0)
}

func (s *userUsecaseSuite) TestCreateSuccessful() {
	user := &models.User{
		Username: utUser.Username,
		Fullname: utUser.Fullname,
		Email:    utUser.Email,
		Password: "Password123!",
	}

	result, err := s.usecase.Create(user)

	assert.NoError(s.T(), err)

	data, isExist := result["data"].(*models.User)
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), utUser.Email, data.Email)
	assert.Equal(s.T(), utUser.Fullname, data.Fullname)
	assert.Equal(s.T(), utUser.Username, data.Username)
	assert.Equal(s.T(), utUser.FollowerCount, data.FollowerCount)
	assert.Equal(s.T(), utUser.FollowingCount, data.FollowingCount)
	assert.NotEmpty(s.T(), utUser.ID)

	meta, isExist := result["meta"].(map[string]interface{})
	assert.True(s.T(), isExist)

	_, isExist = meta["access_token"]
	assert.True(s.T(), isExist)

	_, isExist = meta["expires_at"]
	assert.True(s.T(), isExist)

	_, isExist = meta["refresh_token"]
	assert.True(s.T(), isExist)

	s.userRepo.AssertNumberOfCalls(s.T(), "CreateTransaction", 1)
	s.storageMock.AssertNumberOfCalls(s.T(), "AssignImageURLToUser", 1)
	s.tokenRepo.AssertNumberOfCalls(s.T(), "Create", 1)
}

func (s *userUsecaseSuite) TestLoginIncorrectPassword() {
	response, err := s.usecase.Login("gura", "wrongPassword")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrPasswordIncorrect.Error(), err.Error())
	assert.Nil(s.T(), response)

	s.userRepo.AssertNumberOfCalls(s.T(), "GetByEmailOrUsername", 1)
	s.tokenRepo.AssertNumberOfCalls(s.T(), "Create", 0)
	s.storageMock.AssertNumberOfCalls(s.T(), "AssignImageURLToUser", 0)
}

func (s *userUsecaseSuite) TestLoginSuccessful() {
	response, err := s.usecase.Login("gura", "Password123!")
	assert.NoError(s.T(), err)

	data, isExist := response["data"].(*models.User)
	assert.True(s.T(), isExist)
	assert.Equal(s.T(), utUser.Email, data.Email)
	assert.Equal(s.T(), utUser.Fullname, data.Fullname)
	assert.Equal(s.T(), utUser.Username, data.Username)
	assert.Equal(s.T(), utUser.FollowerCount, data.FollowerCount)
	assert.Equal(s.T(), utUser.FollowingCount, data.FollowingCount)
	assert.NotEmpty(s.T(), utUser.ID)

	meta, isExist := response["meta"].(map[string]interface{})
	assert.True(s.T(), isExist)

	_, isExist = meta["access_token"]
	assert.True(s.T(), isExist)

	_, isExist = meta["expires_at"]
	assert.True(s.T(), isExist)

	_, isExist = meta["refresh_token"]
	assert.True(s.T(), isExist)

	s.userRepo.AssertNumberOfCalls(s.T(), "GetByEmailOrUsername", 1)
	s.tokenRepo.AssertNumberOfCalls(s.T(), "Create", 1)
	s.storageMock.AssertNumberOfCalls(s.T(), "AssignImageURLToUser", 1)
}
