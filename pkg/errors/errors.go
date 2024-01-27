package errors

import "fmt"

func (e MoreArgError) Error() string {
	return fmt.Sprintf("args cannot be more than one, only one %s text must be passed", e.Message)
}

func (e *ConfigRepoError) Error() string {
	return fmt.Sprintf("config-repo: %s", e.Message)
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

func (e *MaterialError) Error() string {
	return e.Message
}

func (e *CLIError) Error() string {
	return e.Message
}
