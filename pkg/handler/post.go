package handler

import (
	"crud/models"
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

// @Summary Create new post
// @Security ApiKeyAuth
// @Tags posts
// @Description create new post
// @Accept json
// @Produce json
// @Param input body models.PostInput true "post data"
// @Success 201 {object} NewObjectResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /posts [post]
func (h *Handler) CreatePost(c echo.Context) error {
	var p models.PostInput

	userID := c.Get("userID").(int)

	if err := c.Bind(&p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := p.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	id, err := h.Services.Post.Create(userID, p)
	if err != nil {
		return newErrorResponse(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, NewObjectResponse{id})
}

// @Summary List posts
// @Tags posts
// @Description get all posts
// @Produce json
// @Produce xml
// @Param format query string false "json or xml produce format"
// @Success 200 {object} []models.Post
// @Failure 500 {object} errorResponse
// @Router /posts [get]
func (h *Handler) GetAllPosts(c echo.Context) error {
	posts, err := h.Services.Post.GetAll()
	if err != nil {
		return newErrorResponse(c, http.StatusInternalServerError, err)
	}

	format := c.QueryParams().Get("format")
	if format == "xml" {
		return c.XML(http.StatusOK, &posts)
	}

	return c.JSON(http.StatusOK, &posts)
}

// @Summary Show signle post
// @Tags posts
// @Description get post by id
// @Produce json
// @Produce xml
// @ID get-post
// @Param id path int true "Post ID"
// @Param format query string false "json or xml produce format"
// @Success 200 {object} models.Post
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /posts:/id [get]
func (h *Handler) GetPostById(c echo.Context) error {
	objID := c.Get("objID").(int)

	post, err := h.Services.Post.GetById(objID)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return newErrorResponse(c, http.StatusInternalServerError, err)
	}

	format := c.QueryParams().Get("format")
	if format == "xml" {
		return c.XML(http.StatusOK, &post)
	}

	return c.JSON(http.StatusOK, &post)
}

// @Summary Update post
// @Security ApiKeyAuth
// @Tags posts
// @Accept  json
// @Description Update post with specified id
// @Param input body models.PostInput true "post data"
// @ID update-post
// @Param id path int true "Post ID"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /posts/:id [patch]
func (h *Handler) UpdatePost(c echo.Context) error {
	objID := c.Get("objID").(int)
	userID := c.Get("userID").(int)

	var data models.PostInput
	if err := c.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if data == (models.PostInput{}) {
		return echo.ErrBadRequest
	}

	err := h.Services.Update(userID, objID, data)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrBadRequest
		}

		return newErrorResponse(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, StatusResponse{"ok"})
}

// @Summary Delete post
// @Security ApiKeyAuth
// @Tags posts
// @Description Delete post with specified id
// @ID delete-post
// @Param id path int true "Post ID"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /posts/:id [delete]
func (h *Handler) DeletePost(c echo.Context) error {
	objID := c.Get("objID").(int)
	userID := c.Get("userID").(int)

	err := h.Services.Post.Delete(userID, objID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrBadRequest
		}

		return newErrorResponse(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, StatusResponse{"ok"})
}
