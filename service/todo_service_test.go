package service

import (
	"Permanent-Storage-Todo/mocks"
	"Permanent-Storage-Todo/model"
	"Permanent-Storage-Todo/repository"
	"errors"
	"testing"
	"time"
)

func Test_create_todo(t *testing.T) {
	//define a struct for usecase.
	type usecase struct {
		name          string
		input_todo    model.Create_todo_request
		input_user    *model.User_claims
		repo_todo     *model.Todo
		repo_err      error
		expected_todo *model.Todo
		expected_err  error
		call_repo     bool
	}
	//create a valid time for success state.
	Time := time.Now().UTC()
	//create a table of usecase.
	tests := []usecase{
		{
			name:          "Error: Empty title",
			input_todo:    model.Create_todo_request{Title: "	"},
			input_user:    &model.User_claims{UserID: 1},
			repo_todo:     nil,
			repo_err:      nil,
			expected_todo: nil,
			expected_err:  ErrEmptyTitle,
			call_repo:     false,
		}, {
			name:          "Error: Empty category",
			input_todo:    model.Create_todo_request{Title: "  task 1 ", Category: "	"},
			input_user:    &model.User_claims{UserID: 1},
			repo_todo:     nil,
			repo_err:      nil,
			expected_todo: nil,
			expected_err:  ErrEmptyCategory,
			call_repo:     false,
		}, {
			name:          "Error: Empty category",
			input_todo:    model.Create_todo_request{Title: "  task 1 ", Category: "	"},
			input_user:    &model.User_claims{UserID: 1},
			repo_todo:     nil,
			repo_err:      nil,
			expected_todo: nil,
			expected_err:  ErrEmptyCategory,
			call_repo:     false,
		}, {
			name:          "Error: Invalid priority",
			input_todo:    model.Create_todo_request{Title: "  task 1 ", Category: " Work ", Priority: " urgent "},
			input_user:    &model.User_claims{UserID: 1},
			repo_todo:     nil,
			repo_err:      nil,
			expected_todo: nil,
			expected_err:  ErrInvalidPriority,
			call_repo:     false,
		}, {
			name:          "Error: Invalid dueDate",
			input_todo:    model.Create_todo_request{Title: "  task 1 ", Category: " Work ", Priority: " high ", DueDate: "  1/2/2000  "},
			input_user:    &model.User_claims{UserID: 1},
			repo_todo:     nil,
			repo_err:      nil,
			expected_todo: nil,
			expected_err:  ErrInvalidDueDate,
			call_repo:     false,
		}, {
			name:          "Error: Time in past",
			input_todo:    model.Create_todo_request{Title: "  task 1 ", Category: " Work ", Priority: " high ", DueDate: "  2010-01-05  "},
			input_user:    &model.User_claims{UserID: 1},
			repo_todo:     nil,
			repo_err:      nil,
			expected_todo: nil,
			expected_err:  ErrTimeInPast,
			call_repo:     false,
		}, {
			name:          "Error: Server error",
			input_todo:    model.Create_todo_request{Title: "  task 1 ", Category: " WOrK ", Priority: "  HIgh "},
			input_user:    &model.User_claims{UserID: 1},
			repo_todo:     &model.Todo{Title: "task 1", Category: "work", Priority: "high", Completed: false, DueDate: nil, CompletedAt: nil, UserID: 1},
			repo_err:      errors.New("Database connection failure"),
			expected_todo: nil,
			expected_err:  ErrServerError,
			call_repo:     true,
		}, {
			name:          "Success: Todo created",
			input_todo:    model.Create_todo_request{Title: "  task 1 ", Category: " WOrK ", Priority: "  HIgh ", Completed: true, DueDate: "  2200-01-05  "},
			input_user:    &model.User_claims{UserID: 1},
			repo_todo:     &model.Todo{Title: "task 1", Category: "work", Priority: "high", Completed: true, DueDate: &Time, CompletedAt: &Time, UserID: 1},
			repo_err:      nil,
			expected_todo: &model.Todo{Title: "task 1", Category: "work", Priority: "high", Completed: true, DueDate: &Time, CompletedAt: &Time, UserID: 1},
			expected_err:  nil,
			call_repo:     true,
		},
	}
	//run tests one by one.
	for _, usecase := range tests {
		t.Run(usecase.name, func(t *testing.T) {
			//create a new mock.
			mock_repo := &mocks.Mock_todo_repo{
				Wanted_err: usecase.repo_err,
			}
			//pass mock to the service.
			service := Init_service(mock_repo)
			//call the service method.
			todo, err := service.Create_todo(usecase.input_todo, usecase.input_user)
			//check the repo call: first check the call state.
			if mock_repo.Was_called != usecase.call_repo {
				t.Errorf("Behaviour: Expected repository call %t got %t", usecase.call_repo, mock_repo.Was_called)
			}
			//compare the struct passed to repo with the expected struct.
			compare_todo_structs(t, "Repo input", mock_repo.Received_todo, usecase.repo_todo)
			//check the response of the service.
			if err != usecase.expected_err {
				t.Errorf("Response: Expected error %v got %v", usecase.expected_err, err)
			}
			//compare the struct in response with expected struct.
			compare_todo_structs(t, "Response", todo, usecase.expected_todo)
		})
	}
}

