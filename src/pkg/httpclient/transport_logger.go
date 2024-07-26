package httpclient

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
)

type LoggingTransport struct{}

func (s *LoggingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	bytes, _ := httputil.DumpRequestOut(r, true)

	resp, err := http.DefaultTransport.RoundTrip(r)
	// err is returned after dumping the response

	respBytes, _ := httputil.DumpResponse(resp, true)
	bytes = append(bytes, respBytes...)

	slog.Debug("loggingTransport.round_trip", slog.Any("resp", string(bytes)))

	return resp, err
}
