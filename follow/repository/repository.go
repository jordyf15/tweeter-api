package repository

import (
	"time"

	"github.com/jordyf15/tweeter-api/follow"
	"github.com/jordyf15/tweeter-api/models"
	"gorm.io/gorm"
)

type followRepository struct {
	DB *gorm.DB
}

func NewFollowRepo(db *gorm.DB) follow.Repository {
	return &followRepository{DB: db}
}

func (repo *followRepository) Create(followerID, followingID string) error {
	follow := &models.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
		CreatedAt:   time.Now(),
	}

	return repo.DB.Create(follow).Error
}
