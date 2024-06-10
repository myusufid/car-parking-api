package exception

import (
	"car-parking-api/pkg/logger"
	"car-parking-api/pkg/response"
	"net/http"

	"github.com/cockroachdb/errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	_, ok := err.(ValidationError)
	if ok {
		errorResponse := response.ErrorResponse{
			Error: response.ErrorDetailResponse{
				Code:    400,
				Status:  "BAD_REQUEST",
				Message: err.Error(),
			},
		}
		return ctx.Status(http.StatusBadRequest).JSON(errorResponse)
	}

	if validationErrors, ok := err.(validation.Errors); ok {
		errorsMap := make(map[string]interface{})

		for field, fieldErrors := range validationErrors {
			errorsMap[field] = fieldErrors.Error()
		}

		responseJson := map[string]interface{}{
			"error": map[string]interface{}{
				"code":    422,
				"fields":  getFieldsFromErrors(errorsMap),
				"message": errorsMap,
			},
		}

		return ctx.Status(422).JSON(responseJson)
	}

	var e *fiber.Error
	if errors.As(err, &e) {
		statusMapping := map[int]int{
			fiber.StatusServiceUnavailable: http.StatusServiceUnavailable,
			fiber.StatusUnauthorized:       http.StatusUnauthorized,
			// Add more mappings as needed
		}

		if mappedStatus, exists := statusMapping[e.Code]; exists {
			return ctx.Status(mappedStatus).JSON(response.ErrorResponse{
				Error: response.ErrorDetailResponse{
					Code:    e.Code,
					Status:  http.StatusText(mappedStatus),
					Message: err.Error(),
				},
			})
		}
	}

	logger.WithContext(ctx).Error("error:", zap.Error(errors.WithStack(err)))

	return ctx.Status(http.StatusInternalServerError).JSON(response.ErrorResponse{
		Error: response.ErrorDetailResponse{
			Code:    500,
			Status:  "Internal Server Error",
			Message: err.Error(),
		},
	})
}

func getFieldsFromErrors(errorsMap map[string]interface{}) []string {
	fields := make([]string, 0, len(errorsMap))
	for field := range errorsMap {
		fields = append(fields, field)
	}
	return fields
}
