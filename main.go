package main

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title   string         `json:"title"`
	Authors pq.StringArray `json:"authors" gorm:"type:text[]"`
}

var (
	db  *gorm.DB
	err error
)

func CreateBook(c *gin.Context) {
	var book Book

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if result := db.Create(&book); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to save in db",
		})
		slog.Error(result.Error.Error())
		return
	}

	c.JSON(http.StatusCreated, book)
}

func GetBook(c *gin.Context) {
	bookId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var book Book
	if result := db.First(&book, bookId); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}

func main() {
	r := gin.Default()

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	r.POST("/book", CreateBook)
	r.GET("/book/:id", GetBook)

	dsn := "host=localhost user=postgres password= dbname=test port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&Book{}); err != nil {
		panic(err)
	}

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
