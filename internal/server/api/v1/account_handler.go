package v1

import (
	"github.com/GGmaz/wallet-arringo/internal/db/model"
	"github.com/GGmaz/wallet-arringo/pkg/utils"
	"github.com/GGmaz/wallet-arringo/pkg/wire"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAccount(router *gin.Engine) {
	v1 := router.Group("/accounts")
	{
		v1.POST("/", createAccount)
	}
}

func createAccount(c *gin.Context) {
	var createAcc model.Account
	err := c.BindJSON(&createAcc)
	utils.Handle(err)
	if createAcc.AccNumber == "" || createAcc.UserId <= 0 {
		c.PureJSON(400, gin.H{"error": "account number and valid userId are required"})
		return
	}

	db := c.MustGet("transaction").(*gorm.DB)
	account, err := wire.Svc.AccountService.CreateAccount(db, createAcc.AccNumber, createAcc.UserId)
	if err != nil {
		c.PureJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.PureJSON(200, account)
}
