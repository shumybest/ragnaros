package main

import (
	"github.com/shumybest/ragnaros"
	"github.com/shumybest/ragnaros/feign"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
	"net/http"
)

type Product struct {
	gorm.Model `json:",omitempty"`
	Code  string `json:"code,omitempty" binding:"required"`
	Price string `json:"price,omitempty" binding:"required"`
}

func DemoController(r *ragnaros.Context) {
	r.MySQLMigration(&Product{})
	demo := r.RouterGroup("/demo")

	demo.GET("/products", func(c *gin.Context) {
		var products []Product
		r.DB.Find(&products)
		c.JSON(http.StatusOK, products)
	})

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

	demo.POST("/product", func(c *gin.Context) {
		var product Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		r.DB.Create(&product)
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

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
