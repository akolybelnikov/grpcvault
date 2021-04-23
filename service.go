// Package vault provides a gRPC service to hash and validate hashed passwords
package vault

import "context"

type vaultService struct {

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
	return "", nil
}

func (vaultService) Validate(ctx context.Context, password, hash string) (bool, error)  {
	return false, nil
}