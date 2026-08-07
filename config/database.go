package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"Permanent-Storage-Todo/model"
)

func Init_database(connection_data string) (*gorm.DB, error) {

	//first get connection to the database.
	db, err := gorm.Open(postgres.Open(connection_data), &gorm.Config{})
	//check if something went wrong.
	if err != nil{
		return nil, err
	}
	//create the table with AutoMigrate.
	err = db.AutoMigrate(&model.Todo{})
	if err != nil{
		return nil, err
	}
	//create a table for users.
	err = db.AutoMigrate(&model.User{})
	if err != nil{
		return nil, err
	}
	//if the connection and autoMigration succeeded, return the database object.
	return db, nil
}