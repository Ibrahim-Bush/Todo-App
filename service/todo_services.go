package service

import (
	"Permanent-Storage-Todo/model"
	"Permanent-Storage-Todo/repository"
	"errors"
	"strconv"
	"strings"
	"time"
)

var (
	ErrNotFound        = errors.New("Todo not found")
	ErrEmptyTitle      = errors.New("Title cannot be empty")
	ErrEmptyCategory   = errors.New("Category cannot be empty")
	ErrInvalidPriority = errors.New("priority must be low, meduim or high")
	ErrInvalidDueDate  = errors.New("Invalid Due date")
	ErrTimeInPast      = errors.New("Time must be in the future")
	ErrInvalidId       = errors.New("Invalid Todo ID")
	ErrInvalidStatus   = errors.New("Invalid Todo Status")
	ErrServerError     = errors.New("Datebase Error")
)

type Todo_service interface {
	Modify_todo_service
	Get_todo_service
}

type Modify_todo_service interface {
	Create_todo(input_struct model.Create_todo_request, user *model.User_claims) (*model.Todo, error)
	Update_todo_by_id(input_id string, input_struct model.Update_todo_request, user *model.User_claims) (*model.Todo, error)
	Update_todos_by_category(input_category string, status bool, user *model.User_claims) error
	Delete_todo_by_id(input_id string, user *model.User_claims) error
	Delete_all_todos() error
}

type Get_todo_service interface {
	Get_all_todos(user *model.User_claims) ([]model.Todo, error)
	Get_todo_by_id(input_id string, user *model.User_claims) (*model.Todo, error)
	Get_todos_by_category(input_category string, user *model.User_claims) ([]model.Todo, error)
	Get_todos_by_status(input_status string, user *model.User_claims) ([]model.Todo, error)
	Get_todos_by_title(input_title string, user *model.User_claims) ([]model.Todo, error)
}

type todo_service struct {
	repository repository.Todo_repo
}

func Init_service(repo repository.Todo_repo) Todo_service {
	var new_service = todo_service{repository: repo}
	return &new_service
}

func (service *todo_service) Create_todo(input_struct model.Create_todo_request, user *model.User_claims) (*model.Todo, error) {
	//define a Todo object for new todo.
	var todo = model.Todo{
		Title:     input_struct.Title,
		Category:  input_struct.Category,
		Priority:  input_struct.Priority,
		Completed: input_struct.Completed,
		UserID:    user.UserID,
	}
	//Normalize all fields and validate them.
	err := Normalize_todo_fields(&todo)
	if err != nil {
		return nil, err
	}
	//check the dueDate if provided.
	due_date, err := Parse_standard_date(input_struct.DueDate)
	if err != nil {
		return nil, err
	}
	//check if the field is empty.
	if due_date.IsZero() != true {
		todo.DueDate = &due_date
	}
	//finally, create a record for new todo.
	err = service.repository.Create_record(&todo)
	if err != nil {
		return nil, ErrServerError
	}
	return &todo, nil
}

func Parse_standard_date(input_date string) (time.Time, error) {
	//first, clean empty spaces at the start and end.
	clean_date := strings.TrimSpace(input_date)
	if clean_date == "" {
		return time.Time{}, nil
	}
	//convert string to standard date.
	standard_date, err := time.Parse(time.RFC3339, clean_date)
	if err != nil {
		//try to convert it to date only without time.
		standard_date, err = time.Parse("2006-01-02", clean_date)
		if err != nil {
			return time.Time{}, ErrInvalidDueDate
		}
	}
	//convert date to UTC.
	standard_date = standard_date.UTC()
	//check if the time is in the past.
	if standard_date.Before(time.Now().UTC()) {
		return time.Time{}, ErrTimeInPast
	}
	//if the time was valid return standard time.
	return standard_date, nil
}

func Normalize_todo_fields(todo *model.Todo) error {
	//validate title field.
	todo.Title = strings.TrimSpace(todo.Title)
	if todo.Title == "" {
		return ErrEmptyTitle
	}
	//validate category field.
	todo.Category = strings.ToLower(strings.TrimSpace(todo.Category))
	if todo.Category == "" {
		return ErrEmptyCategory
	}
	//validate priority field.
	todo.Priority = strings.ToLower(strings.TrimSpace(todo.Priority))
	if todo.Priority != "low" && todo.Priority != "medium" && todo.Priority != "high" {
		return ErrInvalidPriority
	}
	//validate completed field.
	if todo.Completed == true {
		current_time := time.Now().UTC()
		todo.CompletedAt = &current_time
	}
	//After checking fields return nil after success.
	return nil
}

func (service *todo_service) Update_todo_by_id(input_id string, input_struct model.Update_todo_request, user *model.User_claims) (*model.Todo, error) {
	//first convert id to integer.
	id, err := strconv.Atoi(input_id)
	if err != nil {
		return nil, ErrInvalidId
	}
	//validate fields before update.
	err = Normalize_updates(&input_struct)
	if err != nil {
		return nil, err
	}
	//then get the required todo to update it.
	todo_ptr, err := service.repository.Get_record_by_id(id, user)
	if err != nil {
		//check if record not found or server fails.
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrServerError
	}
	//apply updates to the original todo.
	err = Apply_changes(&input_struct, todo_ptr)
	//check if something was invalid.
	if err != nil {
		return nil, err
	}
	//save the todo with updates.
	err = service.repository.Update_record_by_id(todo_ptr)
	if err != nil {
		return nil, ErrServerError
	}
	//after successfull update return the modified todo.
	return todo_ptr, nil
}

