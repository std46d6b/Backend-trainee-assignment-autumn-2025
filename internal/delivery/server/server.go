package server

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/delivery/http/handlers"
)

type Server struct {
	echo *echo.Echo
}

func NewServer(
	teamService handlers.TeamService,
	userService handlers.UserService,
	pullRequestService handlers.PullRequestService,
) *Server {
	e := echo.New()

	api := e.Group("")

	handlers.RegisterTeamRoutes(api, teamService)
	handlers.RegisterUserRoutes(e, userService)
	handlers.RegisterPullRequestRoutes(e, pullRequestService)

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	return &Server{
		echo: e,
	}
}

func (s *Server) Start(addr string) error {
	return s.echo.Start(addr)
}

func (s *Server) Stop(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}
