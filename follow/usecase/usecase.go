package usecase

import (
	"github.com/jordyf15/tweeter-api/custom_errors"
	"github.com/jordyf15/tweeter-api/follow"
	"github.com/jordyf15/tweeter-api/user"
)

type followUsecase struct {
	followRepo follow.Repository
	userRepo   user.Repository
}

func NewFollowUsecase(followRepo follow.Repository, userRepo user.Repository) follow.Usecase {
	return &followUsecase{followRepo: followRepo, userRepo: userRepo}
}

func (usecase *followUsecase) FollowUser(followerID, followingID string) error {
	if followerID == followingID {
		return custom_errors.ErrMatchedFollowerIDAndFollowingID
	}

	isExist, err := usecase.userRepo.IsIDExist(followingID)
	if err != nil {
		return err
	}

	if !isExist {
		return custom_errors.ErrRecordNotFound
	}

	return usecase.followRepo.Create(followerID, followingID)
}
