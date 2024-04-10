package v1

import (
	"github.com/GGmaz/wallet-arringo/internal/db/model"
	"github.com/GGmaz/wallet-arringo/pkg/utils"
	"github.com/GGmaz/wallet-arringo/pkg/wire"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUser(router *gin.Engine) {
	v1 := router.Group("/users")
	{
		v1.POST("/", createUser)
		v1.GET("/:email/balance/:acc", getBalance)
	}
}

func createUser(c *gin.Context) {
	var createUser model.User
	err := c.BindJSON(&createUser)
	utils.Handle(err)
	if createUser.Email == "" {
		c.PureJSON(400, gin.H{"error": "email is required"})
		return
	}

	db := c.MustGet("transaction").(*gorm.DB)
	user, err := wire.Svc.UserService.CreateUser(db, createUser.Email)
	if err != nil {
		c.PureJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.PureJSON(200, user)
}

func getBalance(c *gin.Context) {
	// Parse the user email from the request parameters
	userEmail := c.Param("email")
	accNum := c.Param("acc")
	if userEmail == "" || accNum == "" {
		c.PureJSON(400, gin.H{"error": "email and account number are required"})
		return
	}

	if err != nil {
		c.PureJSON(500, gin.H{"error": err.Error()})
		return
	}

	// Parse the response and send the HTTP response with email and balance
	c.PureJSON(200, gin.H{"email": userEmail, "balance": string(response.Data)})
}
