package middleware

import (
	"a-easy-memo/pkg/utils"
	"errors"
	"fmt"

	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header.Get("Authorization")
		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"question": "请重新登录",
			})
			c.Abort()
			return
		}
		heade := Spilt(header)
		token, err := jwt.ParseWithClaims(heade, &utils.MyClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("有错误")
			}
			return []byte("vivo50"), nil
		})
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"message": "令牌不合格",
			})
			c.Abort()
			return
		}
		if claims, ok := token.Claims.(*utils.MyClaims); ok && token.Valid {
			c.Set("username", claims.UserName)
			c.Next()
		} else {
			c.JSON(401, gin.H{
				"message": "令牌过期或无效",
			})
			c.Abort()
			return
		}
	}
}

func Spilt(s string) string {
	if strings.HasPrefix(strings.ToLower(strings.TrimSpace(s)), "bearer") {
		return strings.TrimSpace(s[len("bearer"):])
	}
	return strings.TrimSpace(s)
}
