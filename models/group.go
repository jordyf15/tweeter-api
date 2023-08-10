package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jordyf15/tweeter-api/custom_errors"
	"gorm.io/gorm"
)

const (
	minNameLength = 3
	maxNameLength = 255
)

type Group struct {
	ID          string `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"type:varchar(255)"`
	Description string `json:"description" gorm:"type:text"`
	Images      `gorm:"type:jsonb,default:'[]'" json:"images"`
	MemberCount uint      `gorm:"default:0" json:"member_count"`
	CreatorID   string    `json:"-"`
	Creator     *User     `gorm:"-" json:"creator"`
	CreatedAt   time.Time `json:"created_at"`
	IsOpen      bool      `json:"is_open"`
}

func (group *Group) VerifyFields() []error {
	errors := make([]error, 0)
	if len(group.Name) < minNameLength {
		errors = append(errors, custom_errors.ErrGroupNameTooShort)
	}

	if len(group.Name) > maxNameLength {
		errors = append(errors, custom_errors.ErrGroupNameTooLong)
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func (group *Group) BeforeSave(tx *gorm.DB) error {
	tempGroup := &[]Group{}
	err := tx.Where("id <> (?) and (name = (?))", group.ID, group.Name).Find(tempGroup).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if len(*tempGroup) > 0 {
		return custom_errors.ErrGroupNameAlreadyExist
	}

	return nil
}

func (group *Group) AssignPointers(columns []string) []interface{} {
	pointers := make([]interface{}, len(columns))

	pointerMap := map[string]interface{}{
		"id":           &group.ID,
		"created_at":   &group.CreatedAt,
		"name":         &group.Name,
		"description":  &group.Description,
		"images":       &group.Images,
		"member_count": &group.MemberCount,
		"creator_id":   &group.CreatorID,
		"is_open":      &group.IsOpen,
	}

	for i, column := range columns {
		pointers[i] = pointerMap[column]
	}

	return pointers
}

func (group *Group) MarshalJSON() ([]byte, error) {
	type Alias Group
	newStruct := &struct {
		Creator   string `json:"creator"`
		CreatedAt string `json:"created_at"`
		*Alias
	}{
		Creator:   group.Creator.Username,
		CreatedAt: group.CreatedAt.Format("2006-01-02T15:04:05-0700"),
		Alias:     (*Alias)(group),
	}

	return json.Marshal(newStruct)
}

func (group *Group) ImagePath(image *Image) string {
	return fmt.Sprintf("uploads/groups/%s/%s", group.ID, image.Filename)
}
