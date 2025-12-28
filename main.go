package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/caovanhoang63/cloud-traveler/internal/config"
	"github.com/caovanhoang63/cloud-traveler/internal/handler"
	"github.com/caovanhoang63/cloud-traveler/internal/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db := storage.NewPostgresDB(storage.DBConfig{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
	})
	defer db.Close()

	if err := storage.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	s3Client, err := storage.NewS3Client(cfg.AWSRegion)
	if err != nil {
		log.Fatalf("Failed to create S3 client: %v", err)
	}

	fileRepo := storage.NewFileRepository(db)

	healthHandler := handler.NewHealthHandler(db)
	uploadHandler := handler.NewUploadHandler(s3Client, cfg.S3BucketName, fileRepo)
	fileHandler := handler.NewFileHandler(fileRepo)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.MaxMultipartMemory = 32 << 20

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	r.SetHTMLTemplate(tmpl)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	api := r.Group("/api")
	{
		api.GET("/health/db", healthHandler.CheckDB)
		api.POST("/upload", uploadHandler.Upload)
		api.GET("/files", fileHandler.List)
		api.GET("/files/:id", fileHandler.Get)
		api.DELETE("/files/:id", fileHandler.Delete)
	}

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
