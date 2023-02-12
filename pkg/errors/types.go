package errors

type AuthError struct {
	Message string
}

type CipherMinKeyError struct {
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}

func (e *CipherMinKeyError) Error() string {
	return e.Message
}
