package routes

import (
	"github.com/gin-gonic/gin"
	"helloworld-api/controller"
)

func Route(r *gin.Engine) {
	r.GET("/api/data/search/by-data", controller.SearchDataByData)
	r.GET("/api/data/search/by-time", controller.SearchDataByTime)
	r.GET("/api/data/count/by-time", controller.CountDataByTime)
}
