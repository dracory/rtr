package rtr

import (
	"maps"
	"net/http"
)

// GetParam retrieves a path parameter from the request context by name.
// Returns the parameter value and true if found, or an empty string and false otherwise.
func GetParam(r *http.Request, name string) (string, bool) {
	if r == nil {
		return "", false
	}

	// Get the params map from the context
	params, ok := r.Context().Value(ParamsKey).(map[string]string)
	if !ok || params == nil {
		return "", false
	}

	// Look up the parameter
	value, exists := params[name]
	return value, exists
}

// MustGetParam retrieves a path parameter from the request context by name.
// Panics if the parameter is not found. Use only when you're certain the parameter exists.
func MustGetParam(r *http.Request, name string) string {
	value, exists := GetParam(r, name)
	if !exists {
		panic("parameter not found: " + name)
	}
	return value
}

// GetParams returns all path parameters as a map.
// Returns an empty map if no parameters exist.
func GetParams(r *http.Request) map[string]string {
	if r == nil {
		return map[string]string{}
	}

	params, ok := r.Context().Value(ParamsKey).(map[string]string)
	if !ok || params == nil {
		return map[string]string{}
	}

	// Return a copy to prevent external modifications
	result := make(map[string]string, len(params))
	maps.Copy(result, params)
	return result
}
