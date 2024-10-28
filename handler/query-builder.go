package handler

import "net/http"

type QueryParams struct {
	Operation string   `json:"operation"`
	Column    []string `json:"columns"`
}

type QueryHandler func(w http.ResponseWriter, column []string, filePath string)

type QueryBuilder struct {
	queries     map[string]QueryHandler
	defaultCase QueryHandler
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		queries: make(map[string]QueryHandler),
	}
}

func (qb *QueryBuilder) AddQuery(operation string, handler QueryHandler) *QueryBuilder {
	qb.queries[operation] = handler
	return qb
}

func (qb *QueryBuilder) Build(w http.ResponseWriter, params QueryParams, filePath string) {
	if handler, exists := qb.queries[params.Operation]; exists {
		handler(w, params.Column, filePath)
	} else if qb.defaultCase != nil {
		qb.defaultCase(w, make([]string, 0), filePath)
	} else {
		http.Error(w, "invalid operation", http.StatusBadRequest)
	}
}

func (qb *QueryBuilder) SetDefaultCase(handler QueryHandler) *QueryBuilder {
	qb.defaultCase = handler
	return qb
}
