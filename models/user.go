package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/jordyf15/tweeter-api/custom_errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	minUsernameLength = 3
	maxUsernameLength = 30
	minFullnameLength = 1
	maxFullnameLength = 255
	minPasswordLength = 8
	maxPasswordLength = 30
)

var (
	emailRegex      = regexp.MustCompile("\\A[\\w+\\-.]+@[a-z\\d\\-.]+\\.[a-z]+\\z")
	passwordRegexes = []*regexp.Regexp{
		regexp.MustCompile(".*[A-Z].*"),           // uppercase letter
		regexp.MustCompile(".*[a-z].*"),           // lowercase letter
		regexp.MustCompile(".*[0-9].*"),           // digit
		regexp.MustCompile(".*[^A-Za-z0-9\\s].*"), //special character
	}
	usernameRegexes = []regexPair{
		{regex: regexp.MustCompile("^[a-z0-9._]+$"), shouldMatch: true},    // characters allowed
		{regex: regexp.MustCompile("^[^_.].+$"), shouldMatch: true},        // must not start with a fullstop or underscore
		{regex: regexp.MustCompile("\\A.*[_.]{2}.*$"), shouldMatch: false}, // must not have consecutive fullstops/unserscores
		{regex: regexp.MustCompile("\\A.*[^_.]$"), shouldMatch: true},      // must not end with a fullstop or underscore
	}
)

type regexPair struct {
	regex       *regexp.Regexp
	shouldMatch bool
}

type User struct {
	ID                string `json:"id" gorm:"primaryKey"`
	Fullname          string `json:"fullname" gorm:"type:varchar(255)"`
	Username          string `json:"username" gorm:"type:varchar(30)"`
	Email             string `json:"email" gorm:"type:varchar(255)"`
	Description       string `json:"description" gorm:"type:text"`
	Password          string `gorm:"-" json:"-"`
	EncryptedPassword string `gorm:"type:text" json:"-"`

	ProfileImages   Images `gorm:"type:jsonb;default:'[]'" json:"profile_images"`
	BackgroundImage Image  `gorm:"type:json;default:'{}'" json:"background_image"`

	FollowerCount  uint `gorm:"default:0" json:"follower_count"`
	FollowingCount uint `gorm:"default:0" json:"following_count"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (user *User) VerifyFields() []error {
	errors := make([]error, 0)
	if !emailRegex.MatchString(user.Email) {
		errors = append(errors, custom_errors.ErrEmailAddressInvalid)
	}

	if len(user.Fullname) < minFullnameLength {
		errors = append(errors, custom_errors.ErrFullnameTooShort)
	}

	if len(user.Fullname) > maxFullnameLength {
		errors = append(errors, custom_errors.ErrFullnameTooLong)
	}

	if len(user.Username) < minUsernameLength {
		errors = append(errors, custom_errors.ErrUsernameTooShort)
	}

	if len(user.Username) > maxUsernameLength {
		errors = append(errors, custom_errors.ErrUsernameTooLong)
	}

	for _, regexPair := range usernameRegexes {
		if regexPair.regex.MatchString(user.Username) != regexPair.shouldMatch {
			errors = append(errors, custom_errors.ErrUsernameInvalid)
			break
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func (user *User) SetPassword(newPassword string) error {
	if len(newPassword) < minPasswordLength {
		return custom_errors.ErrPasswordTooShort
	} else if len(newPassword) > maxPasswordLength {
		return custom_errors.ErrPasswordTooLong
	}

	for _, regex := range passwordRegexes {
		if !regex.MatchString(newPassword) {
			return custom_errors.ErrPasswordInvalid
		}
	}

	hashedNewPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	user.Password = ""
	user.EncryptedPassword = string(hashedNewPassword)

	return nil
}

func (user *User) BeforeSave(tx *gorm.DB) error {
	errors := []error{}
	tempUser := &[]User{}
	err := tx.Where("email = (?) OR username = (?)", user.Email, user.Username).Find(tempUser).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if len(*tempUser) > 0 {
		var usernameExist, emailExist bool
		for _, v := range *tempUser {
			if v.Email == user.Email && !emailExist {
				emailExist = true
				errors = append(errors, custom_errors.ErrEmailAlreadyExist)
			}
			if v.Username == user.Username && !usernameExist {
				usernameExist = true
				errors = append(errors, custom_errors.ErrUsernameAlreadyExist)
			}
			if usernameExist && emailExist {
				break
			}
		}
	}

	if len(errors) > 0 {
		return &custom_errors.MultipleErrors{Errors: errors}
	}

	return nil
}

func (user *User) ImagePath(image *Image) string {
	return fmt.Sprintf("uploads/users/%s/%s", user.ID, image.Filename)
}

func (user *User) AssignPointers(columns []string) []interface{} {
	pointers := make([]interface{}, len(columns))

	pointerMap := map[string]interface{}{
		"id":               &user.ID,
		"fullname":         &user.Fullname,
		"username":         &user.Username,
		"email":            &user.Email,
		"description":      &user.Description,
		"profile_images":   &user.ProfileImages,
		"background_image": &user.BackgroundImage,
		"follower_count":   &user.FollowerCount,
		"following_count":  &user.FollowingCount,
		"created_at":       &user.CreatedAt,
		"updated_at":       &user.UpdatedAt,
	}

	for i, column := range columns {
		pointers[i] = pointerMap[column]
	}

	return pointers
}

func (user *User) Value() (driver.Value, error) {
	return json.Marshal(user)
}

func (user *User) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	err := json.Unmarshal(b, &user)
	return err
}
