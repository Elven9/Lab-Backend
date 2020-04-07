package router

import (
	"github.com/gin-gonic/gin"

	"Elven9/Lab-Backend/router/handler"
)

// SetUpRouter ,初始化 Router
func SetUpRouter(engine *gin.Engine) {
	// System Info Related Route.
	engine.GET("system/hardwareSpec", handler.GetSystemInfo)

	// Allocation
	engine.GET("resource/allocation", handler.GetAllocation)

	// Utilization
	engine.GET("resource/utilization", handler.GetUtilization)

	// Jobs
	engine.GET("job/getJobs", handler.GetJobs)
	engine.GET("job/systemwideStatus", handler.GetSystemwideStatus)
}
