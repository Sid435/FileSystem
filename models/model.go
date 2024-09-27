package models

import (
	"context"
	"fmt"
	"time"

	"github.com/sid/FileSystem/config"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Username string `gorm:"primaryKey" json:"username"`
	Name     string `json:"name"`
	Age      string `json:"age"`
	Password string `json:"password"`
}

var db gorm.DB

func init() {
	db = *config.Connect()
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

type FileMetadata struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	FileName  string    `gorm:"type:varchar(255)"`
	FileType  string    `gorm:"type:varchar(100)"`
	FileSize  int64     `gorm:"type:bigint"`
	UserID    string    `gorm:"type:varchar(255)"`
	S3Url     string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
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

func GetFileMetadataByName(fileName, userID string) (*FileMetadata, error) {
	var metadata FileMetadata
	result := db.Where("file_name = ? AND user_id = ?", fileName, userID).First(&metadata)
	if result.Error != nil {
		return nil, result.Error
	}
	return &metadata, nil
}

func DeleteFileMetadata(id uint) error {
	result := db.Delete(&FileMetadata{}, id)
	return result.Error
}
