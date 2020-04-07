package main

import (
	"log"
	"os"

	"Elven9/Lab-Backend/router"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	// 目前先站時設定 Log 到 Standard Output
	log.SetOutput(os.Stdout)

	// 確認是否 K8s-Config 是否存在
	if file, err := os.Open("/root/.kube/config"); err != nil {
		if !os.IsExist(err) {
			log.Panicln("K8s Configuration File is not exist at path '/root/.kube/config', Put it Their and Try again.")
		}
		log.Panicf("Unhandled Error Happened: %v\n", err)
	} else {
		file.Close()
	}
}

func main() {
	engine := gin.Default()

	// CORS Plugin
	engine.Use(cors.Default())

	router.SetUpRouter(engine)

	engine.Run(":8080")
}
