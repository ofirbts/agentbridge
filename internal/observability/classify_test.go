package observability

import "testing"

func TestClassifyError(t *testing.T) {
	cases := []struct {
		msg  string
		want string
	}{
		{"http crawl: rate limited (429)", "rate_limit"},
		{"search http: timeout status (504)", "timeout"},
		{"http crawl: server error (502)", "server_error"},
		{"http crawl: unexpected status (404)", "client_error"},
		{"search http: request failed: dial tcp", "network"},
		{"mock search failure", "provider"},
		{"something else", "unknown"},
	}
	for _, tc := range cases {
		if got := ClassifyError(tc.msg); got != tc.want {
			t.Fatalf("ClassifyError(%q) = %q, want %q", tc.msg, got, tc.want)
		}
	}
}
