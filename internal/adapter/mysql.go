package adapter

import (
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iqbalbaharum/sol-stalker/internal/utils"
)

var (
	mysqlClient *sql.DB
	mySQLOnce   sync.Once
	DbName      string = "sol_stalker"
	TableName   string = "trade3"
)

type column struct {
	Field string
	Type  string
}

var Column = []column{
	{Field: "amm_id", Type: "TEXT"},
	{Field: "mint", Type: "TEXT"},
	{Field: "action", Type: "TEXT"},
	{Field: "compute_limit", Type: "INT"},
	{Field: "compute_price", Type: "INT"},
	{Field: "amount", Type: "BIGINT"},
	{Field: "signature", Type: "TEXT"},
	{Field: "timestamp", Type: "INT"},
	{Field: "tip", Type: "TEXT"},
	{Field: "tip_amount", Type: "INT"},
	{Field: "status", Type: "TEXT"},
}

func CreateColumn() error {
	query := fmt.Sprintf("SELECT * FROM %s LIMIT 1", TableName)
	rows, err := mysqlClient.Query(query)

	if err != nil {

		return fmt.Errorf(fmt.Sprintf("Failed to get rows: %v", err))

	}

	defer rows.Close()

	// Get column names
	cols, err := rows.Columns()

	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Failed to get columns: %v", err))
	}

	createColumn := fmt.Sprintf(`ALTER TABLE %s `, TableName)

	for _, col := range Column {
		if slices.Contains(cols, col.Field) {
			continue
		}

		addColumn := fmt.Sprintf("ADD %s %s, ", col.Field, col.Type)
		createColumn = createColumn + addColumn
	}

	createColumn = utils.ReplaceLastComma(createColumn, ";")

	_, err = mysqlClient.Exec(createColumn)

	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Failed to create columns: %v", err))
	}

	return nil
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

	createTable := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s `, TableName)
	addColumn := "("

	for _, col := range Column {
		addColumn += fmt.Sprintf("%s %s, ", col.Field, col.Type)
	}

	addColumn = utils.ReplaceLastComma(addColumn, ");")
	createTable += addColumn

	_, err = mysqlClient.Exec(createTable)

	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Failed to create table: %v", err))
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

	return initError
}

func GetMySQLClient() (*sql.DB, error) {
	if mysqlClient == nil {
		return nil, errors.New("MySQL client is not initialized. call InitMySQLClient first")
	}
	return mysqlClient, nil
}
