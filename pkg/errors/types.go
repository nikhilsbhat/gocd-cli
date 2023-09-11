package errors

type AuthError struct {
	Message string
}

type MoreArgError struct {
	Message string
}

type CipherMinKeyError struct {
	Message string
}

type UnknownObjectTypeError struct {
	Name string
}

type ConfigRepoError struct {
	Message string
}

type MaterialError struct {
	Message string
}
