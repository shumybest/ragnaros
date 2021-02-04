package ragnaros

import (
	"github.com/shumybest/ragnaros/repository"
	"github.com/shumybest/ragnaros/web"
	"github.com/sirupsen/logrus"
)

type Context struct {
	*logrus.Entry
	*web.HttpServer
	*repository.MySQLClient
	*repository.RedisClient
}
