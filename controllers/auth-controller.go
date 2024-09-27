package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sid/FileSystem/models"
	"github.com/sid/FileSystem/utils"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	CreateUser := &models.User{}
	utils.ParseBody(c, CreateUser)

	if serch_user, _ := models.GetUserByUsername(CreateUser.Username); serch_user != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "user already exists",
			"some":  serch_user,
		})
		return
	} else {
		u := CreateUser.CreateUser()
		// res, _ := json.Marshal(u)
		c.JSON(http.StatusOK, gin.H{
			"user": u,
		})
		return
	}
}

func LogIn(c *gin.Context) {
	var userDet = &models.User{}
	utils.ParseBody(c, &userDet)
	ex_username := userDet.Username
	ex_pass := userDet.Password
	user_from_data, _ := models.GetUserByUsername(userDet.Username)
	if user_from_data != nil {
		if err := bcrypt.CompareHashAndPassword([]byte(user_from_data.Password), []byte(ex_pass)); err != nil {
			s, err := utils.CreateToken(ex_username)
			if err != nil {
				c.JSON(http.StatusConflict, gin.H{
					"message": "token not created",
				})
				log.Fatal(err)

				return
			}
			c.JSON(http.StatusOK, gin.H{
				"token": s,
			})
			return
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "incorrect password",
			})
			return
		}
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
	}
}
