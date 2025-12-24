package handler

import (
	"net/http"
	"strconv"

	"github.com/caovanhoang63/cloud-traveler/internal/storage"
	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	repo *storage.FileRepository
}

func NewFileHandler(repo *storage.FileRepository) *FileHandler {
	return &FileHandler{repo: repo}
}

func (h *FileHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	files, err := h.repo.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Failed to list files: " + err.Error(),
		})
		return
	}

	count, _ := h.repo.Count(c.Request.Context())

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"files":  files,
			"total":  count,
			"limit":  limit,
			"offset": offset,
		},
	})
}

func (h *FileHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Invalid file ID",
		})
		return
	}

	file, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "File not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   file,
	})
}

func (h *FileHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Invalid file ID",
		})
		return
	}

	s3Key, err := h.repo.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "File not found or already deleted",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "File record deleted",
		"s3_key":  s3Key,
	})
}
