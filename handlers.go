package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func get_all_todos_handler(c *gin.Context) {
	c.JSON(200, get_all())
}

func get_todo_handler(c *gin.Context) {
	//get the id from context.
	var input_id = c.Param("id")
	//convert it to integer.
	id, _ := strconv.Atoi(input_id)
	todo_ptr := get_todo(id)
	//if the todo found.
	if todo_ptr != nil {
		c.JSON(200, *todo_ptr)
	} else {
		c.JSON(404, gin.H{"error": "Todo not found"})
	}
}

func create_todo_handler(c *gin.Context) {
	//get the data from json body.
	data := make(map[string]any)
	c.ShouldBindJSON(&data)
	//call create function to make new todo.
	todo_ptr := create_todo(data["title"].(string), data["completed"].(bool))
	//check if element created or not.
	if todo_ptr != nil {
		c.JSON(200, *todo_ptr)
	} else {
		c.JSON(400, gin.H{"error": "Title cannot be empty"})
	}
}

func update_todo_handler(c *gin.Context) {
	//get the data from json body.
	data := make(map[string]any)
	c.ShouldBindJSON(&data)
	//call update function to update todo.
	todo_ptr, err := update_todo(data["id"].(int), data["title"].(string), data["completed"].(bool))
	//check if element updated or not.
	if todo_ptr != nil {
		c.JSON(200, *todo_ptr)
	} else { //if we have error value.
		var err_val = err.Error()
		//check the error to return the status code.
		if err_val == "empty" {
			c.JSON(400, gin.H{"error": "Title cannot be empty"})
		} else {
			c.JSON(404, gin.H{"error": "Title not found"})
		}
	}
}

func delete_todo_handler(c *gin.Context) {
	//get the id from context.
	var input_id = c.Param("id")
	id, _ := strconv.Atoi(input_id)
	err := delete_todo(id)
	if err == nil {
		c.JSON(200, gin.H{"message": "Todo deleted"})
	} else {
		c.JSON(404, gin.H{"error": "Todo not found"})
	}
}
