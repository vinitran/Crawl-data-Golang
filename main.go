package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"helloworld-api/database"
	_ "helloworld-api/docs"
	"helloworld-api/routes"
	"helloworld-api/tables"
)

func main() {
	rootCtx := context.Background()
	r := gin.New()
	db := database.ConnectDatabase()
	tables.Initial(db, rootCtx)
	routes.Route(r)
	//Crawl()
	r.Run(":8080")

}
