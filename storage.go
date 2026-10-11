package main

import (
	"strings"
	"slices"
	"errors"
)

//create a global Id variable.
var id int = 0
//create a global empty slice for storage using make to avoid null behaviour.
var todo_list = make([]*Todo, 0)

func get_id() int{
	id++
	return id
}

type Todo struct {
	//Entities must be of uppercase to be accessable to json libraries.
	//add JSON tags to ensure consistent lowercase in requests and responses.
	Id        int		`json:"id"`
	Title     string 	`json:"title"`
	Completed bool		`json:"completed"`
}

func create_todo(title string, status bool) *Todo {
	//first we check the title.
	if strings.TrimSpace(title) == "" {
		return nil		//if not valid return nil.
	}
	//if valid create a new item.
	new_id := get_id()
	new_todo := Todo{new_id, title, status}
	todo_list = append(todo_list, &new_todo)
	//return the address of created item.
	return &new_todo
}

func get_all() []*Todo{
	//return the todo_list.
	return todo_list
}

func get_todo(id int) *Todo{
	//iterate until finding the id.
	for _, value := range(todo_list){
		if value.Id == id{
			return value
		}
	}
	return nil			//if not found return nil.
}

func update_todo(id int, new_title string, new_status bool) (*Todo, error){
	//first we check the new title.
	if strings.TrimSpace(new_title) == ""{
		err := errors.New("empty")
		return nil, err
	}
	//if valid search for id and update item.
	for _, value := range(todo_list){
		if value.Id == id{
			value.Title = new_title
			value.Completed = new_status
			return value, nil
		}
	}
	//if not found
	err := errors.New("not found")
	return nil, err
}

func delete_todo(id int) error{
	//iterate to search for id.
	for index, value := range(todo_list){
		//if id found, delete item.
		if value.Id == id{
			todo_list = slices.Delete(todo_list, index, index + 1)
			return nil
		}
	}
	//if id not found.
	err := errors.New("not found")
	return err
}