package adapter

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iqbalbaharum/sol-stalker/internal/config"
	db "github.com/iqbalbaharum/sol-stalker/internal/database"
)

var (
	Database  *db.Database
	mySQLOnce sync.Once
)

func InitMySQLClient(dsn string) error {
	if dsn == "" {
		return errors.New("MySQL DSN is empty")
	}

	var initError error

	var client *sql.DB
	var err error
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

	})

	database, err := db.NewDatabase(client, config.MySqlDbName)
	if err != nil {
		return err
	}

	err = database.CreateDatabaseAndTable()
	if err != nil {
		return err
	}

	Database = database

	return initError
}

func GetMySQLClient() (*sql.DB, error) {
	if Database == nil {
		return nil, errors.New("MySQL client is not initialized. call InitMySQLClient first")
	}

	return Database.MysqlClient, nil
}
