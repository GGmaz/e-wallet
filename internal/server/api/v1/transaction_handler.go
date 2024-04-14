package v1

import (
	"github.com/GGmaz/wallet-arringo/pkg/models"
	"github.com/GGmaz/wallet-arringo/pkg/wire"
	"github.com/gin-gonic/gin"
)

// RegisterTransaction registers transaction-related HTTP endpoints on the provided RouterGroup.
func RegisterTransaction(v1 *gin.RouterGroup) {
	v1 = v1.Group("/transactions")
	{
		v1.POST("/deposit", deposit)
		v1.POST("/withdraw", withdraw)
		v1.POST("/transfer", transferMoney)
	}
}

// deposit handles the HTTP request to deposit money into an account.
// @Tags Transaction
// @Summary Deposit Money
// @Description Deposit money into an account
// @ID deposit money
// @Accept  json
// @Produce  json
// @Param deposit body models.AddMoneyReq true "Deposit" default({"userId":1, "accNumber":"1234", "amount":100})
// @Success 200 {float} float "ok"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /transactions/deposit [post]
func deposit(c *gin.Context) {
	// Parse the request body to extract deposit details
	var req models.AddMoneyReq
	err := c.BindJSON(&req)
	if err != nil {
		c.PureJSON(500, gin.H{"error": err.Error()})
		return
	}

	// Validate the request parameters
	if req.UserID <= 0 || req.Amount <= 0 || req.AccNumber == "" {
		c.PureJSON(400, gin.H{"error": "userId, account number, and amount are required and cannot be <=0 or empty"})
		return
	}

	// Call the TransactionService to deposit money
	res, err := wire.Svc.TransactionService.AddMoney(c, req.UserID, req.Amount, req.AccNumber)
	if err != nil {
		c.PureJSON(500, gin.H{"error": err.Error()})
		return
	}

	// Respond with the updated balance
	c.PureJSON(200, models.AddMoneyRes{Balance: res})
}

// withdraw handles the HTTP request to withdraw money from an account.
// @Tags Transaction
// @Summary Withdraw Money
// @Description Withdraw money from an account
// @ID withdraw money
// @Accept  json
// @Produce  json
// @Param withdraw body models.AddMoneyReq true "Withdraw" default({"userId":1, "accNumber":"1234", "amount":10})
// @Success 200 {float} float "ok"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /transactions/withdraw [post]
func withdraw(c *gin.Context) {
	// Parse the request body to extract withdrawal details
	var req models.AddMoneyReq
	err := c.BindJSON(&req)
	if err != nil {
		c.PureJSON(500, gin.H{"error": err.Error()})
		return
	}

	// Validate the request parameters
	if req.UserID <= 0 || req.Amount <= 0 || req.AccNumber == "" {
		c.PureJSON(400, gin.H{"error": "userId, account number, and amount are required and cannot be <=0 or empty"})
		return
	}

	// Call the TransactionService to withdraw money
	res, err := wire.Svc.TransactionService.Withdraw(c, req.UserID, req.Amount, req.AccNumber)
	if err != nil {
		c.PureJSON(500, gin.H{"error": err.Error()})
		return
	}

	// Respond with the updated balance
	c.PureJSON(200, models.AddMoneyRes{Balance: res})
}

// transferMoney handles the HTTP request to transfer money between accounts.
// @Tags Transaction
// @Summary Transfer Money
// @Description Transfer money between accounts
// @ID transfer money
// @Accept  json
// @Produce  json
// @Param transfer body models.TransferMoneyReq true "Transfer" default({"from_acc_id":"1234", "to_acc_id":"5678", "amount":10})
// @Success 200 "ok"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /transactions/transfer [post]
func transferMoney(c *gin.Context) {
	// Parse the request body to extract transfer details
	var req models.TransferMoneyReq
	err := c.BindJSON(&req)
	if err != nil {
		c.PureJSON(500, gin.H{"error": err.Error()})
		return
	}

	// Validate the request parameters
	if req.From == "" || req.To == "" || req.Amount <= 0 {
		c.PureJSON(400, gin.H{"error": "from_acc_id, to_acc_id, and amount are required and cannot be <=0 or empty"})
		return
	}

	// Call the TransactionService to transfer money
	err = wire.Svc.TransactionService.TransferMoney(c, req.From, req.To, req.Amount)
	if err != nil {
		c.PureJSON(500, gin.H{"error": err.Error()})
		return
	}

	// Respond with success
	c.Status(200)
}
