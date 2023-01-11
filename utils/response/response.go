package response

import (
	"github.com/gin-gonic/gin"
)

func SuccessResponse(c *gin.Context,statusCode int,data map[string]string){

	if statusCode < 200 || statusCode > 299{
		statusCode = 200
	}

	c.JSON(statusCode,gin.H{
		"code" : statusCode,
		"success" : true,
		"data" : data,
	})
}

func ErrorResponse(c *gin.Context,statusCode int,data map[string]string){

	c.JSON(statusCode,gin.H{
		"code" : statusCode,
		"success" : false,
		"data" : data,
	})
}
