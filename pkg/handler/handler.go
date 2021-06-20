package handler

import (
	"crud/pkg/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"
	"strconv"

	_ "crud/docs"
)

type Handler struct {
	Services *service.Service
	//Logger log.Logger
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{s}
}

func (h *Handler) getUserID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return echo.ErrUnauthorized
		}
		userID, err := h.Services.User.ParseToken(token)
		if err != nil {
			return echo.ErrUnauthorized
		}
		c.Set("userID", userID)
		return next(c)
	}
}

func (h *Handler) getObjectID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id < 1 {
			return echo.ErrBadRequest
		}

		c.Set("objID", id)
		return next(c)
	}
}

func (h *Handler) InitRoutes(e *echo.Echo) {

	jwtMiddleware := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(service.SigningKey),
	})

	apiV1 := e.Group("/api/v1")

	posts := apiV1.Group("/posts")

	{
		posts.GET("", h.GetAllPosts)
		posts.GET("/:id", h.GetPostById, h.getObjectID)
		posts.POST("", h.CreatePost, jwtMiddleware, h.getUserID)
		posts.DELETE("/:id", h.DeletePost, jwtMiddleware, h.getUserID, h.getObjectID)
		posts.PATCH("/:id", h.UpdatePost, jwtMiddleware, h.getUserID, h.getObjectID)
	}

	auth := e.Group("/auth")
	{
		auth.POST("/sign-up", h.SignUp)
		auth.POST("/sign-in", h.SignIn)

		auth.GET("/google", h.GetGoogleAuthUrl)
		auth.GET("/callback/google", h.GoogleAuthCallback)
	}

	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
