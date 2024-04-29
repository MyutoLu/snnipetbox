package models

import (
	"database/sql"
	"errors"
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
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, err
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `select ID, TITLE, CONTENT, CREATED, EXPIRES FROM snippets
				where expires > UTC_TIMESTAMP() order by id desc limit 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	Snippets := []*Snippet{}
	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		Snippets = append(Snippets, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return Snippets, nil
}
