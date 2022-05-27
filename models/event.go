package models

// Event represents a change in data in user's storage
type Event struct {
	Operation Operation
	Item      Item
}

// Operation specifies which operation to take on the data.
type Operation string

const (
	OpCreate Operation = "CREATE"
	OpUpdate Operation = "UPDATE"
)

func (o Operation) Valid() bool {
	return o == OpCreate || o == OpUpdate
}
