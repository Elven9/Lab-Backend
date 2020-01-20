package main

import (
	"log"
	"os"

	"Elven9/Lab-Backend/router"

	"github.com/gin-gonic/gin"
)

func init() {
	// 目前先站時設定 Log 到 Standard Output
	log.SetOutput(os.Stdout)
}

func main() {
	engine := gin.Default()

	router.SetUpRouter(engine)

	engine.Run(":8080")
}
