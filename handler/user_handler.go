package handler

import (
	"Permanent-Storage-Todo/model"
	"Permanent-Storage-Todo/service"
	"errors"

	"github.com/gin-gonic/gin"
)

type User_handler struct{
	service service.User_service
}

func Init_user_handler(service service.User_service) *User_handler{
	new_handler := User_handler{service: service}
	return &new_handler
}


// @Summary	 		Register a new user.
// @Description 	Create a new user account with provided username and password and default role (user).
// @Tags			User
// @Accept			json
// @Produce			json
// @Param			request body 	model.Auth_request true 	"User signup credentials"
// @Success			201  {object}	model.User					"Created user"
// @Failure			400	 {object}	map[string]string			"Error message"
// @Failure			500	 {object}	map[string]string			"Server error"
// @Router			/signup [post]
func (handler *User_handler) Create_user_handler(c *gin.Context){
	//get the data form json body.
	var user_request model.Auth_request
	err := c.ShouldBindJSON(&user_request)
	//if json is invalid.
	if err != nil{
		c.JSON(400, gin.H{"error": "Invalid json"})
		return
	}
	//create a new user with the user input.
	user_ptr, err := handler.service.Create_user(user_request)
	//check the result.
	switch {
	case err == nil:
		c.JSON(201, *user_ptr)
	case errors.Is(err, service.ErrEmptyName):
		c.JSON(400, gin.H{"error": "Username cannot be empty"})
	case errors.Is(err, service.ErrEmptyPassword):
		c.JSON(400, gin.H{"error": "Password cannot be empty"})
	case errors.Is(err, service.ErrInvalidPassword):
		c.JSON(400, gin.H{"error": "Invalid password"})
	case errors.Is(err, service.ErrDuplicatedName):
		c.JSON(400, gin.H{"error": "Username is already used"})
	default:
		c.JSON(500, gin.H{"error": "Server error"})
	}
}

// @Summary	 		User login.
// @Description 	Authenticate user credentials and generate a JWT token.
// @Tags			User
// @Accept			json
// @Produce			json
// @Param			request body 	model.Auth_request true 	"User login credentials"
// @Success			200  {object}	map[string]string			"Authentication JWT token"
// @Failure			400	 {object}	map[string]string			"Invalid json format"
// @Failure			401	 {object}	map[string]string			"Invalid credentials"
// @Failure			500	 {object}	map[string]string			"Server error"
// @Router			/login [post]
func (handler *User_handler) Login_user_handler(c *gin.Context){
	//get the user input from the json body.
	var user_request model.Auth_request
	err := c.ShouldBindJSON(&user_request)
	//check the result.
	if err != nil{
		c.JSON(400, gin.H{"error": "Invalid json"})
		return
	}
	//look up the user.
	token, err := handler.service.Login_user(user_request)
	//check the result of process.
	switch {
	case err == nil:
		c.JSON(200, gin.H{"token": token})
	case errors.Is(err, service.ErrUserNotFound), errors.Is(err, service.ErrWrongPassword):
		c.JSON(401, gin.H{"error": "Invalid credentials"})
	default:
		c.JSON(500, gin.H{"error": "Server error"})
	}
}