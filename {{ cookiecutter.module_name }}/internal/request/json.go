package request

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

const MAX_JSON_BYTES = 1_048_576

func DecodeJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	return decodeJSON(w, r, dst, false)
}

func DecodeJSONStrict(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	return decodeJSON(w, r, dst, true)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst interface{}, disallowUnknownFields bool) error {
	r.Body = http.MaxBytesReader(w, r.Body, int64(MAX_JSON_BYTES))

	dec := json.NewDecoder(r.Body)

	if disallowUnknownFields {
		dec.DisallowUnknownFields()
	}

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return MalformedJSONError{Position: int(syntaxError.Offset)}

		case errors.Is(err, io.ErrUnexpectedEOF):
			return MalformedJSONError{Position: -1}

		case errors.As(err, &unmarshalTypeError):
			return IncorrectJSONTypeError{Field: unmarshalTypeError.Field, Position: unmarshalTypeError.Offset}

		case errors.Is(err, io.EOF):
			return BodyEmptyError{}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return UnknownJSONKeyError{Key: strings.Trim(fieldName, "\"")}

		case err.Error() == "http: request body too large":
			return BodyTooLargeError{ByteSize: MAX_JSON_BYTES}

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return MultipleJSONError{}
	}

	return nil
}
