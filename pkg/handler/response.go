package handler

import (
	"github.com/labstack/echo/v4"
)

type StatusResponse struct {
	Status string `json:"status"`
}

type NewObjectResponse struct {
	ID int `json:"id"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c echo.Context, statusCode int, err error, message ...string) error {
	c.Logger().Error(err)
	if len(message) == 0 {
		return echo.NewHTTPError(statusCode)
	}
	return echo.NewHTTPError(statusCode, errorResponse{message[0]})
}
