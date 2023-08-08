package follow

type Repository interface {
	Create(followerID, followingID string) error
	Delete(followerID, followingID string) error
}

type Usecase interface {
	FollowUser(followerID, followingID string) error
	UnfollowUser(followerID, followingID string) error
}
