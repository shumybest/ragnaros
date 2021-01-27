package ragnaros

import (
	"github.com/shumybest/ragnaros/actuator"
	"github.com/shumybest/ragnaros/config"
	"github.com/shumybest/ragnaros/eureka"
	"github.com/shumybest/ragnaros/log"
	"github.com/shumybest/ragnaros/repository"
	"github.com/shumybest/ragnaros/web"
	"github.com/urfave/cli/v2"
	"runtime/debug"
)

type injectedApp func(*Context)
var apps []injectedApp

func InjectApps(injectedFunc ...injectedApp) {
	apps = append(apps, injectedFunc...)
}

// with options
func StartOpts(c *cli.Context) error {
	config.Init(c)
	initComponents()
	run()

	return nil
}

// no options
func Start() {
	config.Init()
	initComponents()
	run()
}

func initComponents() {
	l := log.GetLoggerInstance()
	l.InitConfig()

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
	l := log.GetLoggerInstance()
	context := Context{
		l,
		h,
		m,
		r,
	}

	defer func() {
		l.Warn("panic occurred, try to recover")
		if err := recover(); err != nil {
			l.Error("stacktrace: \n" + string(debug.Stack()))
			l.Fatal("really panic: ", err)
		}
	}()

	for _, app := range apps {
		app(&context)
	}
	go h.Run()
	select {}
}