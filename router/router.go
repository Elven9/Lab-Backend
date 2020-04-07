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
	engine.POST("resource/allocation", handler.GetAllocation)

	// Jobs
	engine.GET("job/getJobs", handler.GetJobs)
	engine.GET("job/systemwideStatus", handler.GetSystemwideStatus)
}
