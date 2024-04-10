package v1

import "github.com/gin-gonic/gin"

func RegisterUser(router *gin.Engine) {
	v1 := router.Group("/api/v1/transactions")
	{
		v1.POST("/add-money", addMoney)
		v1.POST("/transfer-money", transferMoney)
	}
}
