package util

import (
	"errors"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

// IsDuplicateErr checks if the error generated comes from inserting a duplicate record
func IsDuplicateErr(err error) bool {
	pqerr, ok := err.(*pq.Error)
	return ok && pqerr.Code.Name() == "unique_violation"
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
	resp.Data = nil
	resp.Error = "Internal error"
	log.Printf("Database error: %s\n", err.Error())
	return 500, resp
}
