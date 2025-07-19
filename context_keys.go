package rtr

// contextKey is a type for context keys to avoid collisions
type contextKey string

// Context key for storing path parameters in the request context
const (
	// ParamsKey is the key used to store path parameters in the request context
	// Using a more specific key to avoid collisions with other packages
	ParamsKey contextKey = "rtr.path.params"
)
