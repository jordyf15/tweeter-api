package token

import "github.com/jordyf15/tweeter-api/models"

const (
	DefaultTokenLimitPerUser = 5
)

type Repository interface {
	GetTokenSet(userID string, hashedRefreshTokenID string, includeParent bool) (*models.TokenSet, error)

	Save(accessToken *models.AccessToken) error
	Exists(accessToken *models.AccessToken) bool
	Remove(accessToken *models.AccessToken) error

	Create(tokenSet *models.TokenSet) error
	Update(tokenSet *models.TokenSet) error
	Updates(tokenSet *models.TokenSet, changes map[string]interface{}) error
	Delete(tokenSet *models.TokenSet) error

	LimitTokenCount(userID string, limit uint) error
}

type Usecase interface {
	Refresh(token *models.RefreshToken) (*models.AccessToken, error)

	Use(token *models.AccessToken) error

	DeleteRefreshToken(token *models.RefreshToken) error
}
