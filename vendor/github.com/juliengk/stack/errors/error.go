package errors

import (
	"fmt"
)

type Error struct {
	ErrorCode int    `json:"code"`
	Message   string `json:"message"`
}

type Category int

type Reason int

const (
	None Reason = iota
)

func DefaultCategoryErrorString(category string, reason Reason) string {
	return fmt.Sprintf("Unsupported error reason %d under category %s.", reason, category)
}

func DefaultTypeErrorString(category Category) string {
	return fmt.Sprintf("Unsupported error type: %d.", category)
}
