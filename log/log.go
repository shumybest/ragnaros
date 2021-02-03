package log

import (
    "github.com/elastic/go-elasticsearch/v7"
    "github.com/shumybest/ragnaros/config"
    "github.com/sirupsen/logrus"
    "gopkg.in/go-extras/elogrus.v7"
    "sync"
    "time"
)

type Logger struct {
    *logrus.Logger
}

var loggerInstance *Logger
var Once sync.Once

func GetLoggerInstance() *Logger {
    Once.Do(func() {
        loggerInstance = &Logger{}
        loggerInstance.Logger = logrus.New()
    })

    return loggerInstance
}

func (l *Logger) InitConfig() {
    l.SetReportCaller(false)
    level, err := logrus.ParseLevel(config.GetConfigString("logging.level.ROOT"))
    if err == nil {
        l.SetLevel(level)
    } else {
        level = logrus.InfoLevel
    }

    // init elasticsearch logging
    esUrl := config.GetConfigString("ragnaros.elasticsearch.url")
    esHost := config.GetConfigString("ragnaros.elasticsearch.host")
    esPort := config.GetConfigString("ragnaros.elasticsearch.port")
    esUsername := config.GetConfigString("ragnaros.elasticsearch.username")
    esPassword := config.GetConfigString("ragnaros.elasticsearch.password")

    if esUrl == "" {
        if  esHost == "" {
            // no es configured, skip es hook
            return
        } else {
            esUrl = "http://" + esHost + ":" + esPort
        }
    }

    esConfig := elasticsearch.Config{}

    esConfig.Addresses = []string{esUrl}
    if esUsername != "" && esPassword != "" {
        esConfig.Username = esUsername
        esConfig.Password = esPassword
    }

    client, err := elasticsearch.NewClient(esConfig)
    if err != nil {
        l.Error(err)
        return
    }

    indexName := "ragnaros-" + config.Context.ServiceName + "-" + time.Now().Format("2006.01.02")
    hook, err := elogrus.NewAsyncElasticHook(client, config.Context.Hostname, level, indexName)
    if err != nil {
        l.Error(err)
        return
    }

    l.Hooks.Add(hook)
}