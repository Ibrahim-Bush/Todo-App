package main

import (
	"Permanent-Storage-Todo/config"
	"Permanent-Storage-Todo/handler"
	"Permanent-Storage-Todo/repository"
	"Permanent-Storage-Todo/router"
	"Permanent-Storage-Todo/service"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	//init router with default middleware.
	Router := gin.Default()

	//then: get a connection to database.
	connection_data := "host=localhost user=postgres password=admin dbname=todo_app port=5432 sslmode=disable"
	database, err := database.Init_database(connection_data)
	//if something went wrong exit.
	if err != nil {
		log.Fatal(err)
	}

	//linking database to repository layer.
	todo_repo := repository.Init_postgres_repo(database)
	user_repo := repository.Init_repo(database)

	//linking repository layer with service layer.
	todo_service := service.Init_service(todo_repo)
	user_service := service.Init_user_service(user_repo)

	//linking service layer to handler layer.
	todo_handler := handler.Init_handler(todo_service)
	user_handler := handler.Init_user_handler(user_service)

	//linking the handler layer to the router.
	router.Init_router(Router, todo_handler)
	router.Init_user_router(Router, user_handler)

	//run the server at local port ":8080".
	Router.Run(":8080")

}
