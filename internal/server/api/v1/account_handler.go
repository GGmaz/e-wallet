package v1

import (
	"github.com/GGmaz/wallet-arringo/internal/db/model"
	"github.com/GGmaz/wallet-arringo/pkg/wire"
	"github.com/gin-gonic/gin"
)

// RegisterAccount registers account-related HTTP endpoints on the provided RouterGroup.
func RegisterAccount(v1 *gin.RouterGroup) {
	v1 = v1.Group("/accounts")
	{
		v1.POST("", createAccount)
	}
}

// createAccount handles the HTTP request to create a new account.
// @Tags Account
// @Summary Create Account
// @Description Create a new account for user
// @ID create account
// @Accept  json
// @Produce  json
// @Param account body model.Account true "Account" default({"acc_number":"1234", "user_id":1})
// @Success 200 {object} model.Account "ok"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /accounts [post]
func createAccount(c *gin.Context) {
	// Parse the request body to create a new account
	var createAcc model.Account
	err := c.BindJSON(&createAcc)
	if err != nil {
		c.PureJSON(500, gin.H{"error": err.Error()})
		return
	}

	// Validate the request parameters
	if createAcc.AccNumber == "" || createAcc.UserId <= 0 {
		c.PureJSON(400, gin.H{"error": "account number and valid userId are required"})
		return
	}

	// Create the account using the AccountService
	account, err := wire.Svc.AccountService.CreateAccount(c, createAcc.AccNumber, createAcc.UserId)
	if err != nil {
		c.PureJSON(500, gin.H{"error": err.Error()})
		return
	}

	// Respond with the created account
	c.PureJSON(200, account)
}
