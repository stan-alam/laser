package postgres

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// @note: properly abstracting database operations makes sense as we can
//        substitute other DBMS by matching interfaces, and we can control
//        the order of table creation to match dependencies by providing an
//        Init() function, which can return the standardized sql.DB.
func Init(connString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	} else if tx, err := db.Begin(); err != nil {
		log.Printf("failed to create transaction: %s", err)
		return nil, err
	} else if _, err = tx.Exec(SQLCreateExtensionPGCrypto); err != nil {
		log.Printf("failed to create pgcrypto extension: %s", err)
		tx.Rollback()
		return nil, err
	} else if _, err = tx.Exec(SQLCreateTableUsers); err != nil {
		log.Printf("failed to create users table: %s", err)
		tx.Rollback()
		return nil, err
	} else if _, err = tx.Exec(SQLCreateTableTokens); err != nil {
		log.Printf("failed to create users table: %s", err)
		tx.Rollback()
		return nil, err
	} else if _, err = tx.Exec(SQLCreateTableServices); err != nil {
		log.Printf("failed to create users table: %s", err)
		tx.Rollback()
		return nil, err
		// @todo: alters and preseeds can be added from here
	} else if err = tx.Commit(); err != nil {
		log.Printf("failed to commit transaction: %s", err)
		tx.Rollback()
		return nil, err
	}
	return db, nil
}
