package handler

import (
	"Permanent-Storage-Todo/model"
	"Permanent-Storage-Todo/service"
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidClaims = errors.New("Error invalid claims")
)

type Todo_handler struct {
	service service.Todo_service
}

func Init_handler(service service.Todo_service) *Todo_handler {
	new_handler := Todo_handler{service: service}
	return &new_handler
}

func Get_user_claims (c *gin.Context) (*model.User_claims, error) {
	//get the user from gin context.
	value, exists := c.Get("user")
	if !exists {
		return nil, ErrInvalidClaims
	}
	//extract the user struct from interface.
	user, ok := value.(*model.User_claims)
	if !ok {
		return nil, ErrInvalidClaims
	}
	//if the struct extracted successfully.
	return user, nil
}

func (handler *Todo_handler) Create_todo_handler(c *gin.Context) {
	//get the data from json body.
	var data model.Create_todo_request
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid json"})
		return
	}
	//get the userId to attach it with its todo.
	user, err := Get_user_claims(c)
	if err != nil {
		c.JSON(500, gin.H{"error": "Server error"})
		return
	}
	//create new todo with the input data.
	todo_ptr, err := handler.service.Create_todo(data, user)
	//check if something went wrong.
	switch {
	case err == nil:
		c.JSON(200, *todo_ptr)
	case errors.Is(err, service.ErrEmptyTitle):
		c.JSON(400, gin.H{"error": "Title cannot be empty"})
	case errors.Is(err, service.ErrEmptyCategory):
		c.JSON(400, gin.H{"error": "Category cannot be empty"})
	case errors.Is(err, service.ErrInvalidPriority):
		c.JSON(400, gin.H{"error": "Priority must be one of (low, medium, high)"})
	case errors.Is(err, service.ErrInvalidDueDate):
		c.JSON(400, gin.H{"error": "Unsupported due date format"})
	case errors.Is(err, service.ErrTimeInPast):
		c.JSON(400, gin.H{"error": "Due date must be in the future"})
	default:
		c.JSON(500, gin.H{"error": "Server error"})
	}
}

func (handler *Todo_handler) Update_todo_by_id_handler(c *gin.Context) {
	//first take the id param.
	id := c.Param("id")
	//then take the json body with struct.
	var data model.Update_todo_request
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid json"})
		return
	}
	//get the user claims to restrict modification.
	user, err := Get_user_claims(c)
	if err != nil {
		c.JSON(500, gin.H{"error": "Server error"})
		return
	}
	//finally, update required item with update function.
	todo_ptr, err := handler.service.Update_todo_by_id(id, data, user)
	//check if something went wrong.
	switch {
	case err == nil:
		c.JSON(200, *todo_ptr)
	case errors.Is(err, service.ErrInvalidId):
		c.JSON(400, gin.H{"error": "Invalid todo id"})
	case errors.Is(err, service.ErrEmptyTitle):
		c.JSON(400, gin.H{"error": "Title cannot be empty"})
	case errors.Is(err, service.ErrEmptyCategory):
		c.JSON(400, gin.H{"error": "Category cannot be empty"})
	case errors.Is(err, service.ErrInvalidPriority):
		c.JSON(400, gin.H{"error": "Priority must be one of (low, medium, high)"})
	case errors.Is(err, service.ErrNotFound):
		c.JSON(404, gin.H{"error": "Todo not found"})
	case errors.Is(err, service.ErrInvalidDueDate):
		c.JSON(400, gin.H{"error": "Unsupported due date format"})
	case errors.Is(err, service.ErrTimeInPast):
		c.JSON(400, gin.H{"error": "Due date must be in the future"})
	default:
		c.JSON(500, gin.H{"error": "Server error"})
	}
}

func (handler *Todo_handler) Update_todos_by_category_handler(c *gin.Context) {
	//first take the category param.
	category := c.Param("category")
	//get the completed status from json body.
	var data model.Update_todo_request
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid json"})
		return
	}
	//check if completed field not entered.
	if data.Completed == nil {
		c.JSON(400, gin.H{"error": "Missing new completed status"})
		return
	}
	//update todos of passed category.
	new_status := *(data.Completed)
	//get user calims to restrict modification to user's todos.
	user, err := Get_user_claims(c)
	if err != nil {
		c.JSON(500, gin.H{"error": "Server error"})
		return
	}
	err = handler.service.Update_todos_by_category(category, new_status, user)
	//check the result of process.
	switch {
	case err == nil:
		c.JSON(200, gin.H{"message": "Todos updated successfully"})
	case errors.Is(err, service.ErrEmptyCategory):
		c.JSON(400, gin.H{"error": "Category cannot be empty"})	
	default:
		c.JSON(500, gin.H{"error": "Server error"})
	}
}

