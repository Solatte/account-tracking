package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	mysqlClient *sql.DB
	mySQLOnce   sync.Once
	DbName      string = "sol_stalker"
)

// Error description
const (
	ErrPrepareStatement = "failed to prepare SQL statement"
	ErrExecuteStatement = "failed to execute statement"
	ErrExecuteQuery     = "failed to execute query"
	ErrScanData         = "failed to scan data"
	ErrBeginTransaction = "failed to begin transaction"
	ErrRollback         = "failed to rollback transaction"
	ErrCommit           = "failed to commit transaction"
	ErrRetrieveRows     = "failed to retrieve rows affected"
)

var (
	ErrListenerNotFound = errors.New("listener not found")
	ErrTradeNotFound    = errors.New("trade not found")
)

type MySQLFilter struct {
	Query []MySQLQuery
}

type MySQLQuery struct {
	Column string
	Op     string
	Query  string
}

type Column struct {
	Field string
	Type  string
}

func CreateDatabaseAndTable() error {

	createDatabase := `CREATE DATABASE IF NOT EXISTS ` + DbName

	_, err := mysqlClient.Exec(createDatabase)

	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Failed to create db %s: %v", DbName, err))
	}

	useDatabase := `USE ` + DbName

	_, err = mysqlClient.Exec(useDatabase)

	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Failed to use db %s: %v", DbName, err))
	}

	path := "./external/api/database/migrations/"

	entries, err := os.ReadDir(path)

	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		c, err := os.ReadFile(fmt.Sprintf(path + e.Name()))

		if err != nil {
			log.Fatal(err)
		}

		_, err = mysqlClient.Exec(string(c))

		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func InitMySQLClient(dsn string) error {
	if dsn == "" {
		return errors.New("MySQL DSN is empty")
	}

	var initError error
	mySQLOnce.Do(func() {
		client, err := sql.Open("mysql", dsn)
		if err != nil {
			initError = fmt.Errorf("failed to connect to MySQL: %v", err)
			return
		}

		if err := client.Ping(); err != nil {
			initError = fmt.Errorf("failed to ping MySQL: %v", err)
			return
		}

		mysqlClient = client

	})

	CreateDatabaseAndTable()

	return initError
}

func GetMySQLClient() (*sql.DB, error) {
	if mysqlClient == nil {
		return nil, errors.New("MySQL client is not initialized. call InitMySQLClient first")
	}
	return mysqlClient, nil
}
