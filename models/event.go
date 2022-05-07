package models

import "time"

const (
	OpCreate Operation = "CREATE"
	OpUpdate Operation = "UPDATE"
	OpDelete Operation = "DELETE"
)

type (
	Operation string

	Event struct {
		Timestamp time.Time
		Operation Operation
		Item      interface{}
	}
)
