package follow

type Repository interface {
	Create(followerID, followingID string) error
}

type Usecase interface {
	FollowUser(followerID, followingID string) error
}
