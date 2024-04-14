package v1

import (
	"github.com/GGmaz/wallet-arringo/internal/db/model"
	"github.com/GGmaz/wallet-arringo/pkg/wire"
	"github.com/gin-gonic/gin"
)

// RegisterUser registers user-related HTTP endpoints on the provided RouterGroup.
func RegisterUser(v1 *gin.RouterGroup) {
	v1 = v1.Group("/users")
	{
		v1.POST("", createUser)
		v1.GET("/:email/accounts", getAccounts)
	}
}

// createUser handles the HTTP request to create a new user.
func createUser(c *gin.Context) {
	// Parse the request body to create a new user
	var createUser model.User
	err := c.BindJSON(&createUser)
	if err != nil {
		c.PureJSON(500, gin.H{"error": err.Error()})
		return
	}

	// Validate the request parameters
	if createUser.Email == "" {
		c.PureJSON(400, gin.H{"error": "email is required"})
		return
	}

	// Create the user using the UserService
	user, err := wire.Svc.UserService.CreateUser(c, createUser.Email, createUser.FirstName, createUser.LastName, createUser.Address)
	if err != nil {
		c.PureJSON(500, gin.H{"error": err.Error()})
		return
	}

	// Respond with the created user
	c.PureJSON(200, user)
}

// getAccounts handles the HTTP request to retrieve accounts associated with a user.
func getAccounts(c *gin.Context) {
	// Parse the user email from the request parameters
	userEmail := c.Param("email")
	if userEmail == "" {
		c.PureJSON(400, gin.H{"error": "email is required"})
		return
	}

	// Retrieve accounts associated with the user using the UserService
	res, err := wire.Svc.UserService.GetAccounts(c, userEmail)
	if err != nil {
		c.PureJSON(500, gin.H{"error": err.Error()})
		return
	}

	// Respond with the user's accounts
	c.PureJSON(200, res)
}