func Normalize_updates(requestDTO *model.Update_todo_request) error {
	//check the title field.
	if requestDTO.Title != nil {
		//normalize new title.
		clean_title := strings.TrimSpace(*(requestDTO.Title))
		if clean_title == "" {
			return ErrEmptyTitle
		}
		*(requestDTO.Title) = clean_title
	}
	//check the category field.
	if requestDTO.Category != nil {
		//normalize category.
		clean_category := strings.ToLower(strings.TrimSpace(*(requestDTO.Category)))
		if clean_category == "" {
			return ErrEmptyCategory
		}
		*(requestDTO.Category) = clean_category
	}
	//check the Priority field.
	if requestDTO.Priority != nil {
		//normalize Priority.
		clean_priority := strings.ToLower(strings.TrimSpace(*(requestDTO.Priority)))
		if clean_priority != "low" && clean_priority != "medium" && clean_priority != "high" {
			return ErrInvalidPriority
		}
		*(requestDTO.Priority) = clean_priority
	}
	//after successfull validation, retrun nil.
	return nil
}

func Apply_changes(updates_struct *model.Update_todo_request, todo_ptr *model.Todo) error {
	//check and update the DueDate first.
	if updates_struct.DueDate != nil {
		clean_date := strings.TrimSpace(*(updates_struct.DueDate))
		standard_date, err := Parse_standard_date(clean_date)
		if err != nil {
			return err
		}
		//if the date is not empty.
		if standard_date.IsZero() != true {
			todo_ptr.DueDate = &standard_date
		}
	}
	//update title.
	if updates_struct.Title != nil {
		todo_ptr.Title = *(updates_struct.Title)
	}
	//update category.
	if updates_struct.Category != nil {
		todo_ptr.Category = *(updates_struct.Category)
	}
	//update Priority.
	if updates_struct.Priority != nil {
		todo_ptr.Priority = *(updates_struct.Priority)
	}
	//update completed status.
	if updates_struct.Completed != nil {
		//if updating to true update completedAt.
		if *(updates_struct.Completed) == true {
			todo_ptr.Completed = true
			current_time := time.Now().UTC()
			todo_ptr.CompletedAt = &current_time
		} else {
			todo_ptr.Completed = false
			todo_ptr.CompletedAt = nil
		}
	}
	//after update all fields return nil.
	return nil
}

func (service *todo_service) Update_todos_by_category(input_category string, status bool, user *model.User_claims) error {
	//normalize the input category.
	category := strings.ToLower(strings.TrimSpace(input_category))
	//create a struct for updates.
	var todo = model.Todo{
		Completed: status,
	}	
	//check the new status to update completedAt field.
	if status == true {
		current_time := time.Now().UTC()
		todo.CompletedAt = &current_time
	} else{
		todo.CompletedAt = nil
	}
	//update all todos of that category.
	err := service.repository.Update_records_by_category(category, &todo, user)
	if err != nil {
		return ErrServerError
	}
	return nil
}

func (service *todo_service) Delete_todo_by_id(input_id string, user *model.User_claims) error {
	//convert id to integer.
	id, err := strconv.Atoi(input_id)
	if err != nil {
		return ErrInvalidId
	}
	//delete todo of that id.
	err = service.repository.Delete_record_by_id(id, user)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, repository.ErrNotFound):
		return ErrNotFound
	default:
		return err
	}
}

func (service *todo_service) Delete_all_todos() error {
	//delete all todos stored in the repositoy.
	err := service.repository.Delete_all_records()
	if err != nil {
		return ErrServerError
	}
	return nil
}

func (service *todo_service) Get_all_todos(user *model.User_claims) ([]model.Todo, error) {
	//get all records from repository.
	todos, err := service.repository.Get_all_records(user)
	//check if something went wrong.
	if err != nil {
		return nil, ErrServerError
	}
	return todos, nil
}

func (service *todo_service) Get_todo_by_id(input_id string, user *model.User_claims) (*model.Todo, error) {
	//convert id to integer.
	id, err := strconv.Atoi(input_id)
	if err != nil {
		return nil, ErrInvalidId
	}
	//get todo from repository.
	todo_ptr, err := service.repository.Get_record_by_id(id, user)
	//check the result.
	switch {
	case err == nil:
		return todo_ptr, nil
	case errors.Is(err, repository.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, ErrServerError
	}
}

func (service *todo_service) Get_todos_by_category(input_category string, user *model.User_claims) ([]model.Todo, error) {
	//normalize input category.
	category := strings.ToLower(strings.TrimSpace(input_category))
	//get todos of this category from repository.
	todos, err := service.repository.Get_records_by_category(category, user)
	if err != nil {
		return nil, ErrServerError
	}
	return todos, nil
}

func (service *todo_service) Get_todos_by_status(input_status string, user *model.User_claims) ([]model.Todo, error) {
	//clean input from white sapces.
	clean_status := strings.TrimSpace(input_status)
	//convert status to boolean.
	status, err := strconv.ParseBool(clean_status)
	if err != nil {
		return nil, ErrInvalidStatus
	}
	//get todos of that status from repository.
	todos, err := service.repository.Get_records_by_status(status, user)
	if err != nil {
		return nil, ErrServerError
	}
	return todos, nil
}

func (service *todo_service) Get_todos_by_title(input_title string, user *model.User_claims) ([]model.Todo, error) {
	//clean white spaces at the start and end.
	title := strings.TrimSpace(input_title)
	//get todos that match this title from repository.
	todos, err := service.repository.Get_records_by_title(title, user)
	if err != nil {
		return nil, ErrServerError
	}
	return todos, nil
}
