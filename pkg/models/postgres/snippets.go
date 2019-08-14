package postgres

import (
	"contra-design.com/new_snippetbox/pkg/models"
	"database/sql"
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES($1, $2, NOW(), NOW() + INTERVAL '365 days')`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m SnippetModel) Get(id int) (*models.Snippet, error) {
	return nil, nil
}

func (m SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}