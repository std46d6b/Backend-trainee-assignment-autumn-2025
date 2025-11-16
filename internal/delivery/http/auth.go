package http

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/delivery/http/dto"
)

const (
	adminHeader = "X-Admin-Token"
	userHeader  = "X-User-Token"
)

func AdminOnlyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	adminToken := os.Getenv("ADMIN_TOKEN")

	return func(c echo.Context) error {
		token := c.Request().Header.Get(adminHeader)
		if token == "" || token != adminToken {
			return c.JSON(http.StatusUnauthorized,
				dto.NewErrorResponse("BAD_REQUEST", "invalid admin token"),
			)
		}

		return next(c)
	}
}

func AdminOrUserMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	adminToken := os.Getenv("ADMIN_TOKEN")
	userToken := os.Getenv("USER_TOKEN")

	return func(c echo.Context) error {
		a := c.Request().Header.Get(adminHeader)
		u := c.Request().Header.Get(userHeader)

		ok := (adminToken != "" && a == adminToken) || (userToken != "" && u == userToken)
		if !ok {
			return c.JSON(http.StatusUnauthorized,
				dto.NewErrorResponse("BAD_REQUEST", "invalid token"),
			)
		}

		return next(c)
	}
}
