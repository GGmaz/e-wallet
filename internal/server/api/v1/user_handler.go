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
// createUser handles the HTTP request to create a new user.
// @Tags User
// @Summary Create User
// @Description Create a new user
// @ID create user
// @Accept  json
// @Produce  json
// @Param X-Authorization-Sign header string false "X-Authorization-Sign"
// @Param user body model.User true "User" default({"first_name":"srecko", "last_name":"asdas", "email":"srecko@gmail.com", "address":"asdasd"})
// @Success 200 {object} model.User "ok"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /users [post]
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
// @Tags User
// @Summary Get Accounts
// @Description Get accounts associated with a user
// @ID get accounts
// @Accept  json
// @Produce  json
// @Param email path string true "User Email" default(srecko@gmail.com)
// @Success 200 {object} model.User "ok"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /users/{email}/accounts [get]
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
