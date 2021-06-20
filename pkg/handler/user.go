package handler

import (
	"crud/models"
	"crud/pkg/repository"
	"crud/pkg/service"
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

// @Summary Sign Up
// @Tags auth
// @Description sign up (create new user)
// @Accept json
// @Produce json
// @Param input body models.UserInput true "user data"
// @Success 201 {object} StatusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /sign-up [post]
func (h *Handler) SignUp(c echo.Context) error {
	var input models.UserInput

	if err := c.Bind(&input); err != nil {
		return echo.ErrBadRequest
	}

	if err := input.Validate(); err != nil {
		return echo.ErrBadRequest
	}

	if err := h.Services.User.SignUp(input); err != nil {
		if _, err := repository.IsDuplicateKeyError(err); err == nil {
			return echo.NewHTTPError(http.StatusBadRequest, "login is already taken")
		}
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, StatusResponse{"ok"})
}

// @Summary Sign In
// @Tags auth
// @Description return jwt token
// @Accept json
// @Produce json
// @Param input body models.UserInput true "user data"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /sign-in [post]
func (h *Handler) SignIn(c echo.Context) error {
	var input models.UserInput
	if err := c.Bind(&input); err != nil {
		return echo.ErrInternalServerError
	}

	if err := input.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.Services.User.SignIn(input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "wrong user credentials")
		}
		return echo.ErrInternalServerError
	}

	JWTToken, err := h.Services.User.GenerateToken(user)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, TokenResponse{JWTToken})
}

func (h *Handler) GetGoogleAuthUrl(c echo.Context) error {
	url := h.Services.User.GetGoogleAuthUrl()
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) GoogleAuthCallback(c echo.Context) error {
	state := c.FormValue("state")
	if state != service.AuthState {
		return echo.ErrBadRequest
	}

	accessToken, err := h.Services.User.GoogleAuthExchangeCodeForToken(c.Request().Context(), c.FormValue("code"))
	if err != nil {
		return echo.ErrBadGateway
	}

	email, err := h.Services.User.GetEmailFromGoogleAccessToken(c.Request().Context(), accessToken)
	if err != nil {
		return echo.ErrBadGateway
	}

	user, err := h.Services.User.GetOrCreateUserByEmail(email)
	if err != nil {
		return echo.ErrInternalServerError
	}

	JWTToken, err := h.Services.User.GenerateToken(user)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, TokenResponse{JWTToken})
}