func compare_todo_structs(t *testing.T, stage string, received *model.Todo, expected *model.Todo) {
	//first compare pointers.
	if received == nil || expected == nil {
		if received == nil && expected != nil {
			t.Errorf("Boundary (%s): expected struct %+v got nil pointer", stage, expected)
		} else if received != nil && expected == nil {
			t.Errorf("Boundary (%s): expected nil struct got %+v", stage, received)
		}
	} else {
		//compare the todo fields.
		if received.Title != expected.Title {
			t.Errorf("Logic (%s): Expected title %s got %s", stage, expected.Title, received.Title)
		}
		if received.Category != expected.Category {
			t.Errorf("Logic (%s): Expected category %s got %s", stage, expected.Category, received.Category)
		}
		if received.Priority != expected.Priority {
			t.Errorf("Logic (%s): Expected priority %s got %s", stage, expected.Priority, received.Priority)
		}
		if received.Completed != expected.Completed {
			t.Errorf("Logic (%s): Expected Completed status %t got %t", stage, expected.Completed, received.Completed)
		}
		if expected.DueDate == nil && received.DueDate != expected.DueDate {
			t.Errorf("Logic (%s): Expected DueDate %v got %v", stage, expected.DueDate, received.DueDate)
		} else if expected.DueDate != nil && received.DueDate == nil {
			t.Errorf("Logic (%s): Expected not nil DueDate got %v", expected.DueDate, received.DueDate)
		}
		if expected.CompletedAt == nil && received.CompletedAt != expected.CompletedAt {
			t.Errorf("Logic (%s): Expected CompletedAt %v got %v", stage, expected.CompletedAt, received.CompletedAt)
		} else if expected.CompletedAt != nil && received.CompletedAt == nil {
			t.Errorf("Logic (%s): Expected not nil CompletedAt got %v", expected.CompletedAt, received.CompletedAt)
		}
		if received.UserID != expected.UserID {
			t.Errorf("Logic (%s): Expected UserID %d got %d", stage, expected.UserID, received.UserID)
		}
	}
}

