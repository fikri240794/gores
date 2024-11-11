# Go Response View Model (gores)
Go response view model is used as the standard go struct for RESTful API responses in JSON format.

## Installation
```bash
go get github.com/fikri240794/gores
```

## Example Usage
Run this simple HTTP Server:
```go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/fikri240794/gocerr"
	"github.com/fikri240794/gores"
)

type SomeData struct {
	SomeFieldOne int    `json:"some_field_one"`
	SomeFieldTwo string `json:"some_field_two"`
}

func SuccessHandler(w http.ResponseWriter, r *http.Request) {
	response := gores.NewResponseVM[*SomeData]().
		SetCode(http.StatusOK).
		SetData(&SomeData{
			SomeFieldOne: 1,
			SomeFieldTwo: "two",
		})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Code)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	err := errors.New("something error")
	response := gores.NewResponseVM[*SomeData]().
		SetErrorFromError(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Code)

	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func BadRequestHandler(w http.ResponseWriter, r *http.Request) {
	var err error = gocerr.New(
		http.StatusBadRequest,
		http.StatusText(http.StatusBadRequest),
		gocerr.NewErrorField("some_field_one", "some_field_one is required"),
		gocerr.NewErrorField("some_field_two", "some_field_two is not unique"),
	)
	response := gores.NewResponseVM[*SomeData]().
		SetErrorFromError(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Code)

	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/success", SuccessHandler)
	http.HandleFunc("/error", InternalServerErrorHandler)
	http.HandleFunc("/bad", BadRequestHandler)

	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
```

Call `/success` API:
```curl
curl localhost:8080/success
```
Response:
```json
{
    "code": 200,
    "data": {
        "some_field_one": 1,
        "some_field_two": "two"
    }
}
```

Call `/error` API:
```curl
curl localhost:8080/error
```
Response:
```json
{
    "code": 500,
    "error": {
        "message": "something error"
    }
}
```

Call `/bad` API:
```curl
curl localhost:8080/bad
```
Response:
```json
{
    "code": 400,
    "error": {
        "message": "Bad Request",
        "error_fields": [
            {
                "field": "some_field_one",
                "message": "some_field_one is required"
            },
            {
                "field": "some_field_two",
                "message": "some_field_two is not unique"
            }
        ]
    }
}
```