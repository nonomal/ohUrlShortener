package db

import (
	"database/sql"
	"fmt"
	"ohurlshortener/utils"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var dbService = &DatabaseService{}

type DatabaseService struct {
	Connection *sqlx.DB
}

func InitDatabaseService() (*DatabaseService, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		utils.DatabaseConifg.Host, utils.DatabaseConifg.Port, utils.DatabaseConifg.User,
		utils.DatabaseConifg.Password, utils.DatabaseConifg.DbName)
	conn, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return dbService, err
	}
	conn.SetMaxOpenConns(utils.DatabaseConifg.MaxOpenConns)
	conn.SetMaxIdleConns(utils.DatabaseConifg.MaxIdleConn)
	conn.SetConnMaxLifetime(0) //always REUSE
	dbService.Connection = conn
	return dbService, nil
}

func NamedExec(query string, args interface{}) error {
	_, err := dbService.Connection.NamedExec(query, args)
	return err
}

func ExecTx(query string, args ...interface{}) error {
	tx, err := dbService.Connection.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	stmt, err := tx.Prepare(dbService.Connection.Rebind(query))
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(args...)
	if err != nil {
		return err
	}

	return nil
}

func Get(query string, dest interface{}, args ...interface{}) error {
	err := dbService.Connection.Get(dest, query, args...)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

func Select(query string, dest interface{}, args ...interface{}) error {
	return dbService.Connection.Select(dest, query, args...)
}

func Close() {
	dbService.Connection.Close()
}
