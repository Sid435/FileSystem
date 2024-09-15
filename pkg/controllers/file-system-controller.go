package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sid/FileSystem/pkg/config"
	"github.com/sid/FileSystem/pkg/models"
	"github.com/sid/FileSystem/pkg/utils"
	"golang.org/x/crypto/bcrypt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var userDetails = &models.User{}
	utils.ParseBody(r, &userDetails)
	username := userDetails.Username
	password := userDetails.Password
	user_data, _ := models.GetUserByUsername(username)
	w.Header().Set("Content-Type", "application/json")
	if err := bcrypt.CompareHashAndPassword([]byte(user_data.Password), []byte(password)); err == nil {
		s, err := utils.CreateToken(username)
		if err != nil {
			w.WriteHeader(http.StatusConflict)
			log.Println(err)
			return
		}
		w.WriteHeader(http.StatusOK)
		res, _ := json.Marshal(s)
		w.Write(res)
		return
	} else {
		s := "Soemthing is missing"
		res, _ := json.Marshal(s)
		w.Write(res)
		return
	}
}

func Signup(w http.ResponseWriter, r *http.Request) {
	CreateUser := &models.User{}

	utils.ParseBody(r, CreateUser)
	u := CreateUser.CreateUser()

	res, _ := json.Marshal(u)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

var bucketName string = "file-system-mangement"
var uploader *s3manager.Uploader

func init() {
	var accessKey string = "AKIAQXHOIJXXON7TKMET"
	var secretkey string = "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"
	var region string = "eu-north-1"

	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			accessKey,
			secretkey,
			"",
		),
	})

	if err != nil {
		panic(err)
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors})
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
func GetPreSignedURL(c *gin.Context) {
	var accessKey string = "AKIAQXHOIJXXON7TKMET"
	var secretkey string = "3OZRIsuV86jxtWyzbhzRXFKQ4OaoqQUTrRD+9MSs"
	var region string = "eu-north-1" // Replace with your AWS Secret Key
	username := c.GetString("username")
	fileName := c.Query("fileName")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File name is required"})
		return
	}

	// Create the S3 key, combining the folder (username) and file name
	s3Key := fmt.Sprintf("%s/%s", username, fileName)

	// Create a new AWS session with credentials
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region), // Set your region
		Credentials: credentials.NewStaticCredentials(
			accessKey, // AWS Access Key ID
			secretkey, // AWS Secret Access Key
			"",        // Session token, if not applicable leave as empty string
		),
	})

	if err != nil {
		log.Println("Failed to create session,", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create AWS session"})
		return
	}

	// Create an S3 client
	svc := s3.New(sess)

	// Create a presigned request for the file
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName), // Ensure bucketName is defined
		Key:    aws.String(s3Key),
	})

	// Generate the presigned URL with a 5-minute expiration
	url, err := req.Presign(5 * time.Minute)
	if err != nil {
		log.Println("Failed to sign request", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate presigned URL"})
		return
	}

	// Return the presigned download URL
	c.JSON(http.StatusOK, gin.H{"download_url": url})
}
func GetFileURL(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	mapClaims, ok := claims.(*jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}
	username, ok := (*mapClaims)["username"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username not found in token"})
		return
	}
	fileName := c.Query("fileName")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fileName is required"})
		return
	}
	cacheKey := fmt.Sprintf("file_url:%s:%s", username, fileName)
	cachedURL, err := config.RedisClient.Get(context.Background(), cacheKey).Result()
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"url": cachedURL})
		return
	}

	_, err = models.GetFileMetadataByName(fileName, username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	s3URL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s/%s", bucketName, username, fileName)

	err = config.RedisClient.Set(context.Background(), cacheKey, s3URL, 5*time.Minute).Err()
	if err != nil {
		log.Printf("Failed to set cache: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"url": s3URL})
}
