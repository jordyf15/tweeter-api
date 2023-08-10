package repository

import (
	"github.com/jordyf15/tweeter-api/group"
	"github.com/jordyf15/tweeter-api/models"
	"gorm.io/gorm"
)

type groupRepository struct {
	DB *gorm.DB
}

func NewGroupRepository(db *gorm.DB) group.Repository {
	return &groupRepository{DB: db}
}

func (repo *groupRepository) Create(group *models.Group) error {
	return repo.DB.Create(group).Error
}

func (repo *groupRepository) CreateTransaction(fn func(repo group.Repository) error) error {
	return repo.DB.Transaction(func(tx *gorm.DB) error {
		return fn(&groupRepository{DB: tx})
	})
}
