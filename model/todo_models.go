package model

import "time"

type Todo struct {
	ID          int    	   `json:"id" gorm:"primaryKey"`
	Title       string	   `json:"title" gorm:"type:varchar(255);not null"`
	Completed   bool  	   `json:"completed" gorm:"default:false;not null"`
	Category    string	   `json:"category" gorm:"type:varchar(50);not null"`
	Priority    string 	   `json:"priority" gorm:"type:varchar(10);not null"`
	CompletedAt *time.Time `json:"completedAt" gorm:"default:null"`
	DueDate		*time.Time `json:"dueDate" gorm:"default:null"`
	UserID      int		   `json:"-" gorm:"not null"`
}

type Create_todo_request struct{
	Title		string 	 `json:"title"`
	Completed   bool	 `json:"completed"`
	Category	string	 `json:"category"`
	Priority	string	 `json:"priority"`
	DueDate		string	 `json:"dueDate"`
}

type Update_todo_request struct{
	//we use pointers to differentiate default zero-valued from user input.
	Title	  *string	`json:"title"`
	Completed *bool		`json:"completed"`
	Category  *string	`json:"category"`
	Priority  *string	`json:"priority"`
	DueDate	  *string	`json:"dueDate"`
}

