package main

import (
	_ "bytecrowds-database-server/configuration"

	"bytecrowds-database-server/bytecrowds"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(cors.Default())
	router.GET("/get/:bytecrowd", bytecrowds.GetBytecrowd)
	router.POST("/update", bytecrowds.EditBytecrowd)

	router.GET("/getLanguage/:bytecrowd", bytecrowds.GetLanguage)
	router.POST("/updateLanguage", bytecrowds.EditLanguage)

	router.Run()
}