func Test_update_todo_by_id(t *testing.T) {
	//define a struct for usecase.
	type usecase struct {
		name           string
		input_id       string
		input_request  model.Update_todo_request
		input_user     *model.User_claims
		repo_id        int
		repo_todo      *model.Todo
		repo_err       error
		get_by_id_err  error
		expected_todo  *model.Todo
		expected_err   error
		update_call    bool
		get_by_id_call bool
	}
	//create a valid time for time update.
	Time := time.Now().UTC()
	//create values to use their places in update pointers.
	title := "	finish task 1	"
	category := "		WorK	"
	empty := ""
	completed := true
	not_completed := false
	priority := "		LOW	"
	due_date := "2200-10-10"
	//create a table for usecase.
	tests := []usecase{
		{
			name:           "Error: Invalid id",
			input_id:       "number",
			input_request:  model.Update_todo_request{},
			expected_todo:  nil,
			expected_err:   ErrInvalidId,
			update_call:    false,
			get_by_id_call: false,
		}, {
			name:           "Error: Empty title",
			input_id:       "1",
			input_request:  model.Update_todo_request{Title: &empty},
			expected_todo:  nil,
			expected_err:   ErrEmptyTitle,
			update_call:    false,
			get_by_id_call: false,
		}, {
			name:           "Error: Empty category",
			input_id:       "1",
			input_request:  model.Update_todo_request{Title: &title, Category: &empty},
			expected_todo:  nil,
			expected_err:   ErrEmptyCategory,
			update_call:    false,
			get_by_id_call: false,
		}, {
			name:           "Error: Todo not found",
			input_id:       "1",
			input_request:  model.Update_todo_request{Title: &title, Category: &category, Priority: &priority},
			input_user:     &model.User_claims{Role: "user"},
			repo_id:        1,
			get_by_id_err:  repository.ErrNotFound,
			repo_todo:      nil,
			expected_todo:  nil,
			expected_err:   ErrNotFound,
			update_call:    false,
			get_by_id_call: true,
		}, {
			name:           "Error: Invalid dueDate",
			input_id:       "1",
			input_request:  model.Update_todo_request{Title: &title, Category: &category, Priority: &priority, DueDate: &title},
			input_user:     &model.User_claims{Role: "user"},
			repo_id:        1,
			get_by_id_err:  nil,
			repo_todo:      &model.Todo{Title: "pass exam", Completed: false, Category: "study", Priority: "high"},
			expected_todo:  nil,
			expected_err:   ErrInvalidDueDate,
			update_call:    false,
			get_by_id_call: true,
		}, {
			name:           "Error: Server error with update todo",
			input_id:       "1",
			input_request:  model.Update_todo_request{Title: &title, Completed: &completed, Category: &category, Priority: &priority, DueDate: &due_date},
			input_user:     &model.User_claims{Role: "user"},
			repo_id:        1,
			repo_err:       errors.New("Database connection failure"),
			get_by_id_err:  nil,
			repo_todo:      &model.Todo{Title: "pass exam", Completed: false, Category: "study", Priority: "high"},
			expected_todo:  nil,
			expected_err:   ErrServerError,
			update_call:    true,
			get_by_id_call: true,
		}, {
			name:           "Error: Successful update",
			input_id:       "1",
			input_request:  model.Update_todo_request{Title: &title, Completed: &completed, Category: &category, Priority: &priority, DueDate: &due_date},
			input_user:     &model.User_claims{Role: "user"},
			repo_id:        1,
			repo_err:       nil,
			get_by_id_err:  nil,
			repo_todo:      &model.Todo{Title: "pass exam", Completed: false, Category: "study", Priority: "high"},
			expected_todo:  &model.Todo{Title: title, Completed: completed, Category: category, Priority: priority, DueDate: &Time, CompletedAt: &Time},
			expected_err:   nil,
			update_call:    true,
			get_by_id_call: true,
		}, {
			name:           "Error: Successful update status with nil CompletedAt",
			input_id:       "1",
			input_request:  model.Update_todo_request{Completed: &not_completed},
			input_user:     &model.User_claims{Role: "user"},
			repo_id:        1,
			repo_err:       nil,
			get_by_id_err:  nil,
			repo_todo:      &model.Todo{Title: title, Completed: completed, Category: category, Priority: priority, DueDate: &Time, CompletedAt: &Time},
			expected_todo:  &model.Todo{Title: title, Completed: not_completed, Category: category, Priority: priority, DueDate: &Time, CompletedAt: nil},
			expected_err:   nil,
			update_call:    true,
			get_by_id_call: true,
		},
	}
	//run tests one by one.
	for _, usecase := range tests {
		t.Run(usecase.name, func(t *testing.T) {
			//set the mock.
			mock_repo := mocks.Mock_todo_repo{Wanted_err: usecase.repo_err, Wanted_todo: usecase.repo_todo, Get_by_id_err: usecase.get_by_id_err}
			//pass mock repo to initialize service.
			service := Init_service(&mock_repo)
			//call the method.
			todo, err := service.Update_todo_by_id(usecase.input_id, usecase.input_request, usecase.input_user)
			//check the response.
			if err != usecase.expected_err {
				t.Errorf("Response: expected error in update response %v got %v", usecase.expected_err, err)
				//check the todo in response.
				compare_todo_structs(t, "Response", todo, usecase.expected_todo)
				//check the repo call.
				if mock_repo.Get_by_id_call != usecase.get_by_id_call {
					t.Errorf("Behaviour: expected call to the get_todo_by_id %t got %t", usecase.get_by_id_call, mock_repo.Get_by_id_call)
				}
				if mock_repo.Was_called != usecase.update_call {
					t.Errorf("Behaviour: expected call to the update_todo %t got %t", usecase.update_call, mock_repo.Was_called)
				}
				//if get by id was called, check the passed id.
				if mock_repo.Get_by_id_call == true && usecase.get_by_id_call == true {
					if mock_repo.Received_id != usecase.repo_id {
						t.Errorf("Logic: expected id passed to get_todo_by_id method %d got %d", usecase.repo_id, mock_repo.Received_id)
					}
				}
				//in case of success, check the passed todo to the update method.
				if err == nil && usecase.expected_err == nil {
					compare_todo_structs(t, "Repo input", mock_repo.Received_todo, usecase.expected_todo)
				}
			}
		})
	}
}

