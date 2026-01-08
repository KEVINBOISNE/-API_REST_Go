package repository

import (
	"context"
	"api/model"
)

type BookRepository interface {
	Create(ctx context.Context, b *model.Book) error
	GetByID(ctx context.Context, id int) (*model.Book, error)
	List(ctx context.Context) ([]model.Book, error)
	Update(ctx context.Context, b *model.Book) error
	Delete(ctx context.Context, id int) error
}
