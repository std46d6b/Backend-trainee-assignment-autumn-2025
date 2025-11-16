package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	deliveryhttp "github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/delivery/http"
	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/delivery/http/dto"
	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/domain"
)

type TeamService interface {
	CreateTeam(ctx context.Context, up domain.TeamUpsert) (domain.TeamUpsert, error)
	GetTeamWithMembers(ctx context.Context, teamName string) (domain.TeamUpsert, error)
}

func RegisterTeamRoutes(e *echo.Group, s TeamService) {
	e.POST("/team/add", createTeamHandler(s))
	e.GET("/team/get", getTeamHandler(s))
}

// createTeamHandler handles POST /team/add.
func createTeamHandler(s TeamService) echo.HandlerFunc {
	type requestBody = dto.TeamDTO
	type responseBody struct {
		Team dto.TeamDTO `json:"team"`
	}

	return func(c echo.Context) error {
		var req requestBody

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, dto.NewErrorResponse("BAD_REQUEST", "invalid JSON body"))
		}

		upsertTeam := dto.TeamDTOToDomain(req)

		team, err := s.CreateTeam(c.Request().Context(), upsertTeam)
		if err != nil {
			return deliveryhttp.HandleError(c, err)
		}

		return c.JSON(http.StatusCreated, responseBody{
			Team: dto.TeamDomainToDTO(team),
		})
	}
}

// getTeamHandler handles GET /team/get.
func getTeamHandler(s TeamService) echo.HandlerFunc {
	return func(c echo.Context) error {
		teamName := c.QueryParam("team_name")

		if teamName == "" {
			return c.JSON(http.StatusBadRequest, dto.NewErrorResponse("BAD_REQUEST", "team_name is required"))
		}

		team, err := s.GetTeamWithMembers(c.Request().Context(), teamName)
		if err != nil {
			return deliveryhttp.HandleError(c, err)
		}

		return c.JSON(http.StatusOK, dto.TeamDomainToDTO(team))
	}
}
