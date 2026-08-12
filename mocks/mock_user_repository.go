package mocks

import (
	"Permanent-Storage-Todo/model"
)

type Mock_user_repo struct {
	Received_user *model.User
	Received_name string
	Wanted_err    error
	Wanted_user   *model.User
	Was_called    bool
}

func (mock *Mock_user_repo) Create_record(new_user *model.User) error {
	//change was_called to true.
	mock.Was_called = true
	//store the received user in the struct.
	mock.Received_user = new_user
	//return required value.
	return mock.Wanted_err
}

func (mock *Mock_user_repo) Get_record_by_name(name string) (*model.User, error) {
	//change was_called to true.
	mock.Was_called = true
	//store the received name in the struct.
	mock.Received_name = name
	//return required value.
	return mock.Wanted_user, mock.Wanted_err
}

func (repo *Mock_user_repo) Get_record_by_id(id int) (*model.User, error) {
	return nil, nil
}
