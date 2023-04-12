package jsonapi

// ErrorResponse is a common struct for error response
type ErrorResponse struct {
	Errors []string `json:"errors"`
}
