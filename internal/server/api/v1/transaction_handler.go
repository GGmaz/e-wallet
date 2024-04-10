package v1

import (
	"github.com/GGmaz/wallet-arringo/pkg/models"
	"github.com/GGmaz/wallet-arringo/pkg/utils"
	"github.com/GGmaz/wallet-arringo/pkg/wire"
	"github.com/gin-gonic/gin"
)

func RegisterTransaction(router *gin.Engine) {
	v1 := router.Group("/api/v1/transactions")
	{
		v1.POST("/add-money", addMoney)
		v1.POST("/transfer-money", transferMoney)
	}
}

func transferMoney(c *gin.Context) {
	var req models.TransferMoneyReq
	err := c.BindJSON(&req)
	utils.Handle(err)
	if req.From == 0 || req.To == 0 || req.Amount == 0 {
		c.PureJSON(400, gin.H{"error": "from_user_id and to_user_id and amount are required and can not be 0"})
		return
	}

	err = wire.Svc.TransactionService.TransferMoney(c, req.From, req.To, req.Amount)
	if err != nil {
		c.PureJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Status(200)
}

func addMoney(c *gin.Context) {
	var req models.AddMoneyReq
	err := c.BindJSON(&req)
	utils.Handle(err)
	if req.UserID == 0 || req.Amount == 0 {
		c.PureJSON(400, gin.H{"error": "userId and amount are required and can not be 0"})
		return
	}

	res, err := wire.Svc.TransactionService.AddMoney(c, req.UserID, req.Amount)
	if err != nil {
		c.PureJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.PureJSON(200, models.AddMoneyRes{Balance: res})
}
