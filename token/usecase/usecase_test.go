package usecase_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/jordyf15/tweeter-api/models"
	"github.com/jordyf15/tweeter-api/token"
	tokenMocks "github.com/jordyf15/tweeter-api/token/mocks"
	"github.com/jordyf15/tweeter-api/token/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestTokenUsecase(t *testing.T) {
	suite.Run(t, new(tokenUsecaseSuite))
}

type tokenUsecaseSuite struct {
	suite.Suite
	usecase   token.Usecase
	tokenRepo *tokenMocks.Repository
}

func (s *tokenUsecaseSuite) SetupTest() {
	tokenId := "tokenId"
	userId := "userId"
	s.tokenRepo = new(tokenMocks.Repository)
	tokenSet := &models.TokenSet{
		ID:                 tokenId,
		UserID:             userId,
		RefreshTokenID:     "refreshTokenId",
		UpdatedAt:          time.Now(),
		PrevRefreshTokenID: nil,
	}

	s.tokenRepo.On("GetTokenSet", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("bool")).Return(tokenSet, nil)
	s.tokenRepo.On("Update", mock.AnythingOfType("*models.TokenSet")).Return(nil)
	s.tokenRepo.On("Save", mock.AnythingOfType("*models.AccessToken")).Return(nil)
	s.tokenRepo.On("Exists", mock.AnythingOfType("*models.AccessToken")).Return(true)
	s.tokenRepo.On("Updates", mock.AnythingOfType("*models.TokenSet"), mock.AnythingOfType("map[string]interface {}")).Return(nil)
	s.tokenRepo.On("Remove", mock.AnythingOfType("*models.AccessToken")).Return(nil)
	s.tokenRepo.On("Delete", mock.AnythingOfType("*models.TokenSet")).Return(nil)
	s.usecase = usecase.NewTokenUsecase(s.tokenRepo)
}

func (s *tokenUsecaseSuite) TestRefresh() {
	userId := "userId"
	refreshToken := &models.RefreshToken{
		UserID: userId,
	}
	accessToken, err := s.usecase.Refresh(refreshToken)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), accessToken.RefreshTokenID)
	assert.Equal(s.T(), "string", fmt.Sprintf("%T", accessToken.RefreshTokenID))
	assert.Equal(s.T(), userId, accessToken.UserID)
	assert.NotEmpty(s.T(), accessToken.Id)
}

func (s *tokenUsecaseSuite) TestUse() {
	userId := "userId"
	accessToken := &models.AccessToken{
		UserID:         userId,
		RefreshTokenID: "refreshTokenId",
	}
	err := s.usecase.Use(accessToken)
	assert.NoError(s.T(), err)
}

func (s *tokenUsecaseSuite) TestDeleteRefreshToken() {
	userId := "userId"
	refreshToken := &models.RefreshToken{
		UserID: userId,
	}
	err := s.usecase.DeleteRefreshToken(refreshToken)
	assert.NoError(s.T(), err)
}
