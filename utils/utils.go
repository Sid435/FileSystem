package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables")
	}
}

func ParseBody(c *gin.Context, x interface{}) {
	if body, err := ioutil.ReadAll(c.Request.Body); err != nil {
		if err := json.Unmarshal([]byte(body), x); err != nil {
			return
		}
	}
}

func CreateToken(username string) (string, error) {
	jwt_token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().AddDate(0, 0, 30).Unix(), // 30 days validity
		},
	)
	token, err := jwt_token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return token, nil
}
