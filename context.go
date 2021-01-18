package ragnaros

import (
	"github.com/shumybest/ragnaros2/log"
	"github.com/shumybest/ragnaros2/repository"
	"github.com/shumybest/ragnaros2/web"
)

type Context struct {
	*log.Logger
	*web.HttpServer
	*repository.MySQLClient
	*repository.RedisClient
}
