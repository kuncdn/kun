package util

// IsSuccessCode .
func IsSuccessCode(statusCode int) bool {
	return statusCode < 500
}
