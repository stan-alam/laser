package postgres

import (
	"database/sql"
	"log"

	"github.com/Power9-Alpha/laser"
)

type User struct {
	db     *sql.DB
	create *sql.Stmt
	getOne *sql.Stmt
	getAll *sql.Stmt
	update *sql.Stmt
	delete *sql.Stmt
	login  *sql.Stmt
}

func (u *User) Init(db *sql.DB) error {
	u.db = db
	var err error
	if u.create, err = u.db.Prepare(SQLUserInsert); err != nil {
		log.Printf("failed to create prepared statement: %s", err)
		return err
	} else if u.getOne, err = u.db.Prepare(SQLUserSelect); err != nil {
		log.Printf("failed to create prepared statement: %s", err)
		return err
	} else if u.getAll, err = u.db.Prepare(SQLUsersSelect); err != nil {
		log.Printf("failed to create prepared statement: %s", err)
		return err
	} else if u.update, err = u.db.Prepare(SQLUserUpdate); err != nil {
		log.Printf("failed to create prepared statement: %s", err)
		return err
	} else if u.delete, err = u.db.Prepare(SQLUserDelete); err != nil {
		log.Printf("failed to create prepared statement: %s", err)
		return err
	} else if u.login, err = u.db.Prepare(SQLUserLogin); err != nil {
		log.Printf("failed to create prepared statement: %s", err)
		return err
	}
	return err
}

func (u *User) Login(email, password string) (*laser.User, error) {
	user := laser.User{}
	if err := u.login.QueryRow(email, password).Scan(&user.Name, &user.Email, &user.Password); err != nil {
		log.Printf("Login Failed: %s (%s)", email, err)
		return nil, err
	}
	return &user, nil
}

func (u *User) Insert(user *laser.User) error {
	_, err := u.create.Exec(user.Name, user.Email, user.Password)
	// @note: unsure if result check is needed for anything here
	return err
}

func (u *User) SelectOne(id string) (*laser.User, error) {
	user := laser.User{}
	if err := u.getOne.QueryRow(id).Scan(&user.Name, &user.Email, &user.Password); err == sql.ErrNoRows {
		return nil, nil // @note: caller must check for nil model to address 404
	} else if err != nil {
		log.Printf("Select Failed: %s (%s)", id, err)
		return nil, err
	}
	return &user, nil
}

func (u *User) Select() ([]laser.User, error) {
	rows, err := u.getAll.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []laser.User{}

	for rows.Next() {
		user := laser.User{}
		if err := rows.Scan(&user.Name, &user.Email, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (u *User) Update(user *laser.User) error {
	_, err := u.update.Exec(user.Name, user.Email, user.Password)
	return err
}

func (u *User) Delete(id string) error {
	_, err := u.delete.Exec(id)
	return err
}
