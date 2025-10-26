package handler

import "strconv"

// parseIntQueryParam parses an integer query parameter
func parseIntQueryParam(s string) (int, error) {
	return strconv.Atoi(s)
}
