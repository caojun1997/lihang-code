package ctxkey

type ContextKey string

const (
	KeyRequestID ContextKey = "request_id"
	KeyUserId   ContextKey = "user_id"
	KeyUsername ContextKey = "username"
)
