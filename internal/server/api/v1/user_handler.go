package v1

import (
	"github.com/GGmaz/wallet-arringo/internal/db/model"
	"github.com/GGmaz/wallet-arringo/pkg/wire"
	"github.com/gin-gonic/gin"
)

func RegisterUser(v1 *gin.RouterGroup) {
	v1 = v1.Group("/users")
	{
		v1.POST("", createUser)
		v1.GET("/:email/accounts", getAccounts)
	}
}

func createUser(c *gin.Context) {
	var createUser model.User
	err := c.BindJSON(&createUser)
	if err != nil {
		c.PureJSON(500, gin.H{"error": err.Error()})
		return
	}

	if createUser.Email == "" {
		c.PureJSON(400, gin.H{"error": "email is required"})
		return
	}

	user, err := wire.Svc.UserService.CreateUser(c, createUser.Email, createUser.FirstName, createUser.LastName, createUser.Address)
	if err != nil {
		c.PureJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.PureJSON(200, user)
}

func getAccounts(c *gin.Context) {
	// Parse the user email from the request parameters
	userEmail := c.Param("email")
	if userEmail == "" {
		c.PureJSON(400, gin.H{"error": "email is required"})
		return
	}

	res, err := wire.Svc.UserService.GetAccounts(c, userEmail)
	if err != nil {
		c.PureJSON(500, gin.H{"error": err.Error()})
		return
	}

	// Parse the response and send the HTTP response with email and balance
	c.PureJSON(200, res)
}
