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
	var id int
	var hashed_password []byte
	row := m.DB.QueryRow("SELECT id,hashed_password FROM users WHERE email=$1",email)
	err := row.Scan(&id,&hashed_password)
	if err == sql.ErrNoRows{
		return 0,models.ErrInvalidCredentials
		}else if err!=nil{
			return 0,err
		}
	err = bcrypt.CompareHashAndPassword(hashed_password,[]byte(password))

	if err == bcrypt.ErrMismatchedHashAndPassword{
		return 0 , models.ErrInvalidCredentials
	}else if err != nil{
		return 0 , err
	}
	return id,nil
}

func (m *UserModel) Get(id int) (*models.Users, error) {
	return nil, nil
}
