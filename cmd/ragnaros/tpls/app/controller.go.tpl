package app

import (
	"github.com/shumybest/ragnaros"
	"github.com/gin-gonic/gin"
	"net/http"
)

func {{ Export .App.ControllerName }}(r *ragnaros.Context) {
	demo := r.RouterGroup("/{{ .App.ControllerName }}")

	demo.GET("/helloworld", func(c *gin.Context) {
		var response = map[string]string {
			"hello": "world",
		}
		c.JSON(http.StatusOK, response)
	})
}

