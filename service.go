// Package vault provides a gRPC service to hash and validate hashed passwords
package vault

import (
	"context"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type vaultService struct {

}

type hashRequest struct {
	Password string `json:"password"`
}

type hashResponse struct {
	Hash string `json:"hash"`
	Err string `json:"err,omitempty"`
}

type validateRequest struct {
	Password string `json:"password"`
	Hash string `json:"hash"`
}

type validateResponse struct {
	Valid bool `json:"valid"`
	Err string `json:"err,omitempty"`
}

// NewService make sa new Service.
func NewService() Service  {
	return vaultService{}
}

// Service provides password hashing capabilities.
type Service interface {
	Hash(ctx context.Context, password string) (string, error)
	Validate(ctx context.Context, password, hash string) (bool, error)
}

func (vaultService) Hash(ctx context.Context, password string) (string, error)  {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (vaultService) Validate(ctx context.Context, password, hash string) (bool, error)  {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, nil
	}
	return true, nil
}

func decodeHashRequest(ctx context.Context, r *http.Request) (interface{}, error)  {
	var req hashRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeValidateRequest(ctx context.Context, r *http.Request) (interface{}, error)  {
	var req validateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error  {
	return json.NewEncoder(w).Encode(response)
}