package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	deliveryhttp "github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/delivery/http"
	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/delivery/http/dto"
	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/domain"
)

type UserService interface {
	SetIsActive(ctx context.Context, userID string, isActive bool) (domain.User, error)
	ListReviewPRs(ctx context.Context, userID string) ([]domain.PullRequest, error)
}

func RegisterUserRoutes(e *echo.Echo, s UserService) {
	e.POST("/users/setIsActive", deliveryhttp.AdminOnlyMiddleware(setIsActiveHandler(s)))
	e.GET("/users/getReview", deliveryhttp.AdminOrUserMiddleware(getReviewHandler(s)))
}

// setIsActiveHandler handles POST /users/setIsActive.
func setIsActiveHandler(s UserService) echo.HandlerFunc {
	type requestBody struct {
		UserID   string `json:"user_id"`
		IsActive bool   `json:"is_active"`
	}
	type responseBody struct {
		User dto.UserDTO `json:"user"`
	}

	return func(c echo.Context) error {
		var req requestBody

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, dto.NewErrorResponse("BAD_REQUEST", "invalid JSON body"))
		}

		if req.UserID == "" {
			return c.JSON(http.StatusBadRequest, dto.NewErrorResponse("BAD_REQUEST", "user_id is required"))
		}

		user, err := s.SetIsActive(c.Request().Context(), req.UserID, req.IsActive)
		if err != nil {
			return deliveryhttp.HandleError(c, err)
		}

		return c.JSON(http.StatusOK, responseBody{
			User: dto.UserDomainToDTO(user),
		})
	}
}

// getReviewHandler handles GET /users/getReview.
func getReviewHandler(s UserService) echo.HandlerFunc {
	type responseBody struct {
		UserID       string                    `json:"user_id"`
		PullRequests []dto.PullRequestShortDTO `json:"pull_requests"`
	}

	return func(c echo.Context) error {
		userID := c.QueryParam("user_id")

		if userID == "" {
			return c.JSON(http.StatusBadRequest, dto.NewErrorResponse("BAD_REQUEST", "user_id is required"))
		}

		pullRequests, err := s.ListReviewPRs(c.Request().Context(), userID)

		if err != nil {
			return deliveryhttp.HandleError(c, err)
		}

		var pullRequestShorts []dto.PullRequestShortDTO
		for _, pullRequest := range pullRequests {
			pullRequestShorts = append(pullRequestShorts, dto.PullRequestShortDomainToDTO(pullRequest))
		}

		return c.JSON(http.StatusOK, responseBody{
			UserID:       userID,
			PullRequests: pullRequestShorts,
		})
	}
}
