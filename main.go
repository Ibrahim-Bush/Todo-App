package main

import (
	database "Permanent-Storage-Todo/config"
	_ "Permanent-Storage-Todo/docs"
	"Permanent-Storage-Todo/handler"
	"Permanent-Storage-Todo/repository"
	"Permanent-Storage-Todo/router"
	"Permanent-Storage-Todo/service"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title			Todo App API
// @version 		1.0
// @description 	RESTful API for managing todos with GORM and Gin framework.
// @host			localhost:8080
// @BasePath		/

func main() {

	//init router with default middleware.
	Router := gin.Default()

	//add the swagger documentation on router.
	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//then: get a connection to database.
	connection_data := "host=db user=postgres password=admin dbname=todo_app port=5432 sslmode=disable"
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
