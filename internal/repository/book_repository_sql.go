package repository

import (
	"context"
	"database/sql"

	"api/model"
)

type BookRepositorySQL struct {
	db *sql.DB
}

func NewBookRepositorySQL(db *sql.DB) BookRepository {
	return &BookRepositorySQL{db: db}
}

/* ---------- CREATE ---------- */

func (r *BookRepositorySQL) Create(ctx context.Context, b *model.Book) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO books (title, author, years, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id
	`,
		b.Title, b.Author, b.Years, b.CreatedAt, b.UpdatedAt,
	).Scan(&b.Id)
}

/* ---------- GET BY ID ---------- */

func (r *BookRepositorySQL) GetByID(ctx context.Context, id int) (*model.Book, error) {
	var b model.Book
	err := r.db.QueryRowContext(ctx, `
		SELECT id, title, author, years, created_at, updated_at
		FROM books WHERE id=$1
	`, id).Scan(
		&b.Id,
		&b.Title,
		&b.Author,
		&b.Years,
		&b.CreatedAt,
		&b.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

/* ---------- LIST ---------- */

func (r *BookRepositorySQL) List(ctx context.Context) ([]model.Book, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, author, years, created_at, updated_at
		FROM books
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := []model.Book{}
	for rows.Next() {
		var b model.Book
		if err := rows.Scan(
			&b.Id,
			&b.Title,
			&b.Author,
			&b.Years,
			&b.CreatedAt,
			&b.UpdatedAt,
		); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}

/* ---------- UPDATE ---------- */

func (r *BookRepositorySQL) Update(ctx context.Context, b *model.Book) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE books
		SET title=$1, author=$2, years=$3, updated_at=$4
		WHERE id=$5
	`,
		b.Title, b.Author, b.Years, b.UpdatedAt, b.Id,
	)
	return err
}

/* ---------- DELETE ---------- */

func (r *BookRepositorySQL) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `
		DELETE FROM books WHERE id=$1
	`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
