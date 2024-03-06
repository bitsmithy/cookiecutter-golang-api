package request

import "fmt"

type MultipleJSONError struct{}

func (e MultipleJSONError) Error() string {
	return "body must only contain a single JSON value"
}

type BodyEmptyError struct{}

func (e BodyEmptyError) Error() string {
	return "body must not be empty"
}

type MalformedJSONError struct {
	Position int
}

func (e MalformedJSONError) Error() string {
	if e.Position != -1 {
		return fmt.Sprintf("body contains badly-formed JSON (at character %d)", e.Position)
	} else {
		return "body contains badly-formed JSON"
	}
}

type IncorrectJSONTypeError struct {
	Field    string
	Position int64
}

func (e IncorrectJSONTypeError) Error() string {
	if e.Field == "" {
		return fmt.Sprintf("body contains incorrect JSON type at character %d", e.Position)
	} else {
		return fmt.Sprintf("body contains incorrect JSON type for field %q (at character %d)", e.Field, e.Position)
	}
}

type BodyTooLargeError struct {
	ByteSize int
}

func (e BodyTooLargeError) Error() string {
	return fmt.Sprintf("body must not be larger than %d bytes", e.ByteSize)
}

type UnknownJSONKeyError struct {
	Key string
}

func (e UnknownJSONKeyError) Error() string {
	return fmt.Sprintf("body contains unknown key %q", e.Key)
}
