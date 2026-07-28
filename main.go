package main

import "github.com/gin-gonic/gin"

func main(){

	//create a gin router with default middleware.
	router := gin.Default()

	//route requests to its appropriate handlers.
	//the gin framework will automatically call the passed function with the Context object.
	router.GET("/todos", get_all_todos_handler)
	router.GET("/todos/:id", get_todo_handler)
	router.POST("/todos", create_todo_handler)
	router.PUT("/todos/:id", update_todo_handler)
	router.DELETE("/todos/:id", delete_todo_handler)

	//starting API server at local port ":8080"
	router.Run(":8080")
}