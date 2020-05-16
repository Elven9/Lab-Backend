package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"Elven9/Lab-Backend/router"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Program Argument Definition
var toEscapeCheck bool
var port int

func init() {
	// 目前先站時設定 Log 到 Standard Output
	log.SetOutput(os.Stdout)

	// Get Program Argument
	flag.BoolVar(&toEscapeCheck, "escapeCheck", false, "Set if Want to Escape Preflight Check.")
	flag.IntVar(&port, "p", 8080, "Server Bind Port")

	flag.Parse()

	if !toEscapeCheck {
		// 確認是否 K8s-Config 是否存在
		if file, err := os.Open("/root/.kube/config"); err != nil {
			if !os.IsExist(err) {
				log.Panicln("K8s Configuration File is not exist at path '/root/.kube/config', Put it Their and Try again.")
			}
			log.Panicf("Unhandled Error Happened: %v\n", err)
		} else {
			file.Close()
		}
	} else {
		log.Println("Preflight Check Escaped.")
	}
}

func main() {
	engine := gin.Default()

	// CORS Plugin
	engine.Use(cors.Default())

	router.SetUpRouter(engine)

	engine.Run(fmt.Sprintf(":%d", port))
}
