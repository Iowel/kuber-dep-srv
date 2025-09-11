package repository

import (
	"context"
	"database/sql"
)

type Repository interface {
	Save(ctx context.Context)
	Delete(ctx context.Context)
	Get(ctx context.Context)
	Update(ctx context.Context)
}

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) Save(ctx context.Context) {}

func (r *Repo) Delete(ctx context.Context) {}

func (r *Repo) Get(ctx context.Context) {}

func (r *Repo) Update(ctx context.Context) {}
