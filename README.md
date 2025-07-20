# 🚀 GoRes - Go Response

<div align="center">

**A standardized, type-safe Go library for REST API response structures with comprehensive error handling**

</div>

---

## 📋 Overview

GoRes provides a standardized approach to handling HTTP responses in Go REST APIs. It offers type-safe response structures with consistent error handling and validation error mapping.

## ✨ Features

- **Generic Response Types**: Type-safe responses with `ResponseVM[T]`
- **Fluent API**: Method chaining for clean, readable code
- **Error Field Mapping**: Automatic validation error field extraction
- **HTTP Status Integration**: Seamless HTTP status code handling

## 📦 Installation

```bash
go get github.com/fikri240794/gores
```

### Requirements

- Go 1.18 or higher (for generics support)
- [gocerr](https://github.com/fikri240794/gocerr) (latest version)

## 🚀 Quick Start

```go
package main

import (
    "encoding/json"
    "net/http"
    "os"
    "github.com/fikri240794/gores"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

func main() {
    // Success response
    response := gores.NewResponseVM[*User]().
        SetCode(http.StatusOK).
        SetData(&User{ID: 1, Name: "John Doe"})

    json.NewEncoder(os.Stdout).Encode(response)
    // Output: {"code":200,"data":{"id":1,"name":"John Doe"}}
}
```

## 📖 Examples

### Success Response

```go
// Create a successful response with data
user := &User{ID: 1, Name: "Alice", Email: "alice@example.com"}

response := gores.NewResponseVM[*User]().
    SetCode(http.StatusOK).
    SetData(user)

// JSON Output:
// {
//   "code": 200,
//   "data": {
//     "id": 1,
//     "name": "Alice", 
//     "email": "alice@example.com"
//   }
// }
```

### Error Response

```go
// Standard error handling
err := errors.New("database connection failed")

response := gores.NewResponseVM[*User]().
    SetErrorFromError(err)

// JSON Output:
// {
//   "code": 500,
//   "error": {
//     "message": "database connection failed"
//   }
// }
```

### Validation Error Response

```go
// Validation error with field-specific messages
validationErr := gocerr.New(
    http.StatusUnprocessableEntity,
    "Validation failed",
    gocerr.NewErrorField("email", "Invalid email format"),
    gocerr.NewErrorField("age", "Age must be between 18 and 120"),
)

response := gores.NewResponseVM[*User]().
    SetErrorFromError(validationErr)

// JSON Output:
// {
//   "code": 422,
//   "error": {
//     "message": "Validation failed",
//     "error_fields": [
//       {
//         "field": "email",
//         "message": "Invalid email format"
//       },
//       {
//         "field": "age", 
//         "message": "Age must be between 18 and 120"
//       }
//     ]
//   }
// }
```

### Real-world Usage Pattern

```go
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        response := gores.NewResponseVM[*User]().
            SetErrorFromError(gocerr.New(
                http.StatusBadRequest,
                "Invalid JSON format",
            ))
        
        writeJSONResponse(w, response)
        return
    }

    // Validation
    if user.Name == "" || user.Email == "" {
        response := gores.NewResponseVM[*User]().
            SetErrorFromError(gocerr.New(
                http.StatusUnprocessableEntity,
                "Validation failed",
                gocerr.NewErrorField("name", "Name is required"),
                gocerr.NewErrorField("email", "Email is required"),
            ))
        
        writeJSONResponse(w, response)
        return
    }

    // Success response
    response := gores.NewResponseVM[*User]().
        SetCode(http.StatusCreated).
        SetData(&user)
    
    writeJSONResponse(w, response)
}
```

## 🏗️ API Reference

### Core Types

#### `ResponseVM[T]`
```go
type ResponseVM[T comparable] struct {
    Code  int              `json:"code"`
    Error *ResponseErrorVM `json:"error,omitempty"`
    Data  T                `json:"data,omitempty"`
}
```

#### `ResponseErrorVM`
```go
type ResponseErrorVM struct {
    Message     string                  `json:"message"`
    ErrorFields []*ResponseErrorFieldVM `json:"error_fields,omitempty"`
}
```

#### `ResponseErrorFieldVM`
```go
type ResponseErrorFieldVM struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}
```

### Methods

#### ResponseVM Methods
- `NewResponseVM[T]() *ResponseVM[T]` - Create new response instance
- `SetCode(code int) *ResponseVM[T]` - Set HTTP status code
- `SetData(data T) *ResponseVM[T]` - Set response data
- `SetError(err *ResponseErrorVM) *ResponseVM[T]` - Set error manually
- `SetErrorFromError(err error) *ResponseVM[T]` - Parse and set error from Go error

#### ResponseErrorVM Methods
- `NewResponseErrorVM() *ResponseErrorVM` - Create new error instance
- `SetMessage(message string) *ResponseErrorVM` - Set error message
- `AddErrorFields(fields ...*ResponseErrorFieldVM) *ResponseErrorVM` - Add field errors
- `ParseError(err error) *ResponseErrorVM` - Parse error from Go error

#### ResponseErrorFieldVM Methods
- `NewResponseErrorFieldVM(field, message string) *ResponseErrorFieldVM` - Create field error

## 🤝 Integration with gocerr

GoRes seamlessly integrates with [gocerr](https://github.com/fikri240794/gocerr) for enhanced error handling:

```go
// Leverage gocerr helper functions
if gocerr.HasErrorFields(err) {
    fmt.Printf("Field errors count: %d\n", gocerr.ErrorFieldCount(err))
}

if gocerr.HasErrorField(err, "email") {
    fmt.Printf("Email error: %s\n", gocerr.GetErrorFieldMessage(err, "email"))
}

// Automatic integration in responses
response := gores.NewResponseVM[*User]().
    SetErrorFromError(err) // Automatically extracts all gocerr details
```

## 📖 Complete Example

See the [examples](examples/) directory for a comprehensive demonstration including:

- Success and error response patterns
- Validation error handling with field-specific messages
- Method chaining examples
- Advanced gocerr features
- Different HTTP status codes handling

Run the example:

```bash
cd examples
go run main.go
```

This will demonstrate all core features with formatted JSON output showing how each response type is structured.

## 🛠️ Advanced Usage

### Custom Error Types

```go
// Custom error with specific HTTP status
customErr := gocerr.New(
    http.StatusTeapot,
    "I'm a teapot",
    gocerr.NewErrorField("brew_type", "Cannot brew coffee"),
)

response := gores.NewResponseVM[*TeaPot]().
    SetErrorFromError(customErr)
```

### Method Chaining Patterns

```go
// Complex chained response
response := gores.NewResponseVM[*User]().
    SetCode(http.StatusCreated).
    SetData(newUser).
    SetError(nil) // Can be chained even with nil

// Conditional chaining
response := gores.NewResponseVM[*User]()
if hasError {
    response = response.SetErrorFromError(validationErr)
} else {
    response = response.SetCode(http.StatusOK).SetData(user)
}
```

### Generic Type Usage

```go
// Different data types
type UserList struct {
    Users []User `json:"users"`
    Total int    `json:"total"`
}

// List response
listResponse := gores.NewResponseVM[*UserList]().
    SetCode(http.StatusOK).
    SetData(&UserList{Users: users, Total: len(users)})

// Empty response  
emptyResponse := gores.NewResponseVM[interface{}]().
    SetCode(http.StatusNoContent)
```

---