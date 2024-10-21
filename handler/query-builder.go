package handler

import "net/http"

type QueryParams struct {
	Operation string
	Column    string
}

type QueryHandler func(w http.ResponseWriter, column string, filePath string)

type QueryBuilder struct {
	queries     map[string]QueryHandler
	defaultCase QueryHandler
}

// Create a new QueryBuilder
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		queries: make(map[string]QueryHandler),
	}
}

// Add a query case with its handler
func (qb *QueryBuilder) AddQuery(operation string, handler QueryHandler) *QueryBuilder {
	qb.queries[operation] = handler
	return qb
}

// Build the query based on the operation
func (qb *QueryBuilder) Build(w http.ResponseWriter, params QueryParams, filePath string) {
	if handler, exists := qb.queries[params.Operation]; exists {
		handler(w, params.Column, filePath)
	} else if qb.defaultCase != nil {
		qb.defaultCase(w, "", filePath)
	} else {
		http.Error(w, "invalid operation", http.StatusBadRequest)
	}
}

// Set the default case handler
func (qb *QueryBuilder) SetDefaultCase(handler QueryHandler) *QueryBuilder {
	qb.defaultCase = handler
	return qb
}
