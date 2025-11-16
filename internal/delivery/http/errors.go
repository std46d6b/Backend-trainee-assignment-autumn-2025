package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/delivery/http/dto"
	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/domain"
)

func HandleError(c echo.Context, err error) error {
	var domainError *domain.DomainError

	if errors.As(err, &domainError) {
		status := httpStatusCodeMapper(domainError.Code)
		return c.JSON(status, dto.NewErrorResponse(string(domainError.Code), domainError.Message))
	}

	return c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("INTERNAL_SERVER_ERROR", "something went wrong"))
}

func httpStatusCodeMapper(code domain.ErrorCode) int {
	switch code {
	case domain.ErrCodeTeamExists:
		return http.StatusBadRequest
	case domain.ErrCodePRExists:
		return http.StatusConflict
	case domain.ErrCodePRMerged:
		return http.StatusConflict
	case domain.ErrCodeNotAssigned:
		return http.StatusConflict
	case domain.ErrCodeNoCandidate:
		return http.StatusConflict
	case domain.ErrCodeNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
