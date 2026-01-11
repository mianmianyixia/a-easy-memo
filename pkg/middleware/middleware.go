package middleware

import (
	"a-easy-memo/pkg/utils"
	"a-easy-memo/zlog"
	"errors"
	"fmt"
	"time"

	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
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
func ZapLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}
		c.Next()
		end := time.Now()
		latency := end.Sub(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		responseSize := c.Writer.Size()
		errs := c.Errors.ByType(gin.ErrorTypePrivate).String()
		fields := []zap.Field{
			zap.String("clientIP", clientIP),
			zap.Int("statusCode", statusCode),
			zap.Duration("latency", latency),
			zap.Int("bytes", responseSize),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("time", end.Format(time.RFC3339)),
		}
		if errs != "" {
			fields = append(fields, zap.String("errors", errs))
		}
		if statusCode >= 400 && statusCode < 500 {
			zlog.Warn(fmt.Sprintf("%s %s", method, path), fields...)
		} else if statusCode >= 500 {
			zlog.Error(fmt.Sprintf("%s %s", method, path), fields...)
		} else {
			zlog.Info(fmt.Sprintf("%s %s", method, path), fields...)
		}
	}
}

func Spilt(s string) string {
	if strings.HasPrefix(strings.ToLower(strings.TrimSpace(s)), "bearer") {
		return strings.TrimSpace(s[len("bearer"):])
	}
	return strings.TrimSpace(s)
}
