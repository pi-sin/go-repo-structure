package domain

import "context"

// Author ...
type Author struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// AuthorManager represent the author's usecases
type AuthorManager interface {
	GetByID(ctx context.Context, id int64) (Author, error)
}

// AuthorRepository represent the author's repository contract
type AuthorRepository interface {
	GetByID(ctx context.Context, id int64) (Author, error)
}
