package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jordyf15/tweeter-api/models"
	"github.com/jordyf15/tweeter-api/token"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const contextTimeout = time.Second * 30

const (
	RedisKeyFreshAccessTokens = "fresh-access-tokens"
)

type tokenRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewTokenRepository(db *gorm.DB, redis *redis.Client) token.Repository {
	return &tokenRepository{db: db, redis: redis}
}

func (repo *tokenRepository) GetTokenSet(userID, hashedRefreshTokenID string, includeParent bool) (*models.TokenSet, error) {
	tokenSet := &models.TokenSet{}

	var err error
	if includeParent {
		err = repo.db.Where("user_id = (?) AND (rt_id = (?) OR prt_id = (?))", userID, hashedRefreshTokenID).First(tokenSet).Error
	} else {
		err = repo.db.Where("user_id = (?) AND rt_id = (?)", userID, hashedRefreshTokenID).First(tokenSet).Error
	}

	return tokenSet, err
}

func (repo *tokenRepository) Save(accessToken *models.AccessToken) error {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	return repo.redis.ZAdd(ctx, RedisKeyFreshAccessTokens, redis.Z{Score: float64(accessToken.ExpiresAt), Member: accessToken.Id}).Err()
}

func (repo *tokenRepository) Exists(accessToken *models.AccessToken) bool {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	expiry := repo.redis.ZScore(ctx, RedisKeyFreshAccessTokens, accessToken.Id)
	return expiry != nil && expiry.Val() > 0
}

func (repo *tokenRepository) Remove(accessToken *models.AccessToken) error {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	return repo.redis.ZRem(ctx, RedisKeyFreshAccessTokens, accessToken.Id).Err()
}

func (repo *tokenRepository) Create(tokenSet *models.TokenSet) error {
	return repo.db.Create(tokenSet).Error
}

func (repo *tokenRepository) Update(tokenSet *models.TokenSet) error {
	return repo.db.Model(tokenSet).Updates(tokenSet).Error
}

func (repo *tokenRepository) Updates(tokenSet *models.TokenSet, changes map[string]interface{}) error {
	return repo.db.Model(tokenSet).Updates(changes).Error
}

func (repo *tokenRepository) Delete(tokenSet *models.TokenSet) error {
	if len(tokenSet.ID) == 0 {
		if len(tokenSet.UserID) > 0 && len(tokenSet.RefreshTokenID) > 0 {
			return repo.db.Delete(models.TokenSet{}, "user_id = (?) AND rt_id = (?)", tokenSet.UserID, tokenSet.RefreshTokenID).Error
		}

		return errors.New("token set has no primary key")
	}

	return repo.db.Delete(tokenSet).Error
}

func (repo *tokenRepository) LimitTokenCount(userID string, limit uint) error {
	return repo.db.Delete(models.TokenSet{}, `user_id = (?) AND 
	updated_at < (SELECT updated_at FROM token_sets WHERE user_id = (?)
	ORDER BY updated_at DESC OFFSET (?) LIMIT1)`, userID, userID, limit-1).Error
}
