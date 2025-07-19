package middlewares

// DefaultMiddlewares returns a slice of default middlewares that should be used with the router.
// Currently, it only includes the RecoveryMiddleware.
func DefaultMiddlewares() []Middleware {
	return []Middleware{
		RecoveryMiddleware,
	}
}
