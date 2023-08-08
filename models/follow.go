package models

import "time"

type Follow struct {
	FollowingID string `json:"-" gorm:"primaryKey"`

	FollowerID string `json:"-" gorm:"primaryKey"`

	CreatedAt time.Time `json:"-"`
}
