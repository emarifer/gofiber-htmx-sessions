package models

import (
	"time"

	"github.com/emarifer/gofiber-htmx-sessions/db"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID        string    `json:"id,omitempty"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type UserSession struct {
	SID    string `json:"sid"`
	IP     string `json:"ip"`
	Login  string `json:"login"`
	Expiry string `json:"expiry"`
	UA     string `json:"ua"`
}

type Account struct {
	Email    string        `json:"email"`
	Username string        `json:"username"`
	Session  string        `json:"userSession"`
	Sessions []UserSession `json:"sessions"`
}

func (u *User) CreateUser() (User, error) {
	query := `INSERT INTO users (id, email, password, username, created_at)
		VALUES(?, ?, ?, ?, ?) RETURNING *`

	stmt, err := db.Db.Prepare(query)
	if err != nil {
		return User{}, err
	}

	defer stmt.Close()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)
	if err != nil {
		return User{}, err
	}

	var newUser User
	err = stmt.QueryRow(
		uuid.NewString(),
		u.Email,
		string(hashedPassword),
		u.Username,
		time.Now().UTC(),
	).Scan(
		&newUser.ID,
		&newUser.Email,
		&newUser.Password,
		&newUser.Username,
		&newUser.CreatedAt,
	)
	if err != nil {
		return User{}, err
	}

	/* if i, err := result.RowsAffected(); err != nil || i != 1 {
		return errors.New("error: an affected row was expected")
	} */

	return newUser, nil
}

func (n *User) GetUserById() (User, error) {
	query := `SELECT * FROM users
		WHERE id=?`

	stmt, err := db.Db.Prepare(query)
	if err != nil {
		return User{}, err
	}

	defer stmt.Close()

	var recoveredUser User
	err = stmt.QueryRow(
		n.ID,
	).Scan(
		&recoveredUser.ID,
		&recoveredUser.Email,
		&recoveredUser.Password,
		&recoveredUser.Username,
		&recoveredUser.CreatedAt,
	)
	if err != nil {
		return User{}, err
	}

	return recoveredUser, nil
}

func (n *User) GetUserByEmail() (User, error) {
	query := `SELECT * FROM users
		WHERE email=?`

	stmt, err := db.Db.Prepare(query)
	if err != nil {
		return User{}, err
	}

	defer stmt.Close()

	var recoveredUser User
	err = stmt.QueryRow(
		n.Email,
	).Scan(
		&recoveredUser.ID,
		&recoveredUser.Email,
		&recoveredUser.Password,
		&recoveredUser.Username,
		&recoveredUser.CreatedAt,
	)
	if err != nil {
		return User{}, err
	}

	return recoveredUser, nil
}
