package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/iqbalbaharum/sol-stalker/internal/types"
)

type listenerStorage struct {
	client *sql.DB
}

func NewListenerStorage(client *sql.DB) *listenerStorage {
	return &listenerStorage{client: client}
}

func (s *listenerStorage) GetAll() ([]*types.Listener, error) {
	ctx := context.Background()

	stmt, err := s.client.PrepareContext(ctx, `SELECT * FROM `+TABLE_NAME_LISTENER)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrPrepareStatement, err)
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrExecuteQuery, err)
	}

	var listeners []*types.Listener

	for rows.Next() {
		var l types.Listener

		err = rows.Scan(&l.Signer)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", ErrScanData, err)
		}

		listeners = append(listeners, &l)
	}

	if len(listeners) == 0 {
		return nil, ErrListenerNotFound
	}

	return listeners, nil
}

func (s *listenerStorage) Add(listener *types.Listener) (*types.Listener, error) {
	ctx := context.Background()

	query := fmt.Sprintf(`INSERT INTO %s (signer) VALUES (?)`, TABLE_NAME_LISTENER)

	stmt, err := s.client.PrepareContext(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrPrepareStatement, err)
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, listener.Signer)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrExecuteQuery, err)
	}

	return listener, nil
}

func (s *listenerStorage) Remove(signer string) error {
	ctx := context.Background()

	stmt, err := s.client.PrepareContext(ctx, `
	DELETE FROM listener WHERE signer = ? 
	`)

	if err != nil {
		return fmt.Errorf("%s: %w", ErrPrepareStatement, err)
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, signer)

	if err != nil {
		return fmt.Errorf("%s: %w", ErrExecuteStatement, err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf("%s: %w", ErrRetrieveRows, err)
	}

	if rowsAffected == 0 {
		return ErrListenerNotFound
	}

	return nil
}
