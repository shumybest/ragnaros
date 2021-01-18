package actuator

import (
	"github.com/gin-gonic/gin"
)

func Controller(router *gin.RouterGroup) {
	router.GET("/health", healthHandler)
	router.GET("/info", infoHandler)
	router.GET("/env", envHandler)
	router.GET("/configprops", configHandler)
}
