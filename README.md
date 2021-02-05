Ragnaros
========

a go implementation of Spring Cloud Netflix Microservice framework
--------

#### Quick Start

Install the Ragnaros command line tool for quick project generation.
```shell
go get github.com/shumybest/ragnaros/cmd/ragnaros
```

Ragnaros usage reference, <[configuration.json](example/conf.json)> example reference [here](example/conf.json):
```shell
ragnaros help

ragnaros generate -c <configuration.json> -o <destination>
```

Then go into the generated project, just start with simply:
```shell
cd <destination>
make docker
```

Force downloading the project template files again
```shell
ragnaros download --force
```

#### Features
- [x] support Eureka Register and Zuul Routing to this service
- [x] support of Registry (currently jHipster) health check
- [x] Application level injection for framework
- [x] MySQL ORM support
- [x] Redis support
- [x] go implementation of FeignClient
- [x] Dockerized support: [Dockerfile example](http://github.com/shumybest/ragnaros-example)
- [x] SpringCloud Config integration
- [x] go tools for project generation
- [x] async invoke of FeignClient
- [x] logging support
- [x] logging elk support
- [ ] MQ support (like kafka)
- [ ] deeper wrapper of resty for ribbon/hystrix implementation
- [ ] tracing support (like skywalking)
- [ ] monitoring
- [ ] migration go cmd for existing Java microservice
- [ ] more registry support (like Consul)
- [ ] more database support (like mongodb, postgreSQL, clickhouse)

#### Why Ragnaros
golang is significantly better than Java both on less resources usage and higher performance of non-blocking IO, but lots services are written in Java/Spring framework. Ragnaros aims to write the production ready microservice which is compatible with Java spring cloud framework, and also help to migrate the Java services to go.

Comparison | golang | java
--------------------- | ------------------ | ----------------
static executable | tens MB  | hundreds MB
container image | tens MB  | hundreds MB
runtime memory | tens MB | hundreds MB even GB
cross platform | native | by JVM
non-blocking IO performance | 5000 req/sec | 2500 req/sec

#### Requirements
- minimal go version: 1.13.4
- recommend go version: 1.15
- go env -w GOPROXY="https://mirrors.aliyun.com/goproxy/,direct"

#### How To Use
- [usage example](http://github.com/shumybest/ragnaros-example)
- copy java application yml configure files to runtime directory (default directory: **./resources/config**)
- configuration file loading sequence: **bootstrap.yml -> bootstrap-<profile>.yml -> application.yml -> application-<profile>.yml**
- environment variables will overwrite the configuration in file, like SERVER_PORT to overwrite server.port in file
- inject your application level services (one or more) by [InjectApps](bootstrap.go), reference: [main.go](example/main.go)

```go
package main

import (
	"github.com/shumybest/ragnaros"
	"ragnaros-example/app"
)

func main() {
	ragnaros.InjectApps(app.DemoController)

    // also you can inject many application level implementations
	ragnaros.InjectApps(a.Controller, b.Controller, c.Controller)

    // or callback function(s)
	ragnaros.InjectApps(func(r *ragnaros.Context) {
		r.Logger.Println("welcome to use ragnaros")
	})

    // then start the microservice afterwards
	ragnaros.Start()
}
```

- implement the service by http router register and http handling, reference: [demo.go](example/demo.go); use default argument [ragnaros.Context](context.go) to access HTTP Engine (gin), MySQL (gorm) and Redis (go-redis).

```go
package main

import (
	"github.com/shumybest/ragnaros"
	"github.com/shumybest/ragnaros/feign"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"net/http"
)

type Product struct {
	gorm.Model
	Code  string `json:"code" binding:"required"`
	Price string `json:"price" binding:"required"`
}

func DemoController(r *ragnaros.Context) {
    r.MySQLMigration(&Product{})
	demo := r.RouterGroup("/demo")

	demo.GET("/products", func(c *gin.Context) {
		var products []Product
		result := r.DB.Find(&products) // result.RowsAffected
        if result.RowsAffected > 0 {
		    c.JSON(http.StatusOK, products)
        }
	})
}
````

- using redis as cache:

````go
func DemoController(r *ragnaros.Context) {
	demo.GET("/product/:code", func(c *gin.Context) {
		code := c.Param("code")
		var product Product

		cache, err := r.RedisGet("products:" + code)
		if err == redis.Nil {
			result := r.DB.First(&product, "code = ?", code)

			if result.RowsAffected > 0 {
				jsonStr, _ := json.Marshal(product)
				r.RedisSet("products:" + code, jsonStr, 0)
				c.JSON(http.StatusOK, gin.H{"message": "success", "data": product})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "success", "data": nil})
			}
			return
		}

		_ = json.Unmarshal([]byte(cache), &product)
		c.JSON(http.StatusOK, gin.H{"message": "success", "data": product})
	})
}
````

- refer the usage of feign client, the internal http invoke between microservices run upon the app name registered to eureka

````go
func DemoController(r *ragnaros.Context) {
	demo.GET("/feignRevoke", func(c *gin.Context) {
		aiboxClient := feign.App("aibox")
		aiboxClient.SetHeaders(feign.Headers{"Authorization": c.GetHeader("Authorization")})
		resp, err := aiboxClient.Get("/management/health")

		if err == nil {
			var raw map[string]interface{}
			_ = json.Unmarshal([]byte(resp.String()), &raw)
			c.JSON(http.StatusOK, gin.H{"message": raw})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
	})

	demo.GET("/feignAsyncRevoke", func(c *gin.Context) {
		aiboxClient := feign.App("aibox")
		aiboxClient.SetHeaders(feign.Headers{"Authorization": c.GetHeader("Authorization")})
		err := aiboxClient.AsyncGet("/management/health", func(response *resty.Response) {
			r.Logger.Info(response)
		})

		if err == nil {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
	})
}
````

#### Supported Environment Variables
Environment Variables | yml configuration field | default
--------------------- | ------------------ | ----------------
RAGNAROS_CONF_DIR |  | resources/config
SPRING_CLOUD_CONFIG_URI | spring.cloud.config.uri |
SPRING_PROFILES_ACTIVE | spring.profiles.active | dev
EUREKA_CLIENT_SERVICE_URL_DEFAULTZONE | eureka.client.service-url.defaultZone |
SERVER_PORT | servr.port | 8999
SPRING_DATASOURCE_URL | spring.datasource.url |
SPRING_DATASOURCE_USERNAME | spring.datasource.username |
SPRING_DATASOURCE_PASSWORD | spring.datasource.password |
SPRING_REDIS_HOST | spring.redis.host | 127.0.0.1
SPRING_REDIS_PORT | spring.redis.port | 6379
RAGNAROS_ELASTICSEARCH_URL | ragnaros.elasticsearch.url |
RAGNAROS_ELASTICSEARCH_HOST | ragnaros.elasticsearch.host |
RAGNAROS_ELASTICSEARCH_PORT | ragnaros.elasticsearch.port | 9200
RAGNAROS_ELASTICSEARCH_USERNAME | ragnaros.elasticsearch.username |
RAGNAROS_ELASTICSEARCH_PASSWORD | ragnaros.elasticsearch.password |

#### Module list (many thanks to the awesome projects)
- [gin](https://github.com/gin-gonic/gin)
- [gorm](https://gorm.io/)
- [go-redis](https://github.com/go-redis/redis)
- [go-resty](https://github.com/go-resty/resty)
- [mergo](https://github.com/imdario/mergo)
- [yaml](https://gopkg.in/yaml.v3)
- [logrus](https://github.com/sirupsen/logrus)
- [jwt-go](https://github.com/dgrijalva/jwt-go)
- [urfave-cli](https://github.com/urfave/cli)
- [gopsutil](https://github.com/shirou/gopsutil)
- ... and all the dependencies

#### FAQ
- unknown or incorrect time zone: 'Asia/Shanghai'
  - solution: https://dev.mysql.com/downloads/timezones.html

#### Authors
- lijingtao (shumybest@163.com)