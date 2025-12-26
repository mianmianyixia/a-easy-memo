package response

import "github.com/gin-gonic/gin"

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"message": data,
	})
}
func InternalError(c *gin.Context, err error) {
	c.JSON(500, gin.H{
		"message": err.Error(),
	})
}
func RequestError(c *gin.Context, err error) {
	c.JSON(400, gin.H{
		"message": err.Error(),
	})
}
