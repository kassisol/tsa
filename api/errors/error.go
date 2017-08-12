package errors

import (
	"github.com/juliengk/stack/errors"
)

const (
	Success errors.Category = 10000 * iota // 0XXX

	DatabaseError // 1XXX
)

const (
	Unknown    errors.Reason = iota // X000
	ReadFailed                      // X001
)

func New(category errors.Category, reason errors.Reason) *errors.Error {
	errorCode := int(category) + int(reason)
	var msg string
	switch category {
	case DatabaseError:
		switch reason {
		case Unknown:
			msg = "Unknown database error"
		case ReadFailed:
			msg = "Failed to read database"
		}
	default:
		msg = errors.DefaultTypeErrorString(category)
	}

	return &errors.Error{ErrorCode: errorCode, Message: msg}
}
