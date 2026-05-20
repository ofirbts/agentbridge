package observability

import "strings"

func ClassifyError(message string) string {
	lower := strings.ToLower(message)
	switch {
	case strings.Contains(lower, "429"), strings.Contains(lower, "rate limit"):
		return "rate_limit"
	case strings.Contains(lower, "timeout"), strings.Contains(lower, "deadline"):
		return "timeout"
	case strings.Contains(lower, "server error"), strings.Contains(lower, "status (5"):
		return "server_error"
	case strings.Contains(lower, "unexpected status (4"), strings.Contains(lower, "status (4"):
		return "client_error"
	case strings.Contains(lower, "request failed"), strings.Contains(lower, "connection"), strings.Contains(lower, "network"):
		return "network"
	case strings.Contains(lower, "mock"), strings.Contains(lower, "stub"), strings.Contains(lower, "provider"):
		return "provider"
	default:
		return "unknown"
	}
}
