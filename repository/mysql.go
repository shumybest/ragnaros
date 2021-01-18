package repository

import (
	"github.com/shumybest/ragnaros2/config"
	"github.com/shumybest/ragnaros2/eureka"
	"github.com/shumybest/ragnaros2/log"
	"github.com/shumybest/ragnaros2/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"sync"
)

var logger = log.GetLoggerInstance()

type MySQLClient struct {
	*gorm.DB
	Status string
}

var mySQLInstance *MySQLClient
var mySQLOnce sync.Once

func GetMySQLInstance() *MySQLClient {
	mySQLOnce.Do(func() {
		mySQLInstance = &MySQLClient{}
	})
	return mySQLInstance
}

func (m *MySQLClient) InitConnection() {
	connString := strings.TrimSpace(config.GetConfigString("spring.datasource.url"))
	if connString != "" {
		if strings.Contains(connString, "jdbc") {
			connString = utils.JdbcToDSN(connString)
		}
		logger.Infof("connecting to databse: %s\n", connString)
		connString = config.GetConfigString("spring.datasource.username") +
			":" + config.GetConfigString("spring.datasource.password") +
			"@" + connString + "&parseTime=true"

		db, err := gorm.Open(mysql.Open(connString), &gorm.Config{})
		if err != nil {
			panic(err)
			return
		}
		m.DB = db
		m.Status = eureka.UP
		logger.Info("databse connected\n")
	} else {
		logger.Warn("no database configured, continue")
	}
}

func (m *MySQLClient) MySQLMigration(dst ...interface{}) {
	if m.DB != nil {
		_ = m.DB.AutoMigrate(dst...)
	} else {
		panic("no database configured or database not connected, but migration is invoked")
	}
}

