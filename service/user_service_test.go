package service

import (
	"Permanent-Storage-Todo/mocks"
	"Permanent-Storage-Todo/model"
	"Permanent-Storage-Todo/repository"
	"Permanent-Storage-Todo/utils"
	"errors"
	"strings"
	"testing"
)

func Test_create_user(t *testing.T) {

	//define a struct for usecases.
	type usecase struct {
		name          string
		input_data    model.Auth_request
		repo_user     *model.User
		repo_err      error
		expected_err  error
		expected_user *model.User
		call_repo     bool
	}
	//create a table of usecases.
	tests := []usecase{
		{
			name:          "Error: Empty username",
			input_data:    model.Auth_request{Username: "	", Password: ""},
			repo_user:     nil,
			repo_err:      nil,
			expected_err:  ErrEmptyName,
			expected_user: nil,
			call_repo:     false,
		}, {
			name:          "Error: Empty password",
			input_data:    model.Auth_request{Username: " ahmed ", Password: "		"},
			repo_user:     nil,
			repo_err:      nil,
			expected_err:  ErrEmptyPassword,
			expected_user: nil,
			call_repo:     false,
		}, {
			name:          "Error: Invalid password",
			input_data:    model.Auth_request{Username: " ahmed ", Password: strings.Repeat("a", 200)},
			repo_user:     nil,
			repo_err:      nil,
			expected_err:  ErrInvalidPassword,
			expected_user: nil,
			call_repo:     false,
		}, {
			name:          "Error: Duplicate username",
			input_data:    model.Auth_request{Username: "	AhMeD	", Password: "	Admin	"},
			repo_user:     &model.User{Username: "ahmed", Password: "Admin", Role: "user"},
			repo_err:      repository.ErrDuplicatedName,
			expected_err:  ErrDuplicatedName,
			expected_user: nil,
			call_repo:     true,
		}, {
			name:          "Error: Server error",
			input_data:    model.Auth_request{Username: "	AhMeD	", Password: "	Admin	"},
			repo_user:     &model.User{Username: "ahmed", Password: "Admin", Role: "user"},
			repo_err:      errors.New("database connection failure"),
			expected_err:  ErrServerError,
			expected_user: nil,
			call_repo:     true,
		}, {
			name:          "Success: user created",
			input_data:    model.Auth_request{Username: "	AhMeD	", Password: "	Admin	"},
			repo_user:     &model.User{Username: "ahmed", Password: "Admin", Role: "user"},
			repo_err:      nil,
			expected_err:  nil,
			expected_user: &model.User{Username: "ahmed", Password: "Admin", Role: "user"},
			call_repo:     true,
		},
	}
	//run tests one by one.
	for _, usecase := range tests {
		t.Run(usecase.name, func(t *testing.T) {
			//create a new mock repo.
			mock_repo := &mocks.Mock_user_repo{
				Wanted_err: usecase.repo_err,
			}
			//pass mock to the service.
			service := Init_user_service(mock_repo)
			//call method with usecases.
			user, err := service.Create_user(usecase.input_data)
			//check the repository call: first check the call status.
			if mock_repo.Was_called != usecase.call_repo {
				t.Errorf("Behaviour: Expected repository call %t got %t", usecase.call_repo, mock_repo.Was_called)
			}
			//compare expected repo struct to the received struct.
			compare_user_structs(t, "Repo input", mock_repo.Received_user, usecase.repo_user)
			//check responeses.
			if err != usecase.expected_err {
				t.Errorf("Response: Expected error %v got %v", usecase.expected_err, err)
			}
			//check fields in response.
			compare_user_structs(t, "Service response", user, usecase.expected_user)
		})
	}
}

func compare_user_structs(t *testing.T, stage string, received *model.User, expected *model.User) {
	//first compare pointers.
	if received == nil || expected == nil {
		if received == nil && expected != nil {
			t.Errorf("Boundary (%s): expected struct %+v got nil pointer", stage, expected)
		} else if received != nil && expected == nil {
			t.Errorf("Boundary (%s): expected nil struct got %+v", stage, received)
		}
	} else {
		//compare received struct fields to the expected struct fields.
		if received.Username != expected.Username {
			t.Errorf("Logic (%s): expected username %s got %s", stage, expected.Username, received.Username)
		}
		if received.Password == expected.Password {
			t.Errorf("Logic (%s): expected hashed Password got password as plain text", stage)
		} else if utils.Compare_password(expected.Password, received.Password) == false {
			t.Errorf("Logic (%s): Hashed password does not match original password", stage)
		}
		if received.Role != expected.Role {
			t.Errorf("Logic (%s): expected role %s got %s", stage, expected.Role, received.Role)
		}
	}
}

