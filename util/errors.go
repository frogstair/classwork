package util

import (
	"errors"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

// IsDuplicateErr checks if the error generated comes from inserting a duplicate record
func IsDuplicateErr(err error) bool {
	pqerr, ok := err.(*pq.Error) // Get the PostgreSQL error code
	return ok && pqerr.Code.Name() == "unique_violation" // Return true if the 
}

// IsNotFoundErr checks if the returned error signifies a missing record
func IsNotFoundErr(err error) bool {
	return err != nil && errors.Is(err, gorm.ErrRecordNotFound)
}

// Response is the response struct, that will be sent back to the user
type Response struct {
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

// DatabaseError notifies of a database error
func DatabaseError(err error, resp *Response) (int, *Response) {
	resp.Data = nil // Set the resonse data to null which will make it missing from the response JSON
	resp.Error = "Internal error" // Set the error
	log.Printf("Database error: %s\n", err.Error())
	return 500, resp // Return a 500 error
}