func Test_update_todos_by_category(t *testing.T) {
	//define a struct for usecase.
	type usecase struct {
		name           string
		input_category string
		input_status   bool
		input_user     *model.User_claims
		repo_todo      *model.Todo
		repo_category  string
		repo_err       error
		expected_err   error
		repo_call      bool
	}
	//set a valid time to use it in completedAt check.
	Time := time.Now().UTC()
	//create a table of usecase.
	tests := []usecase{
		{
			name:           "Error: Empty category",
			input_category: "		",
			input_status:   true,
			input_user:     &model.User_claims{Role: "user"},
			expected_err:   ErrEmptyCategory,
			repo_call:      false,
		}, {
			name:           "Check completedAt with true status",
			input_category: "	WoRk	",
			input_status:   true,
			input_user:     &model.User_claims{Role: "user"},
			repo_todo:      &model.Todo{Completed: true, CompletedAt: &Time},
			repo_category:  "work",
			repo_err:       errors.New("Database connection failure"),
			expected_err:   ErrServerError,
			repo_call:      true,
		}, {
			name:           "change completedAt to nil false status",
			input_category: "	StUdY	",
			input_status:   false,
			input_user:     &model.User_claims{Role: "user"},
			repo_todo:      &model.Todo{Completed: false, CompletedAt: nil},
			repo_category:  "study",
			repo_err:       nil,
			expected_err:   nil,
			repo_call:      true,
		}, {
			name:           "successful change status to true",
			input_category: "	WoRk	",
			input_status:   true,
			input_user:     &model.User_claims{Role: "user"},
			repo_todo:      &model.Todo{Completed: true, CompletedAt: &Time},
			repo_category:  "work",
			repo_err:       nil,
			expected_err:   nil,
			repo_call:      true,
		},
	}
	//run tests one by one.
	for _, usecase := range tests {
		t.Run(usecase.name, func(t *testing.T) {
			//set the mock.
			mock_repo := mocks.Mock_todo_repo{Wanted_err: usecase.repo_err}
			//pass mock repo to initialize service.
			service := Init_service(&mock_repo)
			//call the method.
			err := service.Update_todos_by_category(usecase.input_category, usecase.input_status, usecase.input_user)
			//check the response.
			if err != usecase.expected_err {
				t.Errorf("Response: expected error in update response %v got %v", usecase.expected_err, err)
			}
			//check the repo call.
			if mock_repo.Was_called != usecase.repo_call {
				t.Errorf("Behaviour: expected repository method call %t got %t", usecase.repo_call, mock_repo.Was_called)
			}
			//check the repo input.
			if mock_repo.Received_category != usecase.repo_category {
				t.Errorf("Logic: expected category passed to repository method %s got %s", usecase.repo_category, mock_repo.Received_category)
			}
			//check the todo passed to repo.
			compare_todo_structs(t, "Repo input", mock_repo.Received_todo, usecase.repo_todo)
		})
	}
}

