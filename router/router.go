package router

import (
	"github.com/gin-gonic/gin"

	"github.com/Elven9/Lab-Backend/router/handler"
)

// SetUpRouter ,初始化 Router
func SetUpRouter(engine *gin.Engine) {
	// System Info Related Route.
	engine.GET("system/hardwareSpec", handler.GetSystemInfo)

	// Queue Statics
	engine.GET("queue/statistic", handler.GetQueueStatics)
	engine.GET("queue/interactive", handler.GetInteractiveInfo)
	engine.GET("queue/train", handler.GetTrainInfo)
	engine.GET("queue/service", handler.GetServiceInfo)
}
