package repository

import (
	"github.com/jordyf15/tweeter-api/group_member"
	"github.com/jordyf15/tweeter-api/models"
	"gorm.io/gorm"
)

type groupMemberRepository struct {
	DB *gorm.DB
}

func NewGroupMemberRepository(db *gorm.DB) group_member.Repository {
	return &groupMemberRepository{DB: db}
}

func (repo *groupMemberRepository) Create(groupMember *models.GroupMember) error {
	return repo.DB.Create(groupMember).Error
}
