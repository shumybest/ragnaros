package logger

import (
    "github.com/elastic/go-elasticsearch/v7"
    "github.com/shumybest/ragnaros/config"
    "github.com/sirupsen/logrus"
    "gopkg.in/go-extras/elogrus.v7"
    "sync"
    "time"
)

var Logger *logrus.Entry
var Once sync.Once

func InitLogger() {
    Once.Do(func() {
        logger := logrus.New()
        logger.SetReportCaller(true)
        level, err := logrus.ParseLevel(config.GetConfigString("logging.level.ROOT"))
        if err == nil {
            logger.SetLevel(level)
        } else {
            level = logrus.InfoLevel
        }

        // init elasticsearch logging, if no es address configured, skip es hook
        esUrl := config.GetConfigString("ragnaros.elasticsearch.url")
        esHost := config.GetConfigString("ragnaros.elasticsearch.host")
        esPort := config.GetConfigString("ragnaros.elasticsearch.port")
        esUsername := config.GetConfigString("ragnaros.elasticsearch.username")
        esPassword := config.GetConfigString("ragnaros.elasticsearch.password")

        if esUrl != "" || esHost != "" {
            if esUrl == "" {
                esUrl = "http://" + esHost + ":" + esPort
            }

            esConfig := elasticsearch.Config{}
            esConfig.Addresses = []string{esUrl}
            if esUsername != "" && esPassword != "" {
                esConfig.Username = esUsername
                esConfig.Password = esPassword
            }

            client, err := elasticsearch.NewClient(esConfig)
            if err != nil {
                logger.Error(err)
                return
            }

            indexName := "ragnaros-" + config.GetConfigString("eureka.instance.appname") +
                "-" + time.Now().Format("2006.01.02")
            hook, err := elogrus.NewAsyncElasticHook(client, config.Context.Hostname, level, indexName)
            if err != nil {
                logger.Error(err)
                return
            }

            logger.Hooks.Add(hook)
        }

        logEntry := logger.WithFields(logrus.Fields{
            "app_name": config.GetConfigString("eureka.instance.appname"),
            "instance": config.GetConfigString("ragnaros.instanceid"),
        })
        Logger = logEntry
    })
}