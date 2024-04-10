package v1

import (
	v1 "github.com/GGmaz/wallet-arringo/internal/server/api/v1"
	"github.com/gin-gonic/gin"
)

func RegisterVersion(router *gin.Engine) {
	router.Group("/api/v1")
	v1.RegisterUser(router)
	v1.RegisterTransaction(router)
	v1.RegisterAccount(router)
}
