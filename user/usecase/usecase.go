package usecase

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jordyf15/tweeter-api/custom_errors"
	"github.com/jordyf15/tweeter-api/models"
	"github.com/jordyf15/tweeter-api/storage"
	"github.com/jordyf15/tweeter-api/token"
	"github.com/jordyf15/tweeter-api/user"
	"github.com/jordyf15/tweeter-api/utils"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo  user.Repository
	tokenRepo token.Repository
	storage   storage.Storage
}

type userInstanceUsecase struct {
	user *models.User
	userUsecase
}

func NewUserUsecase(userRepo user.Repository, tokenRepo token.Repository, storage storage.Storage) user.Usecase {
	return &userUsecase{userRepo: userRepo, tokenRepo: tokenRepo, storage: storage}
}

func (usecase *userUsecase) For(user *models.User) user.InstanceUsecase {
	instanceUsecase := &userInstanceUsecase{user: user, userUsecase: *usecase}
	return instanceUsecase
}

func (usecase *userUsecase) Create(_user *models.User) (map[string]interface{}, error) {
	var err error
	errors := make([]error, 0)
	validateFieldErrors := _user.VerifyFields()
	if len(validateFieldErrors) > 0 {
		errors = append(errors, validateFieldErrors...)
	}

	err = _user.SetPassword(_user.Password)
	if err != nil {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return nil, &custom_errors.MultipleErrors{Errors: errors}
	}

	err = usecase.userRepo.CreateTransaction(func(repo user.Repository) error {
		_user.ID = uuid.New().String()

		uploadChannels := make(chan error, 3)
		var wg sync.WaitGroup
		wg.Add(3)
		defaultProfileImgFile, err := os.Open("./assets/images/default-profile.png")
		if err != nil {
			return err
		}
		defer defaultProfileImgFile.Close()

		profileImgNamedFileReader := utils.NewNamedFileReader(defaultProfileImgFile, fmt.Sprintf("%s.%s", utils.RandString(8), utils.GetFileExtension(defaultProfileImgFile.Name())))
		_user.ProfileImages = make([]*models.Image, len(user.ProfilePictureSizes))
		for i, width := range user.ProfilePictureSizes {
			image := &models.Image{}
			filename := utils.RandFileName("", "."+utils.GetFileExtension(defaultProfileImgFile.Name()))
			image.Filename = filename
			image.Width = width
			image.Height = width

			_user.ProfileImages[i] = image
		}

		defaultBannerImgFile, _ := os.Open("./assets/images/default-banner.jpg")
		defer defaultBannerImgFile.Close()

		_user.BackgroundImage = models.Image{
			Filename: utils.RandFileName("", "."+utils.GetFileExtension(defaultBannerImgFile.Name())),
			Width:    user.BannerPictureWidth,
			Height:   user.BannerPictureHeight,
		}

		err = usecase.userRepo.Create(_user)
		if err != nil {
			return err
		}

		for _, img := range _user.ProfileImages {
			resizedImageFile, err := utils.ResizeImage(profileImgNamedFileReader, int(img.Width), int(img.Height))
			if err != nil {
				return err
			}

			defer os.Remove(resizedImageFile.Name())
			go usecase.storage.UploadFile(uploadChannels, &wg, resizedImageFile, _user.ImagePath(img), nil)
		}

		go usecase.storage.UploadFile(uploadChannels, &wg, defaultBannerImgFile, _user.ImagePath(&_user.BackgroundImage), nil)

		wg.Wait()
		close(uploadChannels)

		for err := range uploadChannels {
			if err != nil {
				fmt.Println(err)
				return err
			}
		}

		return <-uploadChannels
	})

	if err != nil {
		switch actualErr := err.(type) {
		case *custom_errors.MultipleErrors:
			errors = append(errors, actualErr.Errors...)
		default:
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return nil, &custom_errors.MultipleErrors{Errors: errors}
	}

	usecase.storage.AssignImageURLToUser(_user)

	accessToken, refreshToken, _ := usecase.For(_user).GenerateTokens()

	response := utils.DataResponse(_user, map[string]interface{}{
		"access_token":  accessToken.ToJWTString(),
		"refresh_token": refreshToken.ToJWTString(),
		"expires_at":    accessToken.ExpiresAt,
	})

	return response, nil
}

func (usecase *userUsecase) Login(login, password string) (map[string]interface{}, error) {
	user, err := usecase.userRepo.GetByEmailOrUsername(login)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, custom_errors.ErrPasswordIncorrect
		}
		return nil, err
	}

	accessToken, refreshToken, err := usecase.For(user).GenerateTokens()
	if err != nil {
		return nil, err
	}

	usecase.storage.AssignImageURLToUser(user)

	response := utils.DataResponse(user, map[string]interface{}{
		"access_token":  accessToken.ToJWTString(),
		"refresh_token": refreshToken.ToJWTString(),
		"expires_at":    accessToken.ExpiresAt,
	})

	return response, nil
}

