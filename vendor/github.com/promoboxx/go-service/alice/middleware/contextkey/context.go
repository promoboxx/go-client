package contextkey

type ContextKey int

// context keys
const (
	ContextKeyLogger ContextKey = iota
	ContextKeyRequestID
	ContextKeyInsecureUserID
	ContextKeyClaims
	ContextKeyDBConn
	ContextKeyProducer
	ContextKeyJWT
	ContextKeyDB
	ContextKeyCanary
)
