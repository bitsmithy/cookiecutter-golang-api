package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	th "{{ cookiecutter.module_path }}/internal/testhelper"
)

func testRequest(t *testing.T, method, path string, params io.Reader) (*httptest.ResponseRecorder, map[string]any) {
	t.Helper()

	app := &application{}
	mux := app.mux()

	w := httptest.NewRecorder()
	r, err := http.NewRequestWithContext(context.TODO(), method, path, params)
	if err != nil {
		t.Fatalf("Could not create request: %s", err)
	}

	r.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(w, r)

	resp := make(map[string]any)
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Could not unmarshal response body: %s", err)
	}

	return w, resp
}

func TestAPI(t *testing.T) {
	t.Run("/status", func(t *testing.T) {
		want := map[string]any{"Status": "OK"}
		_, got := testRequest(t, "GET", "/status", nil)

		if !th.MapsEq(want, got) {
			t.Errorf("API handler response body did not match | want: %+v, got %+v", want, got)
		}
	})

	t.Run("GenericErrors", func(t *testing.T) {
		testCases := []struct {
			method   string
			path     string
			params   io.Reader
			want     map[string]any
			wantCode int
		}{
			{"GET", "/foobar", nil, map[string]any{"Error": "The requested resource could not be found"}, 404},
			{"POST", "/status", nil, map[string]any{"Error": "The POST method is not supported for this resource"}, 405},
		}

		for _, tc := range testCases {
			w, got := testRequest(t, tc.method, tc.path, tc.params)

			if tc.wantCode != w.Code {
				t.Errorf("API handler response code did not match | got: %d, want: %d", w.Code, tc.wantCode)
			}

			if !th.MapsEq(tc.want, got) {
				t.Errorf("API handler response body did not match | got: %+v, want %+v", got, tc.want)
			}
		}
	})
}