func (usecase *userUsecase) EditUserProfile(userID string, updates map[string]string, profileImageReader, backgroundImageReader utils.NamedFileReader, willRemoveProfileImage, willRemoveBackgroundImage bool) (*models.User, error) {
	errors := make([]error, 0)

	_user, err := usecase.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	if newFullname, isExist := updates["fullname"]; isExist && newFullname != _user.Fullname {
		_user.Fullname = newFullname
	}

	if newUsername, isExist := updates["username"]; isExist && newUsername != _user.Username {
		_user.Username = newUsername
	}

	if newEmail, isExist := updates["email"]; isExist && newEmail != _user.Email {
		_user.Email = newEmail
	}

	if newDescription, isExist := updates["description"]; isExist && newDescription != _user.Description {
		_user.Description = newDescription
	}

	validateFieldErrors := _user.VerifyFields()
	if len(validateFieldErrors) > 0 {
		errors = append(errors, validateFieldErrors...)
	}

	if profileImageReader != nil {
		switch utils.GetFileExtension(profileImageReader.Name()) {
		case "jpg", "jpeg", "png":
			break
		default:
			errors = append(errors, custom_errors.ErrProfileImageInvalidFormat)
		}
	}

	if backgroundImageReader != nil {
		switch utils.GetFileExtension(backgroundImageReader.Name()) {
		case "jpg", "jpeg", "png":
			break
		default:
			errors = append(errors, custom_errors.ErrBackgroundImageInvalidFormat)
		}
	}

	if len(errors) > 0 {
		return nil, &custom_errors.MultipleErrors{Errors: errors}
	}

	previousProfileImages := _user.ProfileImages
	previousBackgroundImage := _user.BackgroundImage
	var newProfileImage utils.NamedFileReader
	var newBackgroundImage utils.NamedFileReader

	if profileImageReader != nil && !willRemoveProfileImage {
		newProfileImage = profileImageReader
		_user.ProfileImages = make([]*models.Image, len(user.ProfilePictureSizes))

		for idx, res := range user.ProfilePictureSizes {
			image := &models.Image{}
			filename := utils.RandFileName("", "."+utils.GetFileExtension(profileImageReader.Name()))
			image.Filename = filename
			image.Width = res
			image.Height = res

			_user.ProfileImages[idx] = image
		}
	}

	if backgroundImageReader != nil && !willRemoveBackgroundImage {
		newBackgroundImage = backgroundImageReader
		image := models.Image{}
		image.Filename = utils.RandFileName("", "."+utils.GetFileExtension(backgroundImageReader.Name()))
		image.Width = user.BannerPictureWidth
		image.Height = user.BannerPictureHeight

		_user.BackgroundImage = image
	}

	if willRemoveProfileImage && profileImageReader == nil {
		defaultProfileImgFile, err := os.Open("./assets/images/default-profile.png")
		if err != nil {
			return nil, err
		}
		defer defaultProfileImgFile.Close()

		newProfileImage = utils.NewNamedFileReader(defaultProfileImgFile, fmt.Sprintf("%s.%s", utils.RandString(8), utils.GetFileExtension(defaultProfileImgFile.Name())))
		_user.ProfileImages = make([]*models.Image, len(user.ProfilePictureSizes))
		for i, width := range user.ProfilePictureSizes {
			image := &models.Image{}
			filename := utils.RandFileName("", "."+utils.GetFileExtension(defaultProfileImgFile.Name()))
			image.Filename = filename
			image.Width = width
			image.Height = width

			_user.ProfileImages[i] = image
		}
	}

	if willRemoveBackgroundImage && backgroundImageReader == nil {
		defaultBannerImgFile, _ := os.Open("./assets/images/default-banner.jpg")
		defer defaultBannerImgFile.Close()

		newBackgroundImage = utils.NewNamedFileReader(defaultBannerImgFile, fmt.Sprintf("%s.%s", utils.RandString(8), utils.GetFileExtension(defaultBannerImgFile.Name())))
		_user.BackgroundImage = models.Image{
			Filename: utils.RandFileName("", "."+utils.GetFileExtension(defaultBannerImgFile.Name())),
			Width:    user.BannerPictureWidth,
			Height:   user.BannerPictureHeight,
		}
	}

	err = usecase.userRepo.CreateTransaction(func(repo user.Repository) error {
		err = usecase.userRepo.Update(_user)
		if err != nil {
			return err
		}

		var wg sync.WaitGroup
		fileStorageChannelSizes := 0

		if willRemoveBackgroundImage || backgroundImageReader != nil {
			// 1 for saving and 1 for deleting
			fileStorageChannelSizes += 2
			wg.Add(2)
		}
		if willRemoveProfileImage || profileImageReader != nil {
			// 2 for saving and 2 for deleting
			fileStorageChannelSizes += 4
			wg.Add(4)
		}

		fileStorageChannels := make(chan error, fileStorageChannelSizes)

		if willRemoveBackgroundImage || backgroundImageReader != nil {
			resizedImageFile, err := utils.ResizeImage(newBackgroundImage, int(user.BannerPictureWidth), int(user.BannerPictureHeight))
			if err != nil {
				return err
			}

			defer os.Remove(resizedImageFile.Name())
			go usecase.storage.UploadFile(fileStorageChannels, &wg, resizedImageFile, _user.ImagePath(&_user.BackgroundImage), nil)

			go usecase.storage.RemoveFile(fileStorageChannels, &wg, _user.ImagePath(&previousBackgroundImage))
		}

		if willRemoveProfileImage || profileImageReader != nil {
			for _, img := range _user.ProfileImages {
				resizedImageFile, err := utils.ResizeImage(newProfileImage, int(img.Width), int(img.Height))
				if err != nil {
					return err
				}

				defer os.Remove(resizedImageFile.Name())
				go usecase.storage.UploadFile(fileStorageChannels, &wg, resizedImageFile, _user.ImagePath(img), nil)
			}

			for _, img := range previousProfileImages {
				go usecase.storage.RemoveFile(fileStorageChannels, &wg, _user.ImagePath(img))
			}
		}

		wg.Wait()
		close(fileStorageChannels)

		for err := range fileStorageChannels {
			if err != nil {
				return err
			}
		}

		return <-fileStorageChannels
	})

	if err != nil {
		switch actualErr := err.(type) {
		case *custom_errors.MultipleErrors:
			errors = append(errors, actualErr.Errors...)
		default:
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return nil, &custom_errors.MultipleErrors{Errors: errors}
	}

	usecase.storage.AssignImageURLToUser(_user)

	return _user, nil
}

func (usecase *userUsecase) ChangeUserPassword(userId, oldPassword, newPassword string) error {
	user, err := usecase.userRepo.GetByID(userId)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(oldPassword))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return custom_errors.ErrPasswordIncorrect
		}
		return err
	}

	err = user.SetPassword(newPassword)
	if err != nil {
		return err
	}

	err = usecase.userRepo.Update(user)
	if err != nil {
		return err
	}

	return nil
}

func (usecase *userInstanceUsecase) GenerateTokens() (*models.AccessToken, *models.RefreshToken, error) {
	refreshToken := (&models.RefreshToken{UserID: usecase.user.ID})
	refreshToken.Id = utils.RandString(8)
	accessToken := (&models.AccessToken{UserID: usecase.user.ID}).SetExpiration(time.Now().Add(time.Hour * 1))

	tokenSet := &models.TokenSet{UserID: usecase.user.ID, RefreshTokenID: accessToken.RefreshTokenID}
	err := usecase.tokenRepo.Create(tokenSet)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}
