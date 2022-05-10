package models

// Event represents a change in data in user's storage
type Event struct {
	Operation Operation
	// Item should be one of these objects:
	//	- TextItem
	//	- BinaryItem
	//	- PasswordItem
	//	- CardItem
	Item interface{}
}

// Operation specifies which operation to take on the data.
type Operation string

const (
	OpCreate Operation = "CREATE"
	OpUpdate Operation = "UPDATE"
	OpDelete Operation = "DELETE"
)

func (o Operation) Valid() bool {
	return o == OpCreate ||
		o == OpUpdate ||
		o == OpDelete
}
