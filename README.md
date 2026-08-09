# Permanent-Storage-Todo-App

## Overview

This is a simple RESTful API for a Todo App using Go and the `gin-gonic/gin` framework for routing. The app will manage todo items stored in a database (using GORM) and support basic CRUD operations (Create, Read, Update, Delete)  with features like categorization, prioritization, completion tracking, due dates, status filtering, title search, and bulk updates. The app is secured by adding user authentication (signup/login), JWT-based session handling, and role-based authorization.

## Features

Our app has the following features:

### Todos Endpoints:
  - **GET `/todos`**: Retrieve all todos as a JSON array.
  - **GET `/todos/:id`**: Retrieve a specific todo by its ID.
  - **GET `/todos/category/:category`**: Retrieve all todos in a specific category.
  - **GET `/todos/status/:status`**: Retrieve all todos with a specific completion status.
  - **GET `/todos/search`**: Retrieve todos with titles containing a query parameter `q` (case-insensitive partial match).
  - **POST `/todos`**: Create a new todo with title, completed status, category, priority, and optional due date.
  - **PUT `/todos/:id`**: Update a todo’s title, completed status, category, priority, or due date by its ID.
  - **PUT `/todos/category/:category`**: Update the `completed` status of all todos in a specific category.
  - **DELETE `/todos/:id`**: Delete a todo by its ID.
  - **DELETE `/todos`**: Delete all todos in the database.
  
### Authentication Endpoints:
  - **POST `/signup`**: Registers in the system with unique username and password.
  - **POST `/login`**: Login to the system with your username and password.

  
 ### Input Handling:

  - For GET , PUT, DELETE requests:
    - Validate the input ID, ensure it is a valid number and return (400, Invalid id) if it is invalid.  
  - For POST and PUT requests: 
    - Check the JSON body and validate incoming requests and return (400, "invalid json") if the JSON structure or types are incompatible.
    - Validate input to ensure:
          `title` is non-empty,
          `priority` is one of "Low", "Medium", or "High",
          and `dueDate`, if provided, is a valid ISO 8601 date and not in the past (relative to the current date).
    - Return (400, "Field cannot be empty") if there is no input for required field or the field is empty.
  -  For requests that omit an optional field:
    - it automatically defaults to its zero-value (false, NULL).
### Error Handling: 
  Return appropriate HTTP status codes with clear messages to explain the error:
  - 200: Successful operations.
  - 201: User created (on signup).
  - 400: Invalid input or JSON.
  - 401: Unauthorized (missing or invalid JWT, invalid login).
  - 403: Forbidden (user lacks role permissions).
  - 404: Resource not found (e.g., deleting non-existent todo).
  - 500: Internal server/database error.
  
### Permanent Data Storage
  - Replaced the temporary in-memory slice implementation with a robust postgreSQL database integration using GORM ORM, ensuring all data persists across server restarts.
  - Integrated GORM's `AutoMigrate` feature to automatically generate and maintain the database schema directly from Go structs without requiring manual SQL table creation.
  - Assigns incremental IDs to new todos and users (e.g., 1, 2, 3, ...).
  
### Clean Code & Architecture
  This project is built following clean code principles and modular architecture to ensure high maintainability, readability and scalable growth.
  
#### Key Architectural Highlights:
  - **Separation of Concerns**: The application logic is decoupled into **distinct layers** (model, repository, service, handler, and database configurations) to keep business logic isolated and easy to test.
  -  **Interface Based Architecture**: Lower layers define clear **Go interfaces** instead of relying on fixed implementations. This decouples the layers, make it easy to replace database implementations without touching the other layers, and allows simple mocking for unit tests.
  -  **Layer-by-Layer Validation**: Data is validated at each layer before calling the next lower layer. This ensure that every layer only receives correct and safe data, detecting errors as early as possible.
  -  **Clear Error Handling**: Errors are handled at their origin, returning simple **JSON responses** with standard HTTP status codes for easy debugging.

### Authentication & Authorization:
  This project provides a secure **JWT authentication** and **Role-Based Access Control** (RBAC) to ensure strict data privacy and isolation.

  - **JWT Authentication**: Protects all todo endpoints by requiring a **valid JSON Web Token** generated upon Login.
  - **Strict Data Isolation**: Restrict database queries to the logged-in user's ID, ensuring regular users can only view, update, or delete their own todo items.
  - **Least Privilege**: Automatically assigns all **new registrations** the standard **user role** by default, ensuring accounts start with the minimal permissions.
  - **Full Secrecy**: Returns `404 Not found` response for any item not belonging to the current user, completely **hiding unowned data** to prevent users from discovering system or attempting **unauthorized tampering**.
  - **Admin Privileges**: Grants administrative accounts **system-wide access** to inspect and manage records across all users.
  
## How To Run

Follow these steps to run the code on your local device:
- Installing **Go** compiler on your device.
- Downloading the **Gin** framework for routing and the **GORM** package on your device to act as a bridge between Go code and postgreSQL.   
  To download them run the following command in your terminal:
 ```
 go mod download
 ```
- Create a postgreSQL database with the following configurations:
  - Host: localhost
  - Port: 5432
  - User: postgres
  - Password: admin
  - Database Name: todo_app

  
After downloading the previous dependencies run the following command in your terminal:
```
go run .
``` 
After this command the API server will start and run locally on
`http://localhost:8080` listening for incoming HTTP requests.

### Practical Example:

- First: your must signup in the system: send a **POST** request to `http://localhost:8080/signup` with the  JSON body:
  ```
  {"username": "your_name", "password": "your_secret_password"}
  ```
- The system will store your account and you will use your credentials to **login to the system** at any time.

- Then **Login** to the system: send a **POST** request to `http://localhost:8080/login` with the  JSON body:
  ```
  {"username": "your_name", "password": "your_secret_password"}
  ```
- The system will respond if your credentials is valid with a **JWT token** and you will attach it to any request by putting it in the authorization header.

- Try to send a **POST** request to `http://localhost:8080/todos` with the JSON body:
  ```
  {"title": "Read book", "category": "study", "priority": "medium", completed": false}
  ```
- The system assigns an ID (e.g., 1) and stores the todo in the database.
- The response will be:
  ```
  {
    "id": 1,
    "title": "Read book",
    "category": study,
    "priority": medium,
    "completed": false,
    "completedAt": null,
    "dueDate": null
  }
  ```