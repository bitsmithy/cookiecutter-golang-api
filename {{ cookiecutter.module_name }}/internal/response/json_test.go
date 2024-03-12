package response_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"slices"
	"testing"

	"{{ cookiecutter.module_path }}/internal/response"
)

type Foo struct {
	Num  int
	Str  string
	List []string
}

func TestJSONWithHeaders(t *testing.T) {
	testCases := []struct {
		name    string
		data    Foo
		want    string
		status  int
		headers http.Header
	}{
		{
			"happy path",
			Foo{Num: 2, Str: "foo", List: []string{"foo", "bar", "baz"}},
			`{"Num": 2, "Str": "foo", "List": ["foo", "bar", "baz"]}`,
			200,
			map[string][]string{"Foo": {"Bar", "Baz"}},
		},
	}

	for _, tc := range testCases {
		w := httptest.NewRecorder()

		err := response.JSONWithHeaders(w, tc.status, tc.data, tc.headers)
		if err != nil {
			t.Fatalf("JSONWithHeaders (test case: %s) errored: %q", tc.name, err)
		}

		resp := w.Result()
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("JSONWithHeaders (test case: %s) failed to read response: %q", tc.name, err)
		}

		var got interface{}
		var want interface{}

		err = json.Unmarshal(body, &got)
		if err != nil {
			t.Fatalf("JSONWithHeaders (test case: %s) failed to unmarshal into `got`: %q", tc.name, err)
		}

		err = json.Unmarshal([]byte(tc.want), &want)
		if err != nil {
			t.Fatalf("JSONWithHeaders (test case: %s) failed to unmarshal into `got`: %q", tc.name, err)
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("JSONWithHeaders (test case: %s) result does not match | want: %+v, got %+v", tc.name, want, got)
		}

		gotStatus := resp.StatusCode
		if gotStatus != tc.status {
			t.Errorf("JSONWithHeaders (test case: %s) status does not match | want: %d, got %d", tc.name, tc.status, gotStatus)
		}

		gotHeader := resp.Header
		if gotHeader.Get("Content-Type") != "application/json" {
			t.Errorf("JSONWithHeaders (test case: %s) did not have a Content-Type of application/json", tc.name)
		}

		for k, v := range tc.headers {
			headerVal := gotHeader.Get(k)
			if !slices.Contains(v, headerVal) {
				t.Errorf("JSONWithHeaders (test case: %s) did not match a header for %s | got: %q, want: %q", tc.name, k, headerVal, v)
			}
		}
	}
}

func TestJSON(t *testing.T) {
	testCases := []struct {
		name    string
		data    Foo
		want    string
		status  int
		headers http.Header
	}{
		{
			"happy path",
			Foo{Num: 2, Str: "foo", List: []string{"foo", "bar", "baz"}},
			`{"Num": 2, "Str": "foo", "List": ["foo", "bar", "baz"]}`,
			200,
			map[string][]string{"Foo": {"Bar", "Baz"}},
		},
	}

	for _, tc := range testCases {
		w := httptest.NewRecorder()

		err := response.JSON(w, tc.status, tc.data)
		if err != nil {
			t.Fatalf("JSON (test case: %s) errored: %q", tc.name, err)
		}

		resp := w.Result()
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("JSON (test case: %s) failed to read response: %q", tc.name, err)
		}

		var got interface{}
		var want interface{}

		err = json.Unmarshal(body, &got)
		if err != nil {
			t.Fatalf("JSON (test case: %s) failed to unmarshal into `got`: %q", tc.name, err)
		}

		err = json.Unmarshal([]byte(tc.want), &want)
		if err != nil {
			t.Fatalf("JSON (test case: %s) failed to unmarshal into `got`: %q", tc.name, err)
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("JSON (test case: %s) result does not match | want: %+v, got %+v", tc.name, want, got)
		}

		gotStatus := resp.StatusCode
		if gotStatus != tc.status {
			t.Errorf("JSON (test case: %s) status does not match | want: %d, got %d", tc.name, tc.status, gotStatus)
		}

		gotHeader := resp.Header
		if gotHeader.Get("Content-Type") != "application/json" {
			t.Errorf("JSON (test case: %s) did not have a Content-Type of application/json", tc.name)
		}

		for k, v := range tc.headers {
			headerVal := gotHeader.Get(k)
			if slices.Contains(v, headerVal) {
				t.Errorf("JSON (test case: %s) did matched a header for %s, but should not have any headers | got: %q, want: %q", tc.name, k, headerVal, v)
			}
		}
	}
}
