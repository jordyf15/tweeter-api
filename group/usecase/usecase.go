package usecase

import (
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/jordyf15/tweeter-api/custom_errors"
	"github.com/jordyf15/tweeter-api/group"
	"github.com/jordyf15/tweeter-api/group_member"
	"github.com/jordyf15/tweeter-api/models"
	"github.com/jordyf15/tweeter-api/storage"
	"github.com/jordyf15/tweeter-api/user"
	"github.com/jordyf15/tweeter-api/utils"
)

type groupUsecase struct {
	groupRepo       group.Repository
	userRepo        user.Repository
	groupMemberRepo group_member.Repository
	storage         storage.Storage
}

func NewGroupUsecase(groupRepo group.Repository, groupMemberRepo group_member.Repository, userRepo user.Repository, storage storage.Storage) group.Usecase {
	return &groupUsecase{groupRepo: groupRepo, groupMemberRepo: groupMemberRepo, userRepo: userRepo, storage: storage}
}

func (usecase *groupUsecase) Create(_group *models.Group, groupImageReader utils.NamedFileReader) (*models.Group, error) {
	var err error
	errors := make([]error, 0)

	validateFieldErrors := _group.VerifyFields()
	if len(validateFieldErrors) > 0 {
		errors = append(errors, validateFieldErrors...)
	}

	switch utils.GetFileExtension(groupImageReader.Name()) {
	case "jpg", "jpeg", "png":
		break
	default:
		errors = append(errors, custom_errors.ErrGroupImageInvalidFormat)
	}

	if len(errors) > 0 {
		return nil, &custom_errors.MultipleErrors{Errors: errors}
	}

	err = usecase.groupRepo.CreateTransaction(func(repo group.Repository) error {
		_group.ID = uuid.New().String()
		_group.Images = make([]*models.Image, 2)

		uploadChannels := make(chan error, 2)
		var wg sync.WaitGroup
		wg.Add(2)
		thumbnailImg := &models.Image{}
		thumbnailImg.Filename = utils.RandFileName("", "."+utils.GetFileExtension(groupImageReader.Name()))
		thumbnailImg.Width = group.ThumbnailPictureRes
		thumbnailImg.Height = group.ThumbnailPictureRes
		_group.Images[0] = thumbnailImg

		bannerImg := &models.Image{}
		bannerImg.Filename = utils.RandFileName("", "."+utils.GetFileExtension(groupImageReader.Name()))
		bannerImg.Width = group.BannerPictureWidth
		bannerImg.Height = group.BannerPictureHeight
		_group.Images[1] = bannerImg

		_group.Creator, err = usecase.userRepo.GetByID(_group.CreatorID)
		if err != nil {
			return err
		}

		err = usecase.groupRepo.Create(_group)
		if err != nil {
			return err
		}

		groupMember := &models.GroupMember{
			GroupID:  _group.ID,
			MemberID: _group.CreatorID,
			Role:     models.GroupMemberRoleAdmin,
		}

		err = usecase.groupMemberRepo.Create(groupMember)
		if err != nil {
			return err
		}

		resizedThumbnailImg, err := utils.ResizeImage(groupImageReader, int(thumbnailImg.Width), int(thumbnailImg.Height))
		if err != nil {
			return err
		}
		defer os.Remove(resizedThumbnailImg.Name())
		go usecase.storage.UploadFile(uploadChannels, &wg, resizedThumbnailImg, _group.ImagePath(thumbnailImg), nil)

		resizedBannerImg, err := utils.ResizeImage(groupImageReader, int(bannerImg.Width), int(bannerImg.Height))
		if err != nil {
			return err
		}
		defer os.Remove(resizedBannerImg.Name())
		go usecase.storage.UploadFile(uploadChannels, &wg, resizedBannerImg, _group.ImagePath(bannerImg), nil)

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

	usecase.storage.AssignImageURLToGroup(_group)
	_group.MemberCount = 1

	return _group, nil
}
