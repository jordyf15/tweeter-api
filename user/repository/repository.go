package repository

import (
	"github.com/jordyf15/tweeter-api/models"
	"github.com/jordyf15/tweeter-api/user"
	"gorm.io/gorm"
)

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.Repository {
	return &userRepository{DB: db}
}

func (repo *userRepository) Create(user *models.User) error {
	return repo.DB.Create(user).Error
}

func (repo *userRepository) CreateTransaction(fn func(repo user.Repository) error) error {
	return repo.DB.Transaction(func(tx *gorm.DB) error {
		return fn(&userRepository{DB: tx})
	})
}

func (repo *userRepository) GetByEmailOrUsername(str string) (*models.User, error) {
	user := &models.User{}
	err := repo.DB.Table("users").Where("email = (?) OR username = (?)", str, str).First(user).Error

	return user, err
}

func (repo *userRepository) GetByID(id string) (*models.User, error) {
	user := &models.User{}

	err := repo.DB.Table("users").Where("id = ?", id).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *userRepository) Update(user *models.User) error {
	return repo.DB.Model(user).Select("*").Updates(user).Error
}
