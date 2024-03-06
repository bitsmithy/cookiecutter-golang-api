package request_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"testing/iotest"

	"{{ cookiecutter.module_path }}/internal/request"
)

type Foo struct {
	Num  int
	Str  string
	List []string
}

func TestDecodeJSON(t *testing.T) {
	testCases := []struct {
		name string
		body io.Reader
		want Foo
		err  error
	}{
		{
			"happy path",
			strings.NewReader(`{"num": 2, "str": "foo", "list": ["bar", "baz"]}`),
			Foo{Num: 2, Str: "foo", List: []string{"bar", "baz"}},
			nil,
		},
		{
			"multiple json",
			strings.NewReader(`{"num": 2, "str": "foo", "list": ["bar", "baz"]}{"foo": "bar"}`),
			Foo{Num: 2, Str: "foo", List: []string{"bar", "baz"}},
			request.MultipleJSONError{},
		},
		{
			"body empty",
			strings.NewReader(``),
			Foo{},
			request.BodyEmptyError{},
		},
		{
			"malformed with position",
			strings.NewReader(`{"foo":}`),
			Foo{},
			request.MalformedJSONError{Position: 8},
		},
		{
			"malformed without position",
			iotest.ErrReader(io.ErrUnexpectedEOF),
			Foo{},
			request.MalformedJSONError{Position: -1},
		},
		{
			"incorrect type",
			strings.NewReader(`{"num": "2"}`),
			Foo{},
			request.IncorrectJSONTypeError{Field: "Num", Position: 11},
		},
		{
			"unknown key",
			strings.NewReader(`{"blah": "blah"}`),
			Foo{},
			nil,
		},
		{
			"large size",
			strings.NewReader(fmt.Sprintf(`{"foo": "%s"}`, bytes.Repeat([]byte("x"), request.MAX_JSON_BYTES-10))), // -10 because the rest of the JSON body is 10 bytes
			Foo{},
			request.BodyTooLargeError{ByteSize: request.MAX_JSON_BYTES},
		},
	}

	for _, tc := range testCases {
		var got Foo
		r, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, "http://example.com", tc.body)
		if err != nil {
			t.Fatalf("Could not create request: %s", err)
		}
		w := httptest.NewRecorder()

		err = request.DecodeJSON(w, r, &got)
		if !errors.Is(err, tc.err) {
			t.Fatalf("DecodeJSON (test case: %s) error does not match | want: %q, got %q", tc.name, tc.err, err)
		}

		if !reflect.DeepEqual(tc.want, got) {
			t.Errorf("DecodeJSON (test case: %s) result does not match | want: %+v, got %+v", tc.name, tc.want, got)
		}
	}
}

func TestDecodeJSONStrict(t *testing.T) {
	testCases := []struct {
		name string
		body io.Reader
		want Foo
		err  error
	}{
		{
			"happy path",
			strings.NewReader(`{"num": 2, "str": "foo", "list": ["bar", "baz"]}`),
			Foo{Num: 2, Str: "foo", List: []string{"bar", "baz"}},
			nil,
		},
		{
			"multiple json",
			strings.NewReader(`{"num": 2, "str": "foo", "list": ["bar", "baz"]}{"foo": "bar"}`),
			Foo{Num: 2, Str: "foo", List: []string{"bar", "baz"}},
			request.MultipleJSONError{},
		},
		{
			"body empty",
			strings.NewReader(``),
			Foo{},
			request.BodyEmptyError{},
		},
		{
			"malformed with position",
			strings.NewReader(`{"foo":}`),
			Foo{},
			request.MalformedJSONError{Position: 8},
		},
		{
			"malformed without position",
			iotest.ErrReader(io.ErrUnexpectedEOF),
			Foo{},
			request.MalformedJSONError{Position: -1},
		},
		{
			"incorrect type",
			strings.NewReader(`{"num": "2"}`),
			Foo{},
			request.IncorrectJSONTypeError{Field: "Num", Position: 11},
		},
		{
			"unknown key",
			strings.NewReader(`{"blah": "blah"}`),
			Foo{},
			request.UnknownJSONKeyError{Key: "blah"},
		},
		{
			"large size",
			strings.NewReader(fmt.Sprintf(`{"foo": "%s"}`, bytes.Repeat([]byte("x"), request.MAX_JSON_BYTES-10))), // -10 because the rest of the JSON body is 10 bytes
			Foo{},
			request.BodyTooLargeError{ByteSize: request.MAX_JSON_BYTES},
		},
	}

	for _, tc := range testCases {
		var got Foo
		r, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, "http://example.com", tc.body)
		if err != nil {
			t.Fatalf("Could not create request: %s", err)
		}
		w := httptest.NewRecorder()

		err = request.DecodeJSONStrict(w, r, &got)
		if !errors.Is(err, tc.err) {
			t.Fatalf("DecodeJSONStrict (test case: %s) error does not match | want: %q, got %q", tc.name, tc.err, err)
		}

		if !reflect.DeepEqual(tc.want, got) {
			t.Errorf("DecodeJSONStrict (test case: %s) result does not match | want: %+v, got %+v", tc.name, tc.want, got)
		}
	}
}
