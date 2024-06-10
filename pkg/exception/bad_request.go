package exception

import (
	"car-parking-api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type BadRequestError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
	ErrorCode  string `json:"error_code"`
}

func (validationError BadRequestError) Error() string {
	return validationError.Message
}

// NewBadRequestError creates a new instance of BadRequestError with custom status code and error code
func NewBadRequestError(message string, statusCode int, errorCode string) BadRequestError {
	return BadRequestError{
		Message:    message,
		StatusCode: statusCode,
		ErrorCode:  errorCode,
	}
}

func HandleBadRequestError(ctx *fiber.Ctx, message string, statusCode int, errorCode string) error {
	err := NewBadRequestError(message, statusCode, errorCode)

	errorResponse := response.ErrorResponse{
		Error: response.ErrorDetailResponse{
			Code:    err.StatusCode,
			Status:  err.ErrorCode,
			Message: err.Message,
		},
	}

	return ctx.Status(err.StatusCode).JSON(errorResponse)
}
