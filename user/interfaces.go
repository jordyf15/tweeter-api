package user

import "github.com/jordyf15/tweeter-api/models"

var (
	ProfilePictureSizes = []uint{100, 400}
	BannerPictureWidth  = uint(1500)
	BannerPictureHeight = uint(500)
)

type Usecase interface {
	For(user *models.User) InstanceUsecase
	Create(*models.User) (map[string]interface{}, error)
	Login(login, password string) (map[string]interface{}, error)
}

type InstanceUsecase interface {
	GenerateTokens() (*models.AccessToken, *models.RefreshToken, error)
}

type Repository interface {
	Create(user *models.User) error
	CreateTransaction(fn func(repo Repository) error) error
	GetByEmailOrUsername(str string) (*models.User, error)
}
