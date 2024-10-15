package user

import (
	"database/sql"
	"fmt"
)

type Service[T any, K comparable] interface {
	itemExists(obj *T, q func(obj *T) (string, []interface{})) (bool, error)
	insertItem(obj *T, q func(obj *T) (string, []interface{})) error
	getItem(k []interface{}, query string, scan func(t *T, rows *sql.Rows) error) (*T, error)
	getAllItems(k []interface{}, query string, scan func(t *T, rows *sql.Rows) error) ([]*T, error)
}
type GenericService[T any, K comparable] struct {
	db *sql.DB
}

func (s *GenericService[T, K]) itemExists(obj *T, q func(obj *T) (string, []interface{})) (bool, error) {
	query, args := q(obj)
	var exists bool
	err := s.db.QueryRow(query, args...).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error whilst checking existence: %w", err)
	}
	fmt.Println(exists)
	return exists, nil
}
func (s *GenericService[T, K]) insertItem(obj *T, q func(obj *T) (string, []interface{})) error {
	query, args := q(obj)
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(args...)
	return err
}
func (s *GenericService[T, K]) getItem(k []interface{}, query string, scan func(t *T, rows *sql.Rows) error) (*T, error) {
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(k...)
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
func (s *GenericService[T, K]) getAllItems(k []interface{}, query string, scan func(t *T, rows *sql.Rows) error) ([]*T, error) {
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(k...)
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
