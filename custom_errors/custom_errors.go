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
