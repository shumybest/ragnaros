package ragnaros

import (
	"github.com/shumybest/ragnaros/log"
	"github.com/shumybest/ragnaros/repository"
	"github.com/shumybest/ragnaros/web"
)

type Context struct {
	*log.Logger
	*web.HttpServer
	*repository.MySQLClient
	*repository.RedisClient
}
