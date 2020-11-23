package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sgkul2000/go-rest-api/controller"
	"github.com/sgkul2000/go-rest-api/entity"
	"github.com/sgkul2000/go-rest-api/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	videoService    service.VideoService       = service.New()
	videoController controller.VideoController = controller.New(videoService)
)

func main() {
	server := gin.Default()
	server.GET("/", handleFindVideos)
	server.POST("/", handleCreateVideo)
	server.PUT("/:id", handleUpdateVideo)
	server.DELETE("/:id", handleDeleteVideo)
	server.Run(":8080")
	defer entity.CloseDB()
}

func handleCreateVideo(c *gin.Context) {
	var video entity.Video
	if err := c.ShouldBindJSON(&video); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}
	id, err := entity.Create(&video)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func handleFindVideos(c *gin.Context) {
	// var videos []entity.Video
	searchTitle, _ := c.GetQuery("search")
	videos, err := entity.Find(searchTitle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, videos)
}

func handleUpdateVideo(c *gin.Context) {
	var video entity.Video
	id := c.Param("id")
	if err := c.ShouldBindJSON(&video); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}
	err := entity.Update(&id, &video)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func handleDeleteVideo(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
		})
		return
	}
	status, err := entity.Delete(id)

	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
		})
		return
	}
	c.JSON(200, gin.H{
		"success": status,
	})
}
