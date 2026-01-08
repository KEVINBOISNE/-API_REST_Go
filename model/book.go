package model

import (
	"encoding/json"
	"errors"
	"io"
	"time"
)

type Book struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Years     int    `json:"years"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func CreateBookFromRequest(r io.Reader) (*Book, error) {
	var b Book

	if err := json.NewDecoder(r).Decode(&b); err != nil {
		return nil, err
	}

	// Validation simple
	if b.Title == "" || b.Author == "" {
		return nil, errors.New("title and author are required")
	}
	if b.Years < 1400 || b.Years > time.Now().Year()+1 {
		return nil, errors.New("invalid year")
	}

	return &b, nil
}

func UpdateBookFromRequest(b *Book, r io.Reader) error {
	var req Book
	if err := json.NewDecoder(r).Decode(&req); err != nil {
		return err
	}

	if req.Title != "" {
		b.Title = req.Title
	}
	if req.Author != "" {
		b.Author = req.Author
	}
	if req.Years != 0 {
		b.Years = req.Years
	}

	return nil
}
