package errors

import "fmt"

type AuthError struct {
	Message string
}

type MoreArgError struct {
	Message string
}

func (e MoreArgError) Error() string {
	return fmt.Sprintf("args cannot be more than one, only one %s text must be passed", e.Message)
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

func (e *ConfigRepoError) Error() string {
	return fmt.Sprintf("args cannot be more than one. At once, %v", e.Message)
}

func (e *AuthError) Error() string {
	return e.Message
}

func (e *CipherMinKeyError) Error() string {
	return e.Message
}

func (e *UnknownObjectTypeError) Error() string {
	return fmt.Sprintf("unknow object type %s", e.Name)
}
