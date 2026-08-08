package router

import (
	"Permanent-Storage-Todo/handler"
	"Permanent-Storage-Todo/middleware"

	"github.com/gin-gonic/gin"
)

func Init_router(router *gin.Engine, handler *handler.Todo_handler) {

	//create a routes group for todos.
	todosAPI := router.Group("/todos")

	//apply authentication middleware to todos group.
	todosAPI.Use(middleware.Auth_middleware)

	//register routes.
	todosAPI.GET("", handler.Get_all_todos_handler)
	todosAPI.GET("/search", handler.Get_todos_by_title_handler)
	todosAPI.GET("/category/:category", handler.Get_todos_by_category_handler)
	todosAPI.GET("/status/:status", handler.Get_todos_by_status_handler)
	todosAPI.GET("/:id", handler.Get_todo_by_id_handler)
	todosAPI.POST("", handler.Create_todo_handler)
	todosAPI.PUT("/:id", handler.Update_todo_by_id_handler)
	todosAPI.PUT("/category/:category", handler.Update_todos_by_category_handler)
	todosAPI.DELETE("/:id", handler.Delete_todo_by_id_handler)
	todosAPI.DELETE("", middleware.Authorization_middleware, handler.Delete_all_todos_handler)

}
