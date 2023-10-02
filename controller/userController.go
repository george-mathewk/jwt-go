package controller

import (
	"jwt/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SignUp(c *gin.Context) {
	db, err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.User{})
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error is getting data",
		})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error in hashing password",
		})
		return
	}

	user.Password = string(hash)

	resp := db.Create(&user)
	if resp.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error is creating user",
		})
		return
	}

	c.JSON(http.StatusOK, "User succesfully created")
}

func Login(c *gin.Context) {
	db, err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.User{})

	var user models.User

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error is getting data",
		})
		return
	}

	var logger models.User

	resp := db.Find(&logger, "email = ?", user.Email)

	if resp.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error in finding email",
		})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(logger.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Wrong password",
		})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": logger.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte("pigpiugpiugpiugpuigp"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error in signing JWT",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{})
}

func Validate(c *gin.Context){
	user, found := c.Get("user")
	if !found{
		c.JSON(http.StatusBadRequest, gin.H{
			"error":"Error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":user,
	})
}
