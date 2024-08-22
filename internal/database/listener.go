package database

import (
	"context"
	"fmt"
)

type Listener struct {
	Signer string
}

func GetAllListener(ctx context.Context) ([]*Listener, error) {
	s, err := mysqlClient.PrepareContext(ctx, `SELECT * FROM listener`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrPrepareStatement, err)
	}

	defer s.Close()

	rows, err := s.QueryContext(ctx)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrExecuteQuery, err)
	}

	var listeners []*Listener

	for rows.Next() {
		var l Listener

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

func AddListener(ctx context.Context, listener *Listener) (*Listener, error) {
	s, err := mysqlClient.PrepareContext(ctx, `INSERT INTO listener (signer) VALUES (?)`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrPrepareStatement, err)
	}

	defer s.Close()

	_, err = s.ExecContext(ctx, listener.Signer)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrExecuteQuery, err)
	}

	return listener, nil
}

func RemoveListener(ctx context.Context, signer string) error {
	s, err := mysqlClient.PrepareContext(ctx, `
	DELETE FROM listener WHERE signer = ? 
	`)

	if err != nil {
		return fmt.Errorf("%s: %w", ErrPrepareStatement, err)
	}

	defer s.Close()

	result, err := s.ExecContext(ctx, signer)

	fmt.Println(signer)

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