func Test_delete_todo_by_id(t *testing.T) {
	//define a struct for usecase.
	type usecase struct {
		name         string
		input_id     string
		input_user   *model.User_claims
		repo_id      int
		repo_err     error
		expected_err error
		repo_call    bool
	}
	//create a table for usecase.
	tests := []usecase{
		{
			name:         "Error: Invalid id",
			input_id:     "number",
			input_user:   &model.User_claims{Role: "user"},
			expected_err: ErrInvalidId,
			repo_call:    false,
		}, {
			name:         "Error: Server error",
			input_id:     "1",
			input_user:   &model.User_claims{Role: "user"},
			repo_id:      1,
			repo_err:     errors.New("Database connection failure"),
			expected_err: ErrServerError,
			repo_call:    true,
		}, {
			name:         "Error: todo not found",
			input_id:     "1",
			input_user:   &model.User_claims{Role: "user"},
			repo_id:      1,
			repo_err:     repository.ErrNotFound,
			expected_err: ErrNotFound,
			repo_call:    true,
		}, {
			name:         "Sueccess: todo deleted",
			input_id:     "1",
			input_user:   &model.User_claims{Role: "user"},
			repo_id:      1,
			repo_err:     nil,
			expected_err: nil,
			repo_call:    true,
		},
	}
	//run tests one by one.
	for _, usecase := range tests {
		t.Run(usecase.name, func(t *testing.T) {
			//set the mock.
			mock_repo := mocks.Mock_todo_repo{Wanted_err: usecase.repo_err}
			//pass mock repo to initialize service.
			service := Init_service(&mock_repo)
			//call the method.
			err := service.Delete_todo_by_id(usecase.input_id, usecase.input_user)
			//check the response.
			if err != usecase.expected_err {
				t.Errorf("Response: expected error in response %v got %v", usecase.expected_err, err)
			}
			//check the repo call.
			if mock_repo.Was_called != usecase.repo_call {
				t.Errorf("Behaviour: expected repository method call %t got %t", usecase.repo_call, mock_repo.Was_called)
			}
			//check the repo input.
			if mock_repo.Received_id != usecase.repo_id {
				t.Errorf("Logic: expected Id passed to repository method %d got %d", usecase.repo_id, mock_repo.Received_id)
			}
		})
	}
}

func Test_delete_all_todos(t *testing.T) {
	//define a struct for usecase.
	type usecase struct {
		name         string
		repo_err     error
		expected_err error
		repo_call    bool
	}
	//set a table of usecase.
	tests := []usecase{
		{
			name:         "Error: Server error",
			repo_err:     errors.New("Database connection failure"),
			expected_err: ErrServerError,
			repo_call:    true,
		}, {
			name:         "Sucess: todos deleted",
			repo_err:     nil,
			expected_err: nil,
			repo_call:    true,
		},
	}
	//run tests one by one.
	for _, usecase := range tests {
		t.Run(usecase.name, func(t *testing.T) {
			//set the mock.
			mock_repo := mocks.Mock_todo_repo{Wanted_err: usecase.repo_err}
			//pass mock repo to initialize service.
			service := Init_service(&mock_repo)
			//call the method.
			err := service.Delete_all_todos()
			//check the response.
			if err != usecase.expected_err {
				t.Errorf("Response: expected error in response %v got %v", usecase.expected_err, err)
			}
			//check the repo call.
			if mock_repo.Was_called != usecase.repo_call {
				t.Errorf("Behaviour: expected repository method call %t got %t", usecase.repo_call, mock_repo.Was_called)
			}
		})
	}
}

