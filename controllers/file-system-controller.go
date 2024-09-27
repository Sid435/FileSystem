package controllers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sid/FileSystem/config"
	"github.com/sid/FileSystem/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	bucketName string
	uploader   *s3manager.Uploader
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables")
	}
	bucketName = os.Getenv("AWS_BUCKET_NAME")
	if bucketName == "" {
		log.Fatal("AWS_BUCKET_NAME environment variable is not set")
	}

	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")

	if accessKey == "" || secretKey == "" || region == "" {
		log.Fatal("AWS credentials or region are not set in environment variables")
	}

	awsSession, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})

	if err != nil {
		log.Fatalf("Failed to create AWS session: %v", err)
	}

	uploader = s3manager.NewUploader(awsSession)
}

func UploadFile(c *gin.Context) {
	username := c.GetString("username")
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var errors []string
	var uploadedURLs []string
	files := form.File["files"]
	for _, file := range files {
		fileHeader := file

		f, err := fileHeader.Open()
		if err != nil {
			errors = append(errors, fmt.Sprintf("Error opening file %s: %s", fileHeader.Filename, err.Error()))
			continue
		}
		defer f.Close()

		uploadedURL, err := saveFile(f, fileHeader, username)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Error saving file %s: %s", fileHeader.Filename, err.Error()))
		} else {
			uploadedURLs = append(uploadedURLs, uploadedURL)
			err := models.SaveMetadata(fileHeader.Filename, fileHeader.Header.Get("Content-Type"), username, fileHeader.Size, uploadedURL)
			if err != nil {
				errors = append(errors, fmt.Sprintf("Error saving metadata for file %s: %s", fileHeader.Filename, err.Error()))
			}
		}
	}

	if len(errors) > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": errors})
	} else {
		c.JSON(http.StatusOK, gin.H{"urls": uploadedURLs})
	}
}

func saveFile(fileReader io.Reader, fileHeader *multipart.FileHeader, username string) (string, error) {
	s3Key := fmt.Sprintf("%s/%s", username, fileHeader.Filename)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(s3Key),
		Body:   fileReader,
	})
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, s3Key)

	return url, nil
}

func GetSignedUrl(c *gin.Context) {
	username := c.GetString("username")
	filename := c.Query("file_name")

	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "file is missing",
		})
	}

	s3Key := fmt.Sprintf("%s/%s", username, filename)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
	})

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "couldn't create session",
		})
	}

	svc := s3.New(sess)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
		Key:    aws.String(s3Key),
	})

	url, err := req.Presign(5 * time.Minute)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "couldn't upload file",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"url": url,
	})
}

func DeleteFile(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	fileName := c.Query("file_name")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fileName is required"})
		return
	}

	metadata, err := models.GetFileMetadataByName(fileName, username.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	session, _ := session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
	})
	svc := s3.New(session)

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fmt.Sprintf("%s/%s", username, fileName)),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file from S3"})
		return
	}

	err = models.DeleteFileMetadata(metadata.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file metadata"})
		return
	}

	cacheKey := fmt.Sprintf("file_url:%s:%s", username, fileName)
	err = config.RedisClient.Del(context.Background(), cacheKey).Err()
	if err != nil {
		log.Printf("Failed to clear cache: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}
