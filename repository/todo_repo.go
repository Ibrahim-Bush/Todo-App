package repository

import (
	"Todo-App/model"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("Record not found")
)

type Todo_repo interface {
	Todo_writer
	Todo_reader
}

type Todo_writer interface {
	Create_record(todo *model.Todo) error
	Update_record_by_id(new_todo *model.Todo) error
	Update_records_by_category(category string, updates *model.Todo, user *model.User_claims) error
	Delete_record_by_id(id int, user *model.User_claims) error
	Delete_all_records() error
}

type Todo_reader interface {
	Get_all_records(user *model.User_claims) ([]model.Todo, error)
	Get_record_by_id(id int, user *model.User_claims) (*model.Todo, error)
	Get_records_by_category(category string, user *model.User_claims) ([]model.Todo, error)
	Get_records_by_status(status bool, user *model.User_claims) ([]model.Todo, error)
	Get_records_by_title(title string, user *model.User_claims) ([]model.Todo, error)
}

type postgres_repo struct {
	db *gorm.DB
}

func Init_postgres_repo(new_db *gorm.DB) Todo_repo {
	//initialize new repo and return its address as a repo interface.
	var new_repo = postgres_repo{db: new_db}
	return &new_repo
}

func (repo *postgres_repo) Create_record(todo *model.Todo) error {
	//add the new object to the database.
	result := repo.db.Create(todo)
	//return the result.
	return result.Error
}

func (repo *postgres_repo) Update_record_by_id(new_todo *model.Todo) error {
	//save the new modified object.
	result := repo.db.Save(new_todo)
	//return the result of process.
	return result.Error
}

func (repo *postgres_repo) Update_records_by_category(category string, updates *model.Todo, user *model.User_claims) error {
	//define variable to check result.
	var result *gorm.DB
	//differentiate between users based on role.
	if user.Role == "admin" {
		//update completed status of all records in the passed category.
		result = repo.db.Model(&model.Todo{}).Where("category = ?", category).Select("completed", "completed_at").
			Updates(updates)
	} else {
		//update only todos of that category if the user was its creator.
		result = repo.db.Model(&model.Todo{}).Where("category = ? AND user_id = ?", category, user.UserID).
			Select("completed", "completed_at").Updates(updates)
	}

	return result.Error
}

func (repo *postgres_repo) Delete_record_by_id(id int, user *model.User_claims) error {
	//define a result to check errors.
	var result *gorm.DB
	//differentiate queries based on user role.
	if user.Role == "admin" {
		//delete the todo that have the passed id anyway.
		result = repo.db.Where("id = ?", id).Delete(&model.Todo{})
	} else {
		//delete the todo if the user was its creator.
		result = repo.db.Where("id = ? AND user_id = ?", id, user.UserID).Delete(&model.Todo{})
	}
	//check the result
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return ErrNotFound
	} else {
		return nil
	}
}

func (repo *postgres_repo) Delete_all_records() error {
	//delete all records in the database.
	result := repo.db.Exec("DELETE FROM todos")
	return result.Error
}

func (repo *postgres_repo) Get_all_records(user *model.User_claims) ([]model.Todo, error) {
	//initialize a slice for records.
	records := make([]model.Todo, 0)
	//define a variable to check result.
	var result *gorm.DB
	//differentiate users based on role.
	if user.Role == "admin" {
		//get all records in the database.
		result = repo.db.Order("id asc").Find(&records)
	} else {
		//get current user's records only.
		result = repo.db.Order("id asc").Where("user_id = ?", user.UserID).Find(&records)
	}
	//if something went wrong, return empty slice and error value.
	if result.Error != nil {
		return nil, result.Error
	}
	return records, nil
}

func (repo *postgres_repo) Get_record_by_id(id int, user *model.User_claims) (*model.Todo, error) {
	//define an object for the target record.
	var todo model.Todo
	//define variable to check the result.
	var result *gorm.DB
	//differentiate users based on role.
	if user.Role == "admin" {
		//get the record anyway.
		result = repo.db.First(&todo, id)
	} else {
		//get element if the user was its creator only.
		result = repo.db.Where("id = ? AND user_id = ?", id, user.UserID).First((&todo))
	}
	//check the result of process.
	switch {
	case result.Error == nil:
		return &todo, nil
	case errors.Is(result.Error, gorm.ErrRecordNotFound):
		return nil, ErrNotFound
	default:
		return nil, result.Error
	}
}

func (repo *postgres_repo) Get_records_by_category(category string, user *model.User_claims) ([]model.Todo, error) {
	//define a slice for records.
	records := make([]model.Todo, 0)
	//define a variable to check the result.
	var result *gorm.DB
	//differentiate users based on role.
	if user.Role == "admin" {
		//get all records of the target category.
		result = repo.db.Where("category = ?", category).Find(&records)
	} else {
		//get records of the current user only.
		result = repo.db.Where("category = ? AND user_id = ?", category, user.UserID).Find(&records)
	}

	//check the result of process.
	if result.Error != nil {
		//return empty slice and error value.
		return nil, result.Error
	}
	return records, nil
}

func (repo *postgres_repo) Get_records_by_status(status bool, user *model.User_claims) ([]model.Todo, error) {
	//define a slice for the records.
	records := make([]model.Todo, 0)
	//define a variable to check result.
	var result *gorm.DB
	//differentiate users based on role.
	if user.Role == "admin" {
		//get all records that has the required status.
		result = repo.db.Where("completed = ?", status).Find(&records)
	} else {
		//get the current user's todos only.
		result = repo.db.Where("user_id = ? AND completed = ?", user.UserID, status).Find(&records)
	}
	//check the result of the process.
	if result.Error != nil {
		//return empty slice and error value.
		return nil, result.Error
	}
	return records, nil
}

func (repo *postgres_repo) Get_records_by_title(title string, user *model.User_claims) ([]model.Todo, error) {
	//define a slice for records.
	records := make([]model.Todo, 0)
	//define a variable to check result.
	var result *gorm.DB
	//differentiate users based on role.
	if user.Role == "admin" {
		//get all records that match the required title.
		result = repo.db.Where("title ILIKE ?", "%"+title+"%").Find(&records)
	} else {
		//get todos of the current user only.
		result = repo.db.Where("user_id = ?", user.UserID).Where("title ILIKE ?", "%"+title+"%").Find(&records)
	}

	//check the result of the process.
	if result.Error != nil {
		return nil, result.Error
	}
	return records, nil
}
