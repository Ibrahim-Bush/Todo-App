package mocks

import (
	"Permanent-Storage-Todo/model"
	"Permanent-Storage-Todo/utils"
	"testing"
)

type Mock_user_repo struct {
	T             *testing.T
	Expected_user *model.User
	Wanted_err    error
	Wanted_user   *model.User
	Was_called    bool
}

func (mock *Mock_user_repo) Create_record(new_user *model.User) error {
	//change was_called to true.
	mock.Was_called = true
	//first check the passed struct pointer.
	if mock.Expected_user == nil && new_user == nil {
		return mock.Wanted_err
	} else if mock.Expected_user == nil || new_user == nil {
		mock.T.Fatalf("Boundary: Expected struct %v got %v", mock.Expected_user, new_user)
	}
	//compare passed struct fields to the expected struct fields.
	if new_user.Username != mock.Expected_user.Username {
		mock.T.Errorf("Logic: Expected username %s got %s", mock.Expected_user.Username, new_user.Username)
	}
	if new_user.Password == mock.Expected_user.Password {
		mock.T.Errorf("Logic: Expected hashed Password got password as plain text")
	} else if utils.Compare_password(mock.Expected_user.Password, new_user.Password) {
		mock.T.Errorf("Logic: Hashed password does not match original password")
	}
	if new_user.Role != mock.Expected_user.Role {
		mock.T.Errorf("Logic : Expected role %s got %s", mock.Expected_user.Role, new_user.Role)
	}
	//after successful logic testing return required value.
	return mock.Wanted_err
}

func (mock *Mock_user_repo) Get_record_by_name(name string) (*model.User, error) {
	//operate only if the expected user is not nil.
	if mock.Expected_user != nil {
		//check the passed username.
		if name != mock.Expected_user.Username {
			mock.T.Errorf("Logic: Expected username %s got %s", mock.Expected_user.Username, name)
		}
	}
	//return required value.
	return mock.Wanted_user, mock.Wanted_err
}
