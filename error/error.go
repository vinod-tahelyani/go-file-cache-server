package error

import "fmt"

type HTTPError struct {
	ErrorMessage string `json:"errorMessage"`
}

func GetHTTPError(msg string, originalError error) HTTPError {
	errorMessage := fmt.Sprintf("%s: %v", msg, originalError)
	return HTTPError{ErrorMessage: errorMessage}
}

func WrapError(msg string, originalError error) string {
	return fmt.Sprintf("%s: %v", msg, originalError)
}