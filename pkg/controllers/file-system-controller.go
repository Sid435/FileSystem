package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/sid/FileSystem/pkg/models"
	"github.com/sid/FileSystem/pkg/utils"
	"golang.org/x/crypto/bcrypt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
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
	// Define the S3 key with the username as the folder
	s3Key := fmt.Sprintf("%s/%s", username, fileHeader.Filename)

	// Upload the file to S3 using the fileReader
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(s3Key),
		Body:   fileReader,
	})
	if err != nil {
		return "", err
	}

	// Get the URL of the uploaded file
	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, s3Key)

	return url, nil
}

func GetFileURL(c *gin.Context) {
	metadata, username := models.GetMetaData(c)
	if metadata != nil {
		var fileURLs []string
		for _, meta := range metadata {
			// Prefix the S3 path with the username folder
			s3Path := fmt.Sprintf("https://%s.s3.amazonaws.com/%s/%s", bucketName, username, meta.FileName)
			fileURLs = append(fileURLs, s3Path)
		}

		c.JSON(http.StatusOK, gin.H{"urls": fileURLs})
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "No files found with the given metadata"})
}
