package usecase_test

import (
	"testing"

	"github.com/jordyf15/tweeter-api/custom_errors"
	"github.com/jordyf15/tweeter-api/follow"
	followMocks "github.com/jordyf15/tweeter-api/follow/mocks"
	"github.com/jordyf15/tweeter-api/follow/usecase"
	userMocks "github.com/jordyf15/tweeter-api/user/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestFollowUsecase(t *testing.T) {
	suite.Run(t, new(followUsecaseSuite))
}

type followUsecaseSuite struct {
	suite.Suite
	usecase    follow.Usecase
	userRepo   *userMocks.Repository
	followRepo *followMocks.Repository
}

func (s *followUsecaseSuite) SetupTest() {
	s.userRepo = new(userMocks.Repository)
	s.followRepo = new(followMocks.Repository)

	isIdExist := func(userID string) bool {
		return userID != "userID3"
	}

	s.userRepo.On("IsIDExist", mock.AnythingOfType("string")).Return(isIdExist, nil)
	s.followRepo.On("Create", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

	s.usecase = usecase.NewFollowUsecase(s.followRepo, s.userRepo)
}

func (s *followUsecaseSuite) TestFollowUserMatchedFollowerIDAndFollowingID() {
	err := s.usecase.FollowUser("userID1", "userID1")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrMatchedFollowerIDAndFollowingID.Error(), err.Error())
}

func (s *followUsecaseSuite) TestFollowUserFollowingIDNotExist() {
	err := s.usecase.FollowUser("userID1", "userID3")

	assert.Error(s.T(), err)
	assert.Equal(s.T(), custom_errors.ErrRecordNotFound.Error(), err.Error())
}

func (s *followUsecaseSuite) TestFollowUserSuccessful() {
	err := s.usecase.FollowUser("userID1", "userID2")

	assert.NoError(s.T(), err)
}
