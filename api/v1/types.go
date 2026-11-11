// Package v1 holds the public wire types of the v1 HTTP API.
package v1

// Health is the response body of the health check endpoint.
type Health struct {
	Status string `json:"status"`
}
