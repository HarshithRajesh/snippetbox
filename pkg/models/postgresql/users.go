package postgresql

import (
	"database/sql"
	"strings"

	"github.com/HarshithRajesh/snippetbox/pkg/models"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil
	}
	stmt := `INSERT INTO users (name,email,hashed_password,created)
			VALUES ($1,$2,$3,CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Kolkata')`
	_, err = m.DB.Exec(stmt, name, email, string(hashed_password))
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "1062" && strings.Contains(pqErr.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
	}
	return err
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) Get(id int) (*models.Users, error) {
	return nil, nil
}
