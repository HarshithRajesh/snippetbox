package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: Invalid crediendtials")
	ErrDuplicateEmail     = errors.New("models: Email already exists")
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type Users struct {
	ID    			int
	Name  			string
	Email 			string
	HashedPassword	[]byte
	Created		    time.Time
}
