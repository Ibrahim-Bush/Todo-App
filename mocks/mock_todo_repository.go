package mocks

import (
	"Permanent-Storage-Todo/model"
)

type Mock_todo_repo struct {
	Received_todo     *model.Todo
	Received_id       int
	Received_category string
	Received_status   bool
	Received_title    string
	Wanted_err        error
	Wanted_list       []model.Todo
	Wanted_todo       *model.Todo
	Was_called        bool
	Get_by_id_err     error
	Get_by_id_call    bool
}

func (mock *Mock_todo_repo) Create_record(todo *model.Todo) error {
	//change is called value.
	mock.Was_called = true
	//strore the received struct.
	mock.Received_todo = todo
	//return the required error
	return mock.Wanted_err
}

func (mock *Mock_todo_repo) Update_record_by_id(new_todo *model.Todo) error {
	//change is called value.
	mock.Was_called = true
	//strore the received struct.
	mock.Received_todo = new_todo
	//return the required error
	return mock.Wanted_err
}
func (mock *Mock_todo_repo) Update_records_by_category(category string, updates *model.Todo, user *model.User_claims) error {
	//change is called value.
	mock.Was_called = true
	//strore the received struct and category.
	mock.Received_todo = updates
	mock.Received_category = category
	//return the required error
	return mock.Wanted_err
}

func (mock *Mock_todo_repo) Delete_record_by_id(id int, user *model.User_claims) error {
	//change is called value.
	mock.Was_called = true
	//strore the received id.
	mock.Received_id = id
	//return the required error
	return mock.Wanted_err
}
func (mock *Mock_todo_repo) Delete_all_records() error {
	//change is called value.
	mock.Was_called = true
	//return the required error
	return mock.Wanted_err
}

func (mock *Mock_todo_repo) Get_all_records(user *model.User_claims) ([]model.Todo, error) {
	//change is called value.
	mock.Was_called = true
	//return the required values.
	return mock.Wanted_list, mock.Wanted_err
}
func (mock *Mock_todo_repo) Get_record_by_id(id int, user *model.User_claims) (*model.Todo, error) {
	//change was called to true.
	mock.Get_by_id_call = true
	//store the received id.
	mock.Received_id = id
	//return the required values.
	return mock.Wanted_todo, mock.Get_by_id_err
}
func (mock *Mock_todo_repo) Get_records_by_category(category string, user *model.User_claims) ([]model.Todo, error) {
	//change is called value.
	mock.Was_called = true
	//strore the received category.
	mock.Received_category = category
	//return the required error
	return mock.Wanted_list, mock.Wanted_err
}
func (mock *Mock_todo_repo) Get_records_by_status(status bool, user *model.User_claims) ([]model.Todo, error) {
	//change is called value.
	mock.Was_called = true
	//strore the received status.
	mock.Received_status = status
	//return the required values.
	return mock.Wanted_list, mock.Wanted_err
}
func (mock *Mock_todo_repo) Get_records_by_title(title string, user *model.User_claims) ([]model.Todo, error) {
	//change is called value.
	mock.Was_called = true
	//strore the received title.
	mock.Received_title = title
	//return the required values.
	return mock.Wanted_list, mock.Wanted_err
}
