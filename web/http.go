package web

import (
	"github.com/gin-gonic/gin"
	"github.com/shumybest/ragnaros/config"
	. "github.com/shumybest/ragnaros/logger"
	"github.com/shumybest/ragnaros/security"
	"github.com/toorop/gin-logrus"
	"net/http"
	"sync"
)

type HttpServer struct {
	*gin.Engine
}

var instance *HttpServer
var once sync.Once

func GetHTTPInstance() *HttpServer {
	once.Do(func() {
		instance = &HttpServer{}
		gin.SetMode(gin.ReleaseMode)
		instance.Engine = gin.Default()
	})
	return instance
}

func (h *HttpServer) InitRouter() {
	h.Use(ginlogrus.Logger(Logger))
	h.GET("/", func(c *gin.Context){
		c.String(http.StatusOK, "Welcome to use Ragnaros Spring Cloud golang suites.")
	})
}

func (h *HttpServer) RouterGroup(basePath string) *gin.RouterGroup {
	group := h.Group(basePath)
	return group
}

func (h *HttpServer) SecurityRouterGroup(basePath string) *gin.RouterGroup {
	group := h.Group(basePath)
	group.Use(security.AuthorizationInterceptor)
	return group
}

func (h *HttpServer) Run() {
	_ = h.Engine.Run(":" + config.Context.Port)
}
