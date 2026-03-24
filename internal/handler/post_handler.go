package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"body-smith-be/internal/model"
	"body-smith-be/internal/service"
)

type PostHandler struct {
	postService service.PostService
}

func NewPostHandler(postService service.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var req model.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	post, err := h.postService.Create(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidPostInput) {
			respondError(c, http.StatusBadRequest, err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, "failed to create post")
		return
	}

	c.JSON(http.StatusCreated, post)
}

func (h *PostHandler) ListAdminPosts(c *gin.Context) {
	page, perPage := parsePagination(c)

	response, err := h.postService.ListAdmin(c.Request.Context(), page, perPage)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "failed to list posts")
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *PostHandler) GetAdminPost(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid post id")
		return
	}

	post, err := h.postService.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrPostNotFound) {
			respondError(c, http.StatusNotFound, err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, "failed to get post")
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid post id")
		return
	}

	var req model.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	post, err := h.postService.Update(c.Request.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidPostInput):
			respondError(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrPostNotFound):
			respondError(c, http.StatusNotFound, err.Error())
		default:
			respondError(c, http.StatusInternalServerError, "failed to update post")
		}
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid post id")
		return
	}

	if err := h.postService.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, service.ErrPostNotFound) {
			respondError(c, http.StatusNotFound, err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, "failed to delete post")
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *PostHandler) ListPublicPosts(c *gin.Context) {
	page, perPage := parsePagination(c)

	response, err := h.postService.ListPublic(c.Request.Context(), page, perPage)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "failed to list posts")
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *PostHandler) GetPublicPost(c *gin.Context) {
	post, err := h.postService.GetBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		if errors.Is(err, service.ErrPostNotFound) {
			respondError(c, http.StatusNotFound, err.Error())
			return
		}
		respondError(c, http.StatusInternalServerError, "failed to get post")
		return
	}

	c.JSON(http.StatusOK, post)
}

func parsePagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	return page, perPage
}
