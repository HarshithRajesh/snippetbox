package postgresql

import (
	"database/sql"

	"github.com/HarshithRajesh/snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	// stmt := `INSERT INTO snippets (title, content, created, expires)
	// VALUES($1, $2, CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Kolkata', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Kolkata' + INTERVAL $3 DAY)
	// RETURNING id;`
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES($1, $2, CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Kolkata', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Kolkata' + $3 * INTERVAL '1 DAY')
	RETURNING id;`
	var lastInsertedID int
	err := m.DB.QueryRow(stmt, title, content, expires).Scan(&lastInsertedID)
	if err != nil {
		return 0, err
	}

	return lastInsertedID, nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {

	stmt := `SELECT id, title, content, created, expires 
			FROM snippets
			WHERE expires > CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Kolkata' AND id = $1`
	row := m.DB.QueryRow(stmt, id)
	s := &models.Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {

	stmt := `SELECT id, title, content, created, expires FROM snippets
			WHERE expires > CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Kolkata'
			 ORDER BY created DESC LIMIT 10;`
	row, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer row.Close()
	snippets := []*models.Snippet{}

	for row.Next() {
		s := &models.Snippet{}
		err = row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err = row.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
