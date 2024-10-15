package service

import (
	"database/sql"
	"fmt"
)

type service[T any, K comparable] interface {
	updateItem(query string, args []interface{}) error
	deleteItems(query string, args []interface{}) error
	itemExists(obj *T, q func(obj *T) (string, []interface{})) (bool, error)
	insertItem(obj *T, q func(obj *T) (string, []interface{})) error
	getItem(query string, args []interface{}, scan func(t *T, rows *sql.Rows) error) (*T, error)
	getAllItems(query string, args []interface{}, scan func(t *T, rows *sql.Rows) error) ([]*T, error)
}
type genericService[T any, K comparable] struct {
	db *sql.DB
}

func (s *genericService[T, K]) updateItem(query string, args []interface{}) error {
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(args...)
	if err != nil {
		return fmt.Errorf("failed to execute update: %w", err)
	}

	return nil
}
func (s *genericService[T, K]) deleteItems(query string, args []interface{}) error {
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(args...)
	if err != nil {
		return fmt.Errorf("failed to execute delete: %w", err)
	}

	return nil
}
func (s *genericService[T, K]) itemExists(obj *T, q func(obj *T) (string, []interface{})) (bool, error) {
	query, args := q(obj)
	var exists bool
	err := s.db.QueryRow(query, args...).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error whilst checking existence: %w", err)
	}
	fmt.Println(exists)
	return exists, nil
}
func (s *genericService[T, K]) insertItem(obj *T, q func(obj *T) (string, []interface{})) error {
	query, args := q(obj)
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(args...)
	return err
}
func (s *genericService[T, K]) getItem(query string, args []interface{}, scan func(t *T, rows *sql.Rows) error) (*T, error) {
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var item T
		if err := scan(&item, rows); err != nil {
			return nil, err
		}
		return &item, nil
	}
	return nil, nil
}
func (s *genericService[T, K]) getAllItems(query string, args []interface{}, scan func(t *T, rows *sql.Rows) error) ([]*T, error) {
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*T
	for rows.Next() {
		var item T
		if err := scan(&item, rows); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
