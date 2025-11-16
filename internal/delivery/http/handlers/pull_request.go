package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	deliveryhttp "github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/delivery/http"
	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/delivery/http/dto"
	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/domain"
)

type PullRequestService interface {
	CreatePullRequest(ctx context.Context, pr domain.PullRequest) (domain.PullRequest, error)
	MergePullRequest(ctx context.Context, prID string) (domain.PullRequest, error)
	ReassignPullRequest(ctx context.Context, prID string, oldReviewerID string) (domain.PullRequest, string, error)
}

func RegisterPullRequestRoutes(e *echo.Echo, s PullRequestService) {
	e.POST("/pullRequest/create", createPullRequestHandler(s))
	e.POST("/pullRequest/merge", mergePullRequestHandler(s))
	e.POST("/pullRequest/reassign", reassignPullRequestHandler(s))
}

// POST /pullRequest/create
func createPullRequestHandler(s PullRequestService) echo.HandlerFunc {
	type requestBody = dto.PullRequestDTO
	type responseBody struct {
		PullRequest dto.PullRequestDTO `json:"pr"`
	}

	return func(c echo.Context) error {
		var req requestBody

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, dto.NewErrorResponse("BAD_REQUEST", "invalid JSON body"))
		}

		domainPullRequest, err := dto.PullRequestDTOToDomain(req)
		if err != nil {
			return deliveryhttp.HandleError(c, err)
		}

		pr, err := s.CreatePullRequest(c.Request().Context(), domainPullRequest)
		if err != nil {
			return deliveryhttp.HandleError(c, err)
		}

		return c.JSON(http.StatusCreated, responseBody{
			PullRequest: dto.PullRequestDomainToDTO(pr),
		})
	}
}

// POST /pullRequest/merge
func mergePullRequestHandler(s PullRequestService) echo.HandlerFunc {
	type requestBody struct {
		PullRequestID string `json:"pull_request_id"`
	}
	type responseBody struct {
		PullRequest dto.PullRequestDTO `json:"pr"`
	}

	return func(c echo.Context) error {
		var req requestBody

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, dto.NewErrorResponse("BAD_REQUEST", "invalid JSON body"))
		}

		if req.PullRequestID == "" {
			return c.JSON(http.StatusBadRequest, dto.NewErrorResponse("BAD_REQUEST", "pull_request_id is required"))
		}

		pr, err := s.MergePullRequest(c.Request().Context(), req.PullRequestID)
		if err != nil {
			return deliveryhttp.HandleError(c, err)
		}

		return c.JSON(http.StatusOK, responseBody{
			PullRequest: dto.PullRequestDomainToDTO(pr),
		})
	}
}

// POST /pullRequest/reassign
func reassignPullRequestHandler(s PullRequestService) echo.HandlerFunc {
	type requestBody struct {
		PullRequestID string `json:"pull_request_id"`
		OldUserId     string `json:"old_user_id"`
	}
	type responseBody struct {
		PullRequest dto.PullRequestDTO `json:"pr"`
		ReplacedBy  string             `json:"replaced_by"`
	}

	return func(c echo.Context) error {
		var req requestBody

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, dto.NewErrorResponse("BAD_REQUEST", "invalid JSON body"))
		}

		if req.PullRequestID == "" || req.OldUserId == "" {
			return c.JSON(http.StatusBadRequest, dto.NewErrorResponse("BAD_REQUEST", "pull_request_id and old_user_id are required"))
		}

		pr, replacedBy, err := s.ReassignPullRequest(c.Request().Context(), req.PullRequestID, req.OldUserId)
		if err != nil {
			return deliveryhttp.HandleError(c, err)
		}

		return c.JSON(http.StatusOK, responseBody{
			PullRequest: dto.PullRequestDomainToDTO(pr),
			ReplacedBy:  replacedBy,
		})
	}
}
