package exception

import (
	"car-parking-api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type ServerError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
	ErrorCode  string `json:"error_code"`
}

func (validationError ServerError) Error() string {
	return validationError.Message
}

// NewServerError creates a new instance of ServerError with custom status code and error code
func NewServerError(message string, statusCode int, errorCode string) ServerError {
	return ServerError{
		Message:    message,
		StatusCode: statusCode,
		ErrorCode:  errorCode,
	}
}

func HandleServerError(ctx *fiber.Ctx, message string, statusCode int, errorCode string) error {
	err := NewServerError(message, statusCode, errorCode)

	errorResponse := response.ErrorResponse{
		Error: response.ErrorDetailResponse{
			Code:    err.StatusCode,
			Status:  err.ErrorCode,
			Message: err.Message,
		},
	}

	return ctx.Status(err.StatusCode).JSON(errorResponse)
}
