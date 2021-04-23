package grpcvault

import "net/http"

import (
	httptransport "github.com/go-kit/kit/transport/http"
)

// NewHTTPServer makes a new Vault HTTP service.
func NewHTTPServer(endpoints Endpoints) http.Handler {
	m := http.NewServeMux()
	m.Handle("/hash", httptransport.NewServer(
		endpoints.HashEndpoint,
		decodeHashRequest,
		encodeResponse,
	))
	m.Handle("/validate", httptransport.NewServer(
		endpoints.ValidateEndpoint,
		decodeValidateRequest,
		encodeResponse,
	))
	return m
}
