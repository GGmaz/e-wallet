package v1

import (
	"github.com/GGmaz/wallet-arringo/pkg/models"
	"github.com/GGmaz/wallet-arringo/pkg/wire"
	"github.com/gin-gonic/gin"
)

func RegisterTransaction(v1 *gin.RouterGroup) {
	v1 = v1.Group("/transactions")
	{
		v1.POST("/deposit", deposit)
		v1.POST("/withdraw", withdraw)
		v1.POST("/transfer", transferMoney)
	}
}

func transferMoney(c *gin.Context) {
	//TODO: check if account are verified
	var req models.TransferMoneyReq
	err := c.BindJSON(&req)
	if err != nil {
		c.PureJSON(500, gin.H{"error": err.Error()})
		return
	}

	if req.From == "" || req.To == "" || req.Amount <= 0 {
		c.PureJSON(400, gin.H{"error": "from_acc_id and to_acc_id and amount are required and can not be <=0 or empty"})
		return
	}

	err = wire.Svc.TransactionService.TransferMoney(c, req.From, req.To, req.Amount)
	if err != nil {
		c.PureJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Status(200)
}

func deposit(c *gin.Context) {
	//TODO: check if account are verified
	var req models.AddMoneyReq
	err := c.BindJSON(&req)
	if err != nil {
		c.PureJSON(500, gin.H{"error": err.Error()})
		return
	}

	if req.UserID <= 0 || req.Amount <= 0 || req.AccNumber == "" {
		c.PureJSON(400, gin.H{"error": "userId, account number and amount are required and can not be <=0 or empty"})
		return
	}

	res, err := wire.Svc.TransactionService.AddMoney(c, req.UserID, req.Amount, req.AccNumber)
	if err != nil {
		c.PureJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.PureJSON(200, models.AddMoneyRes{Balance: res})
}

func withdraw(c *gin.Context) {
	//TODO: check if account are verified
	var req models.AddMoneyReq
	err := c.BindJSON(&req)
	if err != nil {
		c.PureJSON(500, gin.H{"error": err.Error()})
		return
	}

	if req.UserID <= 0 || req.Amount <= 0 || req.AccNumber == "" {
		c.PureJSON(400, gin.H{"error": "userId, account number and amount are required and can not be <=0 or empty"})
		return
	}

	res, err := wire.Svc.TransactionService.Withdraw(c, req.UserID, req.Amount, req.AccNumber)
	if err != nil {
		c.PureJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.PureJSON(200, models.AddMoneyRes{Balance: res})
}
