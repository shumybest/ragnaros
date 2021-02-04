package ragnaros

import (
	"github.com/shumybest/ragnaros/actuator"
	"github.com/shumybest/ragnaros/config"
	"github.com/shumybest/ragnaros/eureka"
	"github.com/shumybest/ragnaros/logger"
	"github.com/shumybest/ragnaros/repository"
	"github.com/shumybest/ragnaros/web"
	"runtime/debug"
)

type injectedApp func(*Context)
var apps []injectedApp

func init() {
	config.InitConfig()
	logger.InitLogger()
}

func InjectApps(injectedFunc ...injectedApp) {
	apps = append(apps, injectedFunc...)
}

func Start() {
	initComponents()
	run()
}

func initComponents() {
	e := eureka.GetClientInstance()
	e.Register()

	m := repository.GetMySQLInstance()
	m.InitConnection()

	r := repository.GetRedisInstance()
	r.InitConnection()

	h := web.GetHTTPInstance()
	h.InitRouter()

	// services
	actuatorRouter := h.RouterGroup(config.Context.Management.BasePath)
	actuator.Controller(actuatorRouter)
}

func run() {
	h := web.GetHTTPInstance()
	m := repository.GetMySQLInstance()
	r := repository.GetRedisInstance()
	context := Context{
		logger.Logger,
		h,
		m,
		r,
	}

	defer func() {
		logger.Logger.Warn("panic occurred, try to recover")
		if err := recover(); err != nil {
			logger.Logger.Error("stacktrace: \n" + string(debug.Stack()))
			logger.Logger.Fatal("really panic: ", err)
		}
	}()

	for _, app := range apps {
		app(&context)
	}
	go h.Run()
	select {}
}