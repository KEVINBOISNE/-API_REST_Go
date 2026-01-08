package model

import (
	"os"
	"io"
	"encoding/json"
	"html/template"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Book struct {
	Name  string
	Author string
	Years int
	CreatedAt string
	UpdatedAt string
}

// booker := Book{"kevin_hart_good_boy", "Kevin_hart", 2026, 7/1/2025, 8/1/2025}
// tmpl, err := template.New("test").Parse("{{.Name}} items are made of {{.Author}} {{.Years}} {{.Created_at}} {{.Updated_at}}")
// if err != nil { panic(err) }
// err = tmpl.Execute(os.Stdout, booker)
// if err != nil { panic(err) }

func New(name string, author string, years int, created_at string, updated_at string) (*Book, error) {

const tmpl = "Nom : {{ .Name }}. Auteur : {{ .Author }}. Annee : {{ .Years }}. Creation : {{ .CreatedAt }}. modification : {{ .UpdatedAt }}"

p := &Book{Name: name,Author: author,Years: years,CreatedAt: created_at,UpdatedAt: updated_at}

t, err := template.New("tmpl").Parse(tmpl)

if err != nil {

return nil, err

}

err = t.Execute(os.Stdout, p)

if err != nil {

return nil, err

}
return p, nil 
}

func CreateBookFromRequest(b *Book, r io.Reader) (*Book, error) {
	var req Book
	if err := json.NewDecoder(r).Decode(&req); err != nil {
		return nil, err
	}

	b.Name = req.Name
	b.Author = req.Author
	b.Years = req.Years
	b.CreatedAt = req.CreatedAt
	b.UpdatedAt = req.UpdatedAt

	const tmplStr = "Nom : {{ .Name }}. Auteur : {{ .Author }}. Annee : {{ .Years }}. Creation : {{ .CreatedAt }}. Modification : {{ .UpdatedAt }}\n"
	tmpl, err := template.New("book").Parse(tmplStr)
	if err != nil {
		return nil, err
	}

	if err := tmpl.Execute(os.Stdout, b); err != nil {
		return nil, err
	}

	return b, nil
}
