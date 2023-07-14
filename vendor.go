package vendor4git

import (
	"context"
)

// Vendor creates and administers git repositories.
type Vendor interface {
	CreateRepo(ctx context.Context, name string, org string, private bool) (*Repository, error)
	RemoveRepo(ctx context.Context, name string, org string) error
}

type Repository struct {
	HTTPSURL *string `json:"https_url,omitempty"` // e.g. https://github.com/user/repo.git
	SSHURL   *string `json:"ssh_url,omitempty"`   // e.g. git@github.com:user/repo.git
}