func Test_get_all_todos(t *testing.T) {
	//define a struct for usecase.
	type usecase struct {
		name           string
		input_user     *model.User_claims
		repo_todos     []model.Todo
		repo_err       error
		expected_todos []model.Todo
		expected_err   error
		repo_call      bool
	}
	//set a table of usecase.
	tests := []usecase{
		{
			name:           "Error: Server error",
			input_user:     &model.User_claims{Role: "user"},
			repo_err:       errors.New("Database connection failure"),
			expected_todos: nil,
			expected_err:   ErrServerError,
			repo_call:      true,
		}, {
			name:           "Sucess: get todos",
			input_user:     &model.User_claims{Role: "user"},
			repo_err:       nil,
			repo_todos:     []model.Todo{{Title: "task 1"}},
			expected_todos: []model.Todo{{Title: "task 1"}},
			expected_err:   nil,
			repo_call:      true,
		},
	}
	//run tests one by one.
	for _, usecase := range tests {
		t.Run(usecase.name, func(t *testing.T) {
			//set the mock.
			mock_repo := mocks.Mock_todo_repo{Wanted_err: usecase.repo_err, Wanted_list: usecase.repo_todos}
			//pass mock repo to initialize service.
			service := Init_service(&mock_repo)
			//call the method.
			todos, err := service.Get_all_todos(usecase.input_user)
			//check the response.
			if err != usecase.expected_err {
				t.Errorf("Response: expected error in response %v got %v", usecase.expected_err, err)
			}
			//check the todos in response.
			if usecase.expected_todos == nil && todos != nil {
				t.Errorf("Response: expected nil slice got %+v todos", todos)
			} else if expected := len(usecase.expected_todos); expected != len(todos) {
				t.Errorf("Response: expected slice of length %d got slice of length %d", expected, len(todos))
			}
			//check the repo call.
			if mock_repo.Was_called != usecase.repo_call {
				t.Errorf("Behaviour: expected repository method call %t got %t", usecase.repo_call, mock_repo.Was_called)
			}
		})
	}
}

func Test_get_todo_by_id(t *testing.T) {
	//define a struct for usecase.
	type usecase struct {
		name          string
		input_id      string
		input_user    *model.User_claims
		repo_id       int
		repo_err      error
		repo_todo     *model.Todo
		expected_err  error
		expected_todo *model.Todo
		repo_call     bool
	}
	//create a table for usecase.
	tests := []usecase{
		{
			name:          "Error: Invalid id",
			input_id:      "number",
			input_user:    &model.User_claims{Role: "user"},
			expected_err:  ErrInvalidId,
			expected_todo: nil,
			repo_call:     false,
		}, {
			name:          "Error: Server error",
			input_id:      "1",
			input_user:    &model.User_claims{Role: "user"},
			repo_id:       1,
			repo_err:      errors.New("Database connection failure"),
			expected_err:  ErrServerError,
			expected_todo: nil,
			repo_call:     true,
		}, {
			name:          "Error: todo not found",
			input_id:      "1",
			input_user:    &model.User_claims{Role: "user"},
			repo_id:       1,
			repo_err:      repository.ErrNotFound,
			expected_err:  ErrNotFound,
			expected_todo: nil,
			repo_call:     true,
		}, {
			name:          "Sueccess: get todo",
			input_id:      "1",
			input_user:    &model.User_claims{Role: "user"},
			repo_id:       1,
			repo_err:      nil,
			repo_todo:     &model.Todo{Title: "finish task 1", Category: "work", Priority: "high"},
			expected_err:  nil,
			expected_todo: &model.Todo{Title: "finish task 1", Category: "work", Priority: "high"},
			repo_call:     true,
		},
	}
	//run tests one by one.
	for _, usecase := range tests {
		t.Run(usecase.name, func(t *testing.T) {
			//set the mock.
			mock_repo := mocks.Mock_todo_repo{Get_by_id_err: usecase.repo_err, Wanted_todo: usecase.repo_todo}
			//pass mock repo to initialize service.
			service := Init_service(&mock_repo)
			//call the method.
			todo, err := service.Get_todo_by_id(usecase.input_id, usecase.input_user)
			//check the response.
			if err != usecase.expected_err {
				t.Errorf("Response: expected error in response %v got %v", usecase.expected_err, err)
			}
			//check the todo in response.
			compare_todo_structs(t, "Response", todo, usecase.expected_todo)
			//check the repo call.
			if mock_repo.Get_by_id_call != usecase.repo_call {
				t.Errorf("Behaviour: expected repository method call %t got %t", usecase.repo_call, mock_repo.Was_called)
			}
			//check the repo input.
			if mock_repo.Received_id != usecase.repo_id {
				t.Errorf("Logic: expected Id passed to repository method %d got %d", usecase.repo_id, mock_repo.Received_id)
			}
		})
	}
}

