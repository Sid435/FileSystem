package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {

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

}
