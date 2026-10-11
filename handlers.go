package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type json_data struct {
	//this struct to ensure the input json has valid data.
	Title     string `json:"title"`
	Completed bool   `json:"completed"` //if it was required false will be refused as it is default value.
}

func get_all_todos_handler(c *gin.Context) {
	c.JSON(200, get_all())
}

func get_todo_handler(c *gin.Context) {
	//get the id from context.
	var input_id = c.Param("id")
	//convert it to integer.
	id, err := strconv.Atoi(input_id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Id must be an integer"})
		return
	}
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
	data := json_data{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid json"})
		return
	}
	//call create function to make new todo.
	todo_ptr := create_todo(data.Title, data.Completed)
	//check if element created or not.
	if todo_ptr != nil {
		c.JSON(200, *todo_ptr)
	} else {
		c.JSON(400, gin.H{"error": "Title cannot be empty"})
	}
}

func update_todo_handler(c *gin.Context) {
	//get id from parameter
	input_id := c.Param("id")
	//convert to integer.
	id, err := strconv.Atoi(input_id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Id must be an integer"})
		return
	}

	//get the new data from json body.
	data := json_data{}
	err = c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid json"})
		return
	}

	//call update function to update todo.
	todo_ptr, err := update_todo(id, data.Title, data.Completed)
	//check if element updated or not.
	if todo_ptr != nil {
		c.JSON(200, *todo_ptr)
	} else { //if we have error value.
		var err_val = err.Error()
		//check the error to return the status code.
		if err_val == "empty" {
			c.JSON(400, gin.H{"error": "Title cannot be empty"})
		} else {
			c.JSON(404, gin.H{"error": "Todo not found"})
		}
	}
}

func delete_todo_handler(c *gin.Context) {
	//get the id from context.
	var input_id = c.Param("id")
	id, err := strconv.Atoi(input_id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Id must be an integer"})
		return
	}
	err = delete_todo(id)
	if err == nil {
		c.JSON(200, gin.H{"message": "Todo deleted"})
	} else {
		c.JSON(404, gin.H{"error": "Todo not found"})
	}
}
