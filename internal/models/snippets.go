package models

import (
	"database/sql"
	"time"
)

// Snippet is a representation of a snippet of code.
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires) VALUES (?, ?,UTC_TIMESTAMP(), 
             date_add(UTC_TIMESTAMP(),interval ? day ))`
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

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `select id, title, content, created, expires from snippets where expires > utc_timestamp() 
        	 and   id = ?`
	row := m.DB.QueryRow(stmt, id)
	s := &Snippet{}
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err == sql.ErrNoRows {
		return nil, ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return s, err
}

func (m *SnippetModel) Latest(id int) ([]*Snippet, error) {
	return nil, nil
}
