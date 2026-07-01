package core

type Validator[T any] interface {
	Validate(T) error
}
