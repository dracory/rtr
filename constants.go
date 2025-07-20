package rtr

// contextKey is a type for context keys to avoid collisions
type contextKey string

// ParamsKey is the key used to store path parameters in the request context
// Using a more specific key to avoid collisions with other packages
// ParamsKey is the key used to store path parameters in the request context
// Using a more specific key to avoid collisions with other packages
const ParamsKey contextKey = "rtr.path.params"

// ExecutionSequenceKey is used to track the execution sequence of middlewares in tests
const ExecutionSequenceKey contextKey = "rtr.execution.sequence"