func Test_login_user(t *testing.T) {
	//define a struct for usecases.
	type usecase struct {
		name           string
		input_data     model.Auth_request
		repo_name      string
		repo_err       error
		repo_user      *model.User
		expected_token string
		expected_err   error
		call_repo      bool
	}
	//create a hashed password for success case.
	hashed_admin, err := utils.Hash_password("admin")
	if err != nil {
		t.Fatalf("Error from testing func: could not make a hashed password for success case")
	}
	//create a table of usecases.
	tests := []usecase{
		{
			name:           "Error: User not found",
			input_data:     model.Auth_request{Username: "  AhMEd	", Password: "		"},
			repo_name:      "ahmed",
			repo_err:       repository.ErrUserNotFound,
			repo_user:      nil,
			expected_token: "",
			expected_err:   ErrUserNotFound,
			call_repo:      true,
		}, {
			name:           "Error: Server error",
			input_data:     model.Auth_request{Username: "  AhMEd	", Password: "	abc	"},
			repo_name:      "ahmed",
			repo_err:       errors.New("Database connection failure"),
			repo_user:      nil,
			expected_token: "",
			expected_err:   ErrServerError,
			call_repo:      true,
		}, {
			name:           "Error: Wrong password",
			input_data:     model.Auth_request{Username: "  AhMEd	", Password: "	abc	"},
			repo_name:      "ahmed",
			repo_err:       nil,
			repo_user:      &model.User{Username: "ahmed", Password: "admin"},
			expected_token: "",
			expected_err:   ErrWrongPassword,
			call_repo:      true,
		}, {
			name:           "Error: Wrong password",
			input_data:     model.Auth_request{Username: "  AhMEd	", Password: "	abc	"},
			repo_name:      "ahmed",
			repo_err:       nil,
			repo_user:      &model.User{Username: "ahmed", Password: "admin"},
			expected_token: "",
			expected_err:   ErrWrongPassword,
			call_repo:      true,
		}, {
			name:           "Success: get token",
			input_data:     model.Auth_request{Username: "  AhMEd	", Password: "	admin	"},
			repo_name:      "ahmed",
			repo_err:       nil,
			repo_user:      &model.User{Username: "ahmed", Password: hashed_admin, ID: 1, Role: "user"},
			expected_token: "not null",
			expected_err:   nil,
			call_repo:      true,
		},
	}
	//run tests one by one.
	for _, usecase := range tests {
		t.Run(usecase.name, func(t *testing.T) {
			//create new mock repo with case data.
			mock_repo := &mocks.Mock_user_repo{
				Wanted_err:  usecase.repo_err,
				Wanted_user: usecase.repo_user,
			}
			//pass mock to the service.
			service := Init_user_service(mock_repo)
			//call method with usecase.
			token, err := service.Login_user(usecase.input_data)
			//check the repo input: first check the call status.
			if mock_repo.Was_called != usecase.call_repo {
				t.Errorf("Behaviour: Expected repository call %t got %t", usecase.call_repo, mock_repo.Was_called)
			}
			//check the fields in the repo call.
			if mock_repo.Received_name != usecase.repo_name {
				t.Errorf("Logic (Repo input): Expected username %s got %s", usecase.repo_name, mock_repo.Received_name)
			}
			//check the service response.
			if err != usecase.expected_err {
				t.Errorf("Response: Expected error %v got %v", usecase.expected_err, err)
			}
			if usecase.expected_token == "" && token != usecase.expected_token {
				t.Errorf("Response: Expected token %s got %s", usecase.expected_token, token)
			} else if usecase.expected_token != "" {
				//verify the token in response.
				claims, err := utils.Verify_token(token)
				if claims != nil && err == nil {
					if claims.Username != usecase.repo_user.Username {
						t.Errorf("Response: Expected username in token %s got %s", usecase.repo_user.Username, claims.Username)
					}
					if claims.UserID != usecase.repo_user.ID {
						t.Errorf("Response: Expected userID in token %d got %d", usecase.repo_user.ID, claims.UserID)
					}
					if claims.Role != usecase.repo_user.Role {
						t.Errorf("Response: Expected Role in token %s got %s", usecase.repo_user.Role, claims.Role)
					}
				}
			}
		})
	}
}
