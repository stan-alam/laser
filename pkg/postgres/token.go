package postgres

import (
	"database/sql"
	"log"

	"github.com/Power9-Alpha/laser"
)

type Token struct {
	db     *sql.DB
	create *sql.Stmt
	get    *sql.Stmt
	delete *sql.Stmt
}

func (t *Token) Init(db *sql.DB) error {
	t.db = db
	var err error
	if t.create, err = t.db.Prepare(SQLTokenInsert); err != nil {
		log.Printf("failed to create prepared statement: %s", err)
		return err
	} else if t.get, err = t.db.Prepare(SQLTokenSelect); err != nil {
		log.Printf("failed to create prepared statement: %s", err)
		return err
	} else if t.delete, err = t.db.Prepare(SQLTokenDelete); err != nil {
		log.Printf("failed to create prepared statement: %s", err)
		return err
	}
	return err
}

func (t *Token) Insert(username string) (string, error) {
	var id string
	err := t.create.QueryRow(username).Scan(&id)
	return id, err
}

func (t *Token) Select(username string) (*laser.Token, error) {
	token := laser.Token{}
	if err := t.get.QueryRow(username).Scan(&token.ID, &token.Username, &token.Created); err == sql.ErrNoRows {
		return nil, nil // @note: caller must check for nil model to address 404
	} else if err != nil {
		log.Printf("Select Failed: %s (%s)", username, err)
		return nil, err
	}
	return &token, nil
}

func (t *Token) Delete(id string) error {
	_, err := t.delete.Exec(id)
	return err
}
