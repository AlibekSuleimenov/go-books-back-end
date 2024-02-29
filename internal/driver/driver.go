package driver

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"time"
)

type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 5
const maxIdleDbConn = 5
const maxDbLifeTime = 5 * time.Minute

func ConnectPostgres(dsn string) (*DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenDbConn)
	db.SetMaxIdleConns(maxIdleDbConn)
	db.SetConnMaxLifetime(maxDbLifeTime)

	err = testDB(db)
	if err != nil {
		return nil, err
	}

	dbConn.SQL = db
	return dbConn, nil
}

func testDB(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}

	fmt.Println("*** Pinged database successfully! ***")
	return nil
}
