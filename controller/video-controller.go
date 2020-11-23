package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sgkul2000/go-rest-api/entity"
	"github.com/sgkul2000/go-rest-api/service"
)

// VideoController is an interface for controller struct
type VideoController interface {
	FindAll() []entity.Video
	Save(ctx *gin.Context) entity.Video
}

type controller struct {
	service service.VideoService
}

// New is a constructor for struct controller
func New(service service.VideoService) VideoController {
	return &controller{
		service: service,
	}
}

func (c *controller) FindAll() []entity.Video {
	return c.service.FindAll()
}
func (c *controller) Save(ctx *gin.Context) entity.Video {
	var video entity.Video
	ctx.BindJSON(&video)
	c.service.Save(video)
	return video
}
