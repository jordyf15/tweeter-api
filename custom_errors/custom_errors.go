package custom_errors

import "strings"

var (
	// general errors
	// ErrRecordNotFound Error when related data is not found
	ErrRecordNotFound = newErr(101, "Record not found")
	// ErrUnknownErrorOccured Error when an unknown error has occured
	ErrUnknownErrorOccured = newErr(102, "Unknown error occured")

	//token errors
	// ErrMalformedRefreshToken Refresh token couldn't be parsed.
	ErrMalformedRefreshToken = newErr(201, "Refresh token is malformed")
	// ErrInvalidRefreshToken refresh token not signed on this server.
	ErrInvalidRefreshToken = newErr(202, "Invalid refresh token")
	// ErrRefreshTokenNotFound refresh token not found in database.
	ErrRefreshTokenNotFound = newErr(203, "Refresh token not found")
	// ErrMalformedAccessToken Refresh token couldn't be parsed.
	ErrMalformedAccessToken = newErr(204, "Access token is malformed")
	// ErrInvalidAccessToken access token not signed on this server.
	ErrInvalidAccessToken = newErr(205, "Invalid access token")
	// ErrAccessTokenExpired access token not signed on this server.
	ErrAccessTokenExpired = newErr(206, "Access token expired")

	// user errors
	// ErrProfileImageTooLarge Error returned when the inputted image file size is too large
	ErrProfileImageTooLarge = newErr(301, "Profile image must be less than 2MB")
	// ErrBackgroundImageTooLarge Error returned when the inputted image file size is too large
	ErrBackgroundImageTooLarge = newErr(302, "Background image must be less than 5MB")
	// ErrEmailAddressInvalid Error returned when the inputted email is invalid or does not match the required pattern
	ErrEmailAddressInvalid = newErr(303, "Invalid email address")
	// ErrFullnameTooShort Error returned when the inputted fullname is too short
	ErrFullnameTooShort = newErr(304, "Fullname is too short")
	// ErrFullnameTooLong Error returned when the inputted fullname is too long
	ErrFullnameTooLong = newErr(305, "Fullname is too long")
	// ErrUsernameTooShort Error returned when the inputted username is too short
	ErrUsernameTooShort = newErr(306, "Username is too short")
	// ErrUsernameTooLong Error returned when the inputted username is too long
	ErrUsernameTooLong = newErr(307, "Username is too long")
	// ErrPasswordTooShort Error returned when the inputted password is too short
	ErrPasswordTooShort = newErr(308, "Password is too short")
	// ErrPasswordTooLong Error returned when the inputted password is too long
	ErrPasswordTooLong = newErr(309, "Password is too long")
	// ErrPasswordInvalid Error returned when the inputted password is invalid or does not match the required pattern
	ErrPasswordInvalid = newErr(310, "Invalid password")
	// ErrEmailAlreadyExist Error returned when the inputted email is already used
	ErrEmailAlreadyExist = newErr(311, "Email already exists")
	// ErrUsernameAlreadyExist Error returned when the inputted username is already used
	ErrUsernameAlreadyExist = newErr(312, "Username already exists")
	// ErrUsernameInvalid Error returned when the inputted username is invalid or does not match the required pattern
	ErrUsernameInvalid = newErr(313, "Invalid username")
	// ErrPasswordIncorrect Error returned when the inputted password is incorrect
	ErrPasswordIncorrect = newErr(314, "Incorrect password")
	// ErrEmptyLogin Error returned when the inputted login is an empty string
	ErrEmptyLogin = newErr(315, "Empty login")
	// ErrEmptyPassword Error returned when the inputted password is an empty string
	ErrEmptyPassword = newErr(316, "Empty password")
)

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (err *Error) Error() string {
	return err.Message
}

func newErr(code int, message string) *Error {
	return &Error{Message: message, Code: code}
}

type MultipleErrors struct {
	Errors []error `json:"errors"`
}

func (multipleErr *MultipleErrors) Error() string {
	messages := make([]string, len(multipleErr.Errors))
	for i, error := range multipleErr.Errors {
		messages[i] = error.Error()
	}
	return strings.Join(messages, ", ")
}
