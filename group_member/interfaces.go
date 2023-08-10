package group_member

import "github.com/jordyf15/tweeter-api/models"

type Repository interface {
	Create(groupMember *models.GroupMember) error
}
