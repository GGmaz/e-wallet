package v1

import (
	"github.com/GGmaz/wallet-arringo/internal/db/model"
	"github.com/GGmaz/wallet-arringo/pkg/wire"
	"github.com/gin-gonic/gin"
)

func RegisterAccount(v1 *gin.RouterGroup) {
	v1 = v1.Group("/accounts")
	{
		v1.POST("", createAccount)
	}
}

func createAccount(c *gin.Context) {
	var createAcc model.Account
	err := c.BindJSON(&createAcc)
	if err != nil {
		c.PureJSON(500, gin.H{"error": err.Error()})
		return
	}

	if createAcc.AccNumber == "" || createAcc.UserId <= 0 {
		c.PureJSON(400, gin.H{"error": "account number and valid userId are required"})
		return
	}

	account, err := wire.Svc.AccountService.CreateAccount(c, createAcc.AccNumber, createAcc.UserId)
	if err != nil {
		c.PureJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.PureJSON(200, account)
}