func Test_get_todos_by_category(t *testing.T) {
	//define a struct for usecase.
	type usecase struct {
		name           string
		input_user     *model.User_claims
		input_category string
		repo_todos     []model.Todo
		repo_category  string
		repo_err       error
		expected_todos []model.Todo
		expected_err   error
		repo_call      bool
	}
	//set a table of usecase.
	tests := []usecase{
		{
			name:           "Error: Empty category",
			input_category: "		",
			input_user:     &model.User_claims{Role: "user"},
			expected_todos: nil,
			expected_err:   ErrEmptyCategory,
			repo_call:      false,
		}, {
			name:           "Error: Server error",
			input_user:     &model.User_claims{Role: "user"},
			input_category: "		wORk	",
			repo_err:       errors.New("Database connection failure"),
			repo_category:  "work",
			expected_todos: nil,
			expected_err:   ErrServerError,
			repo_call:      true,
		}, {
			name:           "Sucess: get todos",
			input_user:     &model.User_claims{Role: "user"},
			input_category: "		wORk	",
			repo_err:       nil,
			repo_todos:     []model.Todo{{Title: "task 1", Category: "work"}, {Title: "task 2", Category: "work"}},
			repo_category:  "work",
			expected_todos: []model.Todo{{Title: "task 1", Category: "work"}, {Title: "task 2", Category: "work"}},
			expected_err:   nil,
			repo_call:      true,
		},
	}
	//run tests one by one.
	for _, usecase := range tests {
		t.Run(usecase.name, func(t *testing.T) {
			//set the mock.
			mock_repo := mocks.Mock_todo_repo{Wanted_err: usecase.repo_err, Wanted_list: usecase.repo_todos}
			//pass mock repo to initialize service.
			service := Init_service(&mock_repo)
			//call the method.
			todos, err := service.Get_todos_by_category(usecase.input_category, usecase.input_user)
			//check the response.
			if err != usecase.expected_err {
				t.Errorf("Response: expected error in response %v got %v", usecase.expected_err, err)
			}
			//check the todos in response.
			if usecase.expected_todos == nil && todos != nil {
				t.Errorf("Response: expected nil slice got %+v todos", todos)
			} else if expected := len(usecase.expected_todos); expected != len(todos) {
				t.Errorf("Response: expected slice of length %d got slice of length %d", expected, len(todos))
			}
			//check the repo call.
			if mock_repo.Was_called != usecase.repo_call {
				t.Errorf("Behaviour: expected repository method call %t got %t", usecase.repo_call, mock_repo.Was_called)
			}
			if mock_repo.Received_category != usecase.repo_category {
				t.Errorf("Logic: expected category passed to repository method %s got %s", usecase.repo_category, mock_repo.Received_category)
			}
		})
	}
}

