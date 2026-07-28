# Todo-App

## Overview

This is a simple RESTful API for a Todo App using Go and the `gin-gonic/gin` framework for routing. The app will manage todo items stored in memory (using a slice) and support basic CRUD operations (Create, Read, Update, Delete).

## Features

Our app has the following features:

- #### Implementing Five Endpoints:
    -  GET /todos: Retrieve all todos as a JSON array 
    -  GET /todos/:id: Retrieve a specific todo by its ID
    -  POST /todos: Create a new todo with a title and completion status
    -  PUT /todos/:id: Update a todo’s title or completion status by its ID
    -  DELETE /todos/:id: Delete a todo by its ID
  
 - #### Input Handling:
    - For `GET /todos/:id` , `PUT /todos/:id` and `DELETE /todos/:id`:
            Validate the input ID, ensure it is a valid number and return 400 if it is invalid.
    - For `POST/todos` , `PUT/todos/:id`: Check the JSON body and validates incoming requests and returns 400 ("invalid json")
    if the JSON structure or types are incompatible or returns 400 ("Title cannot be empty") if there is no title field or the title is empty.
    For requests that omit the completed field, it automatically defaults to false.
- #### Error Handling: 
     Return appropriate HTTP status codes with clear messages to explain the error:
    - 200: Successful GET, POST, PUT, or DELETE.
    - 400: Invalid JSON or empty title or invalid ID.
    - 404: Todo not found for GET, PUT, or DELETE.
- #### In-Memory Storage
    - Store todos in a slice within the Go program.
    - Assign incremental IDs to new todos (e.g., 1, 2, 3, ...).
  
- #### RESTful Interface
    The API allows you to perform these operations smoothly:
    - Create a todo with a title and status.
    - Retrieve all todos or a specific todo.
    - Update a todo’s details.
    - Delete a todo.
    - View error messages for invalid operations.
  
## How To Run

After installing **Go** on your device and downloading the **Gin** framework run the following command in your terminal:
```
go run .
``` 
After this command the API server will start and run locally on
`http://localhost:8080` listening for incoming HTTP requests.

#### Example Usage:
- Try to send a POST request to `http://localhost:8080/todos` with the JSON body:
  ```
  {"title": "Read book", "completed": false}
  ```
- The system assigns an ID (e.g., 1) and stores the todo in the slice.
- The response will be:
  ```
  {
    "id": 1,
    "title": "Read book",
    "completed": false
  }
  ```

  
  

