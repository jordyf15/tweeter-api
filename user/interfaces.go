package user

import (
	"github.com/jordyf15/tweeter-api/models"
	"github.com/jordyf15/tweeter-api/utils"
)

var (
	ProfilePictureSizes = []uint{100, 400}
	BannerPictureWidth  = uint(1500)
	BannerPictureHeight = uint(500)
)

type Usecase interface {
	For(user *models.User) InstanceUsecase
	Create(*models.User) (map[string]interface{}, error)
	Login(login, password string) (map[string]interface{}, error)
	ChangeUserPassword(userID, oldPassword, newPassword string) error
	EditUserProfile(userID string, updates map[string]string, profileImageReader, backgroundImageReader utils.NamedFileReader, willRemoveProfileImage, willRemoveBackgroundImage bool) (*models.User, error)
}

type InstanceUsecase interface {
	GenerateTokens() (*models.AccessToken, *models.RefreshToken, error)
}

type Repository interface {
	Create(user *models.User) error
	CreateTransaction(fn func(repo Repository) error) error
	GetByEmailOrUsername(str string) (*models.User, error)
	GetByID(id string) (*models.User, error)
	IsIDExist(id string) (bool, error)
	Update(user *models.User) error
}
