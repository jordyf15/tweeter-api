package models

import "time"

type GroupMemberRole string

const (
	GroupMemberRoleMember    GroupMemberRole = "member"
	GroupMemberRoleModerator GroupMemberRole = "moderator"
	GroupMemberRoleAdmin     GroupMemberRole = "admin"
)

type GroupMember struct {
	GroupID  string `json:"-" gorm:"primaryKey"`
	MemberID string `json:"-" gorm:"primaryKey"`

	Role GroupMemberRole `json:"role"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
