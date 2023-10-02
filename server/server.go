package server

import (
	"jwt/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Server struct {
	DB *gorm.DB
	R  *gin.Engine
}

func NewServer() *Server {
	db, err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.User{})

	r := gin.Default()

	return &Server{
		DB: db,
		R:  r,
	}
}
