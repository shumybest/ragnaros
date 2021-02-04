package security

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/shumybest/ragnaros/config"
	. "github.com/shumybest/ragnaros/logger"
	"net/http"
	"strings"
)

func ValidateToken(tokenStr string) (*config.Claims) {
	token, err := jwt.ParseWithClaims(tokenStr,
		&config.Context.Security.Claims, func(token *jwt.Token) (interface{}, error) {
		return config.Context.Security.JwtSecret, nil
	})

	if err == nil {
		if claims, ok := token.Claims.(*config.Claims); ok && token.Valid {
			return claims
		}
	}

	Logger.Error(err)
	return nil
}

func AuthorizationInterceptor(c *gin.Context) {
	authorString := c.Request.Header.Get("Authorization")
	Logger.Debug("Checking authentication: " + authorString)

	if authorString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": 401,
			"message": "未授权的访问",
			"data": nil,
		})
		c.Abort()
		return
	}

	authorArray := strings.SplitN(authorString, " ", 2)
	if len(authorArray) != 2 || authorArray[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": 401,
			"message": "未授权的访问",
			"data": nil,
		})
		c.Abort()
		return
	}

	claims := ValidateToken(strings.TrimSpace(authorArray[1]))
	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": 401,
			"message": "未授权的访问",
			"data": nil,
		})
		c.Abort()
		return
	}

	c.Next()
}
