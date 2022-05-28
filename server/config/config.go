package config

type Config struct {
	// MaxDBConnections is the maximum number of goroutines that simultaneously interact with the storage.
	MaxDBConnections int
}
