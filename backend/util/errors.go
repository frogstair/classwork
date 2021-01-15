package util

import (
	"errors"

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
