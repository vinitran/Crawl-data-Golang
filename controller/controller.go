package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/uptrace/bun"
	"helloworld-api/database"
	"helloworld-api/tables"
	"net/http"
)

var db *bun.DB

func init() {
	db = database.ConnectDatabase()
}

type SearchByDataRequest struct {
	TypeData string
	Data     string
}

type SearchByTimeRequest struct {
	Time     string
	TypeData string
}

func SearchDataByData(c *gin.Context) {
	responseData := new(tables.Data)

	var requestBody SearchByDataRequest
	requestBody.Data = c.Query("data")
	requestBody.TypeData = c.Query("type")
	//
	//if err := c.BindJSON(&requestBody); err != nil {
	//	fmt.Println("err", err.Error())
	//	return
	//}

	if requestBody.TypeData == "" {
		err := db.NewSelect().Model(responseData).
			Where("data = ?", requestBody.Data).
			Scan(context.Background())
		if err != nil {
			fmt.Println("err", err.Error())
			return
		}

		c.JSON(http.StatusOK, render.JSON{Data: responseData})
		return
	}

	err := db.NewSelect().Model(responseData).
		Where("type = ?", requestBody.TypeData).
		Where("data = ?", requestBody.Data).
		Scan(context.Background())
	if err != nil {
		fmt.Println("err", err.Error())
		return
	}

	c.JSON(http.StatusOK, render.JSON{Data: responseData})
	return
}

func SearchDataByTime(c *gin.Context) {
	responseData := new([]tables.Data)

	var requestBody SearchByTimeRequest
	requestBody.Time = c.Query("time")
	requestBody.TypeData = c.Query("type")
	//if err := c.BindJSON(&requestBody); err != nil {
	//	fmt.Println("err", err.Error())
	//	return
	//}

	if requestBody.TypeData == "" {
		err := db.NewSelect().Model(responseData).
			Where("time = ?", requestBody.Time).
			Scan(context.Background())
		if err != nil {
			fmt.Println("err", err.Error())
			return
		}

		c.JSON(http.StatusOK, render.JSON{Data: responseData})
		return
	}

	err := db.NewSelect().Model(responseData).
		Where("time = ?", requestBody.Time).
		Where("type = ?", requestBody.TypeData).
		Scan(context.Background())
	if err != nil {
		fmt.Println("err", err.Error())
		return
	}

	c.JSON(http.StatusOK, render.JSON{Data: responseData})
}

func CountDataByTime(c *gin.Context) {
	responseData := new(tables.Data)

	var requestBody SearchByTimeRequest
	requestBody.Time = c.Query("time")
	requestBody.TypeData = c.Query("type")
	//if err := c.BindJSON(&requestBody); err != nil {
	//	fmt.Println("err", err.Error())
	//	return
	//}

	if requestBody.TypeData == "" {
		value, err := db.NewSelect().Model(responseData).
			Where("time = ?", requestBody.Time).
			Count(context.Background())
		if err != nil {
			fmt.Println("err", err.Error())
			return
		}

		c.JSON(http.StatusOK, render.JSON{Data: value})
		return
	}

	value, err := db.NewSelect().Model(responseData).
		Where("time = ?", requestBody.Time).
		Where("type = ?", requestBody.TypeData).
		Count(context.Background())
	if err != nil {
		fmt.Println("err", err.Error())
		return
	}

	c.JSON(http.StatusOK, render.JSON{Data: value})
	return
}
