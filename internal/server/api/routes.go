package v1

import (
	v1 "github.com/GGmaz/wallet-arringo/internal/server/api/v1"
	"github.com/gin-gonic/gin"
)

func RegisterVersion(router *gin.Engine) {
	r := router.Group("/api/v1")
	v1.RegisterUser(r)
	v1.RegisterTransaction(r)
	v1.RegisterAccount(r)
}
