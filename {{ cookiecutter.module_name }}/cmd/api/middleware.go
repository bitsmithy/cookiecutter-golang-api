package main

import (
	"fmt"
	"net/http"
	"time"
)

type (
	// Struct to capture data with the loggingResponseWriter
	capturedInfo struct {
		status int
		size   int
	}

	// Struct to capture response data then forward back to original responseWriter
	loggingResponseWriter struct {
		http.ResponseWriter
		capturedInfo *capturedInfo
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b) // write response using original http.ResponseWriter
	r.capturedInfo.size += size            // capture size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode) // write status code using original http.ResponseWriter
	r.capturedInfo.status = statusCode       // capture status code
}

func (s *server) logRequest(original http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseData := &capturedInfo{}
		lw := loggingResponseWriter{
			ResponseWriter: w,            // original writer
			capturedInfo:   responseData, // map to capture response data in
		}

		start := time.Now()
		original.ServeHTTP(&lw, r)
		duration := time.Since(start)

		s.logger.Info("request received",
			"method", r.Method,
			"path", r.URL.Path,
			"respCode", fmt.Sprint(responseData.status),
			"duration", fmt.Sprintf("%dµs", duration.Microseconds()),
			"respSize", fmt.Sprintf("%d bytes", responseData.size),
		)
	})
}

func (s *server) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event of a panic
		// as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a panic or
			// not.
			if err := recover(); err != nil {
				// If there was a panic, set a "Connection: close" header on the
				// response. This acts as a trigger to make Go's HTTP server
				// automatically close the current connection after a response has been
				// sent.
				w.Header().Set("Connection", "close")
				// The value returned by recover() has the type any, so we use
				// fmt.Errorf() to normalize it into an error and call our
				// serverErrorResponse() helper. In turn, this will log the error using
				// our custom Logger type at the ERROR level and send the client a 500
				// Internal Server Error response.
				s.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
