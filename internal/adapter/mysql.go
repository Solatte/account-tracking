package adapter

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	mysqlClient *sql.DB
	mySQLOnce   sync.Once
)

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

	return initError
}

func GetMySQLClient() (*sql.DB, error) {
	if mysqlClient == nil {
		return nil, errors.New("MySQL client is not initialized. call InitMySQLClient first")
	}
	return mysqlClient, nil
}
