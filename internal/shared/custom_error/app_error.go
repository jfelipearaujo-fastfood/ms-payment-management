package custom_error

import "github.com/labstack/echo/v4"

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

func NewHttpAppError(code int, message string, err error) *echo.HTTPError {
	appError := AppError{
		Code:    code,
		Message: message,
		Details: err.Error(),
	}

	return echo.NewHTTPError(code, appError)
}

func NewHttpAppErrorFromBusinessError(err error) *echo.HTTPError {
	buErr := err.(BusinessError)
	return NewHttpAppError(buErr.Code(), buErr.Title(), err)
}
