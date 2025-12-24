package handler

import (
	"context"
	"net/http"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/caovanhoang63/cloud-traveler/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadHandler struct {
	s3Client   *s3.Client
	bucketName string
	fileRepo   *storage.FileRepository
}

func NewUploadHandler(s3Client *s3.Client, bucketName string, fileRepo *storage.FileRepository) *UploadHandler {
	return &UploadHandler{
		s3Client:   s3Client,
		bucketName: bucketName,
		fileRepo:   fileRepo,
	}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "No file provided: " + err.Error(),
		})
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	key := uuid.New().String() + ext
	contentType := header.Header.Get("Content-Type")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	_, err = h.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(h.bucketName),
		Key:           aws.String(key),
		Body:          file,
		ContentLength: aws.Int64(header.Size),
		ContentType:   aws.String(contentType),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Failed to upload to S3: " + err.Error(),
		})
		return
	}

	uploadedFile := &storage.UploadedFile{
		Filename:    header.Filename,
		S3Key:       key,
		ContentType: contentType,
		SizeBytes:   header.Size,
	}

	if err := h.fileRepo.Create(ctx, uploadedFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "File uploaded to S3 but failed to save metadata: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "File uploaded successfully",
		"data":    uploadedFile,
	})
}