func (handler *Todo_handler) Delete_todo_by_id_handler(c *gin.Context) {
	//get the id param.
	id := c.Param("id")
	//get the user claims to restrict deletion.
	user, err := Get_user_claims(c)
	if err != nil {
		c.JSON(500, gin.H{"error": "Server error"})
		return
	}
	//delete todo of that ID.
	err = handler.service.Delete_todo_by_id(id, user)
	//check the result of process.
	switch {
	case err == nil:
		c.JSON(200, gin.H{"message": "Todo deleted successfully"})
	case errors.Is(err, service.ErrInvalidId):
		c.JSON(400, gin.H{"error": "Invalid todo id"})
	case errors.Is(err, service.ErrNotFound):
		c.JSON(404, gin.H{"error": "Todo not found"})
	default:
		c.JSON(500, gin.H{"error": "Server error"})
	}
}

func (handler *Todo_handler) Delete_all_todos_handler(c *gin.Context) {
	//delete all todos in the system.
	err := handler.service.Delete_all_todos()
	//check the result of process.
	switch {
	case err == nil:
		c.JSON(200, gin.H{"message": "All todos deleted successfully"})
	default:
		c.JSON(500, gin.H{"error": "Server error"})
	}
}

func (handler *Todo_handler) Get_todo_by_id_handler(c *gin.Context) {
	//get id parameter.
	id := c.Param("id")
	//get the user claims to restrict search to user's todos only if the role is user.
	user, err := Get_user_claims(c)
	if err != nil {
		c.JSON(500, gin.H{"error": "Server error"})
		return
	}
	//get the todo of that id.
	todo_ptr, err := handler.service.Get_todo_by_id(id, user)
	//check if something went wrong.
	switch {
	case err == nil:
		c.JSON(200, *todo_ptr)
	case errors.Is(err, service.ErrInvalidId):
		c.JSON(400, gin.H{"error": "Invalid todo id"})
	case errors.Is(err, service.ErrNotFound):
		c.JSON(404, gin.H{"error": "Todo not found"})
	default:
		c.JSON(500, gin.H{"error": "Server error"})
	}
}

func (handler *Todo_handler) Get_todos_by_category_handler(c *gin.Context) {
	//get category param
	category := c.Param("category")
	//get the user claims to restrict todos to user's todos only if the role is user.
	user, err := Get_user_claims(c)
	if err != nil {
		c.JSON(500, gin.H{"error": "Server error"})
		return
	}
	//get all todos of that category.
	todos, err := handler.service.Get_todos_by_category(category, user)
	//check the result of process.
	switch {
	case err == nil:
		c.JSON(200, todos)
	case errors.Is(err, service.ErrEmptyCategory):
		c.JSON(400, gin.H{"error": "Category cannot be empty"})		
	default:
		c.JSON(500, gin.H{"error": "Server error"})
	}
}

func (handler *Todo_handler) Get_todos_by_status_handler(c *gin.Context) {
	//get status param.
	status := c.Param("status")
	//get the user claims to restrict todos to user's todos only if the role is user.
	user, err := Get_user_claims(c)
	if err != nil {
		c.JSON(500, gin.H{"error": "Server error"})
		return
	}
	//get todos of that status.
	todos, err := handler.service.Get_todos_by_status(status, user)
	//check the result of process.
	switch {
	case err == nil:
		c.JSON(200, todos)
	case errors.Is(err, service.ErrInvalidStatus):
		c.JSON(400, gin.H{"error": "Invalid todo status"})
	default:
		c.JSON(500, gin.H{"error": "Server error"})
	}
}

func (handler *Todo_handler) Get_todos_by_title_handler(c *gin.Context) {
	//get the query param.
	title, exist := c.GetQuery("q")
	//check if it exist or not.
	if !exist {
		c.JSON(400, gin.H{"error": "Missing query param for searching"})
		return
	}
	//get the user claims to restrict todos to user's todos only if the role is user.
	user, err := Get_user_claims(c)
	if err != nil {
		c.JSON(500, gin.H{"error": "Server error"})
		return
	}
	//get todos that contain this query param.
	todos, err := handler.service.Get_todos_by_title(title, user)
	//check the result of process.
	switch {
	case err == nil:
		c.JSON(200, todos)
	default:
		c.JSON(500, gin.H{"error": "Server error"})
	}
}

func (handler *Todo_handler) Get_all_todos_handler(c *gin.Context) {
	//get the user claims to restrict todos to user's todos only if the role is user.
	user, err := Get_user_claims(c)
	if err != nil {
		c.JSON(500, gin.H{"error": "Server error"})
		return
	}
	//Get all todos in system.
	todos, err := handler.service.Get_all_todos(user)
	//check the result of process.
	switch {
	case err == nil:
		c.JSON(200, todos)
	case errors.Is(err, service.ErrEmptyTitle):
		c.JSON(400, gin.H{"error": "Title cannot be empty"})	
	default:
		c.JSON(500, gin.H{"error": "Server error"})
	}
}