func Test_get_todos_by_status(t *testing.T) {
	//define a struct for usecase.
	type usecase struct {
		name           string
		input_user     *model.User_claims
		input_status   string
		repo_todos     []model.Todo
		repo_status    bool
		repo_err       error
		expected_todos []model.Todo
		expected_err   error
		repo_call      bool
	}
	//set a table of usecase.
	tests := []usecase{
		{
			name:           "Error: Invalid status",
			input_status:   "	hello	",
			input_user:     &model.User_claims{Role: "user"},
			expected_todos: nil,
			expected_err:   ErrInvalidStatus,
			repo_call:      false,
		}, {
			name:           "Error: Server error",
			input_user:     &model.User_claims{Role: "user"},
			input_status:   "		true	",
			repo_err:       errors.New("Database connection failure"),
			repo_status:    true,
			expected_todos: nil,
			expected_err:   ErrServerError,
			repo_call:      true,
		}, {
			name:           "Sucess: get todos",
			input_user:     &model.User_claims{Role: "user"},
			input_status:   "		false		",
			repo_err:       nil,
			repo_todos:     []model.Todo{{Title: "task 1", Completed: false}, {Title: "task 2", Completed: false}},
			repo_status:    false,
			expected_todos: []model.Todo{{Title: "task 1", Completed: false}, {Title: "task 2", Completed: false}},
			expected_err:   nil,
			repo_call:      true,
		},
	}
	//run tests one by one.
	for _, usecase := range tests {
		t.Run(usecase.name, func(t *testing.T) {
			//set the mock.
			mock_repo := mocks.Mock_todo_repo{Wanted_err: usecase.repo_err, Wanted_list: usecase.repo_todos}
			//pass mock repo to initialize service.
			service := Init_service(&mock_repo)
			//call the method.
			todos, err := service.Get_todos_by_status(usecase.input_status, usecase.input_user)
			//check the response.
			if err != usecase.expected_err {
				t.Errorf("Response: expected error in response %v got %v", usecase.expected_err, err)
			}
			//check the todos in response.
			if usecase.expected_todos == nil && todos != nil {
				t.Errorf("Response: expected nil slice got %+v todos", todos)
			} else if expected := len(usecase.expected_todos); expected != len(todos) {
				t.Errorf("Response: expected slice of length %d got slice of length %d", expected, len(todos))
			}
			//check the repo call.
			if mock_repo.Was_called != usecase.repo_call {
				t.Errorf("Behaviour: expected repository method call %t got %t", usecase.repo_call, mock_repo.Was_called)
			}
			if mock_repo.Received_status != usecase.repo_status {
				t.Errorf("Logic: expected status passed to repository method %t got %t", usecase.repo_status, mock_repo.Received_status)
			}
		})
	}
}

func Test_get_todos_by_title(t *testing.T) {
	//define a struct for usecase.
	type usecase struct {
		name           string
		input_user     *model.User_claims
		input_title    string
		repo_todos     []model.Todo
		repo_title     string
		repo_err       error
		expected_todos []model.Todo
		expected_err   error
		repo_call      bool
	}
	//set a table of usecase.
	tests := []usecase{
		{
			name:           "Error: Empty title",
			input_title:    "		",
			input_user:     &model.User_claims{Role: "user"},
			expected_todos: nil,
			expected_err:   ErrEmptyTitle,
			repo_call:      false,
		}, {
			name:           "Error: Server error",
			input_user:     &model.User_claims{Role: "user"},
			input_title:    "		task	",
			repo_err:       errors.New("Database connection failure"),
			repo_title:     "task",
			expected_todos: nil,
			expected_err:   ErrServerError,
			repo_call:      true,
		}, {
			name:           "Sucess: get todos",
			input_user:     &model.User_claims{Role: "user"},
			input_title:    "		task	",
			repo_err:       nil,
			repo_todos:     []model.Todo{{Title: "task 1", Category: "work"}, {Title: "task 2", Category: "work"}},
			repo_title:     "task",
			expected_todos: []model.Todo{{Title: "task 1", Category: "work"}, {Title: "task 2", Category: "work"}},
			expected_err:   nil,
			repo_call:      true,
		},
	}
	//run tests one by one.
	for _, usecase := range tests {
		t.Run(usecase.name, func(t *testing.T) {
			//set the mock.
			mock_repo := mocks.Mock_todo_repo{Wanted_err: usecase.repo_err, Wanted_list: usecase.repo_todos}
			//pass mock repo to initialize service.
			service := Init_service(&mock_repo)
			//call the method.
			todos, err := service.Get_todos_by_title(usecase.input_title, usecase.input_user)
			//check the response.
			if err != usecase.expected_err {
				t.Errorf("Response: expected error in response %v got %v", usecase.expected_err, err)
			}
			//check the todos in response.
			if usecase.expected_todos == nil && todos != nil {
				t.Errorf("Response: expected nil slice got %+v todos", todos)
			} else if expected := len(usecase.expected_todos); expected != len(todos) {
				t.Errorf("Response: expected slice of length %d got slice of length %d", expected, len(todos))
			}
			//check the repo call.
			if mock_repo.Was_called != usecase.repo_call {
				t.Errorf("Behaviour: expected repository method call %t got %t", usecase.repo_call, mock_repo.Was_called)
			}
			if mock_repo.Received_title != usecase.repo_title {
				t.Errorf("Logic: expected title passed to repository method %s got %s", usecase.repo_title, mock_repo.Received_title)
			}
		})
	}
}
