package group

import (
	"github.com/jordyf15/tweeter-api/models"
	"github.com/jordyf15/tweeter-api/utils"
)

var (
	ThumbnailPictureRes = uint(400)
	BannerPictureWidth  = uint(900)
	BannerPictureHeight = uint(350)
)

type Usecase interface {
	Create(group *models.Group, groupImage utils.NamedFileReader) (*models.Group, error)
}

type Repository interface {
	Create(group *models.Group) error
	CreateTransaction(fn func(repo Repository) error) error
}
