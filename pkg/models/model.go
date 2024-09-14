package models

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/sid/FileSystem/pkg/config"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var db gorm.DB

type User struct {
	Username string `gorm:"primaryKey" json:"username"`
	Name     string `json:"name"`
	Age      string `json:"age"`
	Password string `json:"password"`
}

type FileMetadata struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	FileName  string    `gorm:"type:varchar(255)"`
	FileType  string    `gorm:"type:varchar(100)"`
	FileSize  int64     `gorm:"type:bigint"`
	UserID    string    `gorm:"type:varchar(255)"` // Foreign key to the user if you have a user table
	S3Url     string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func init() {
	config.Connect()     // establish connection
	db = *config.GetDB() // initializing the db to a variable

	db.AutoMigrate(&User{})
	db.AutoMigrate(&FileMetadata{})
}

func (u *User) CreateUser() *User {
	hashPass, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	u.Password = string(hashPass)
	db.Create(&u)
	return u
}

func GetUserByUsername(username string) (*User, *gorm.DB) {

	var user User
	db := db.Where("Username=?", username).Find(&user)
	return &user, db
}

func SaveMetadata(fileName, fileType, userID string, fileSize int64, s3Url string) error {
	metadata := FileMetadata{
		FileName: fileName,
		FileType: fileType,
		FileSize: fileSize,
		UserID:   userID,
		S3Url:    s3Url,
	}
	if err := db.Create(&metadata).Error; err != nil {
		return fmt.Errorf("failed to save file metadata: %w", err)
	}

	cacheKey := fmt.Sprintf("file_metadata:%s:%s:%s", userID, fileName, fileType)
	config.RedisClient.Del(context.Background(), cacheKey).Err()

	return nil
}

func GetMetaData(c *gin.Context) ([]FileMetadata, string) {
	claims, _ := c.Get("claims")
	userClaims := claims.(*jwt.MapClaims)
	username := (*userClaims)["username"].(string)

	fileName := c.Query("fileName")
	fileType := c.Query("fileType")

	cacheKey := fmt.Sprintf("file_metadata:%s:%s:%s", username, fileName, fileType)

	var metadata []FileMetadata
	cacheData, err := config.RedisClient.Get(context.Background(), cacheKey).Result()
	if err == redis.Nil {

		var dbMetadata []FileMetadata
		query := db.Model(&FileMetadata{}).Where("user_id = ?", username)

		if fileName != "" {
			query = query.Where("file_name = ?", fileName)
		}
		if fileType != "" {
			query = query.Where("file_type = ?", fileType)
		}

		if err := query.Find(&dbMetadata).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving files"})
			return nil, ""
		}

		if len(dbMetadata) > 0 {
			cacheData, _ := json.Marshal(dbMetadata)
			config.RedisClient.Set(context.Background(), cacheKey, cacheData, 5*time.Minute).Err()
		}
		metadata = dbMetadata
	} else if err == nil {

		_ = json.Unmarshal([]byte(cacheData), &metadata)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching from Redis"})
		return nil, username
	}

	if len(metadata) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No files found with the given metadata"})
		return nil, username
	}

	return metadata, username
}
