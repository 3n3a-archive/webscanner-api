package validate

import "github.com/gin-gonic/gin"

func JsonError(err error, statusCode int, c *gin.Context)  {
	c.JSON(statusCode, gin.H{
		"message": err.Error(),
	})
}

func IsErrorState(err error) bool {
	if err != nil  {
		return true
	}
	return false
}