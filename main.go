package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"helloworld-api/database"
	"helloworld-api/routes"
	"helloworld-api/tables"
	"net/http"

	_ "helloworld-api/docs"
)

func main() {
	rootCtx := context.Background()
	r := gin.New()
	db := database.ConnectDatabase()
	tables.Initial(db, rootCtx)
	routes.Route(r)
	r.GET("/books", func(c *gin.Context) {
		c.JSON(http.StatusOK, render.JSON{Data: "123123"})
	})
	//Crawl()
	r.Run(":8080")

}
