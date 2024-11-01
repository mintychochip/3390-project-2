package middleware

import (
	"api-3390/container/predicate"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strconv"
)

func URLParam(key string, predicates ...predicate.Predicate[string]) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			val := chi.URLParam(r, key)

			for _, p := range predicates {
				if !p.Test(val) {
					http.Error(w, p.ErrorMessage(val), http.StatusNotFound)
					return
				}
			}
			ctx := context.WithValue(r.Context(), key, val)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func InterceptJson(m map[string][]predicate.Predicate[string]) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Content-Type") != "application/json" {
				http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
				return
			}

			var payload map[string]interface{}
			decoder := json.NewDecoder(r.Body)

			if err := decoder.Decode(&payload); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			for key, predicates := range m {
				if val, exists := payload[key]; exists {
					str, ok := val.(string)
					if !ok {
						http.Error(w, fmt.Sprintf("Invalid type for key '%s'", key), http.StatusBadRequest)
						return
					}

					for _, p := range predicates {
						if !p.Test(str) {
							http.Error(w, p.ErrorMessage(str), http.StatusBadRequest)
							return
						}
					}
				}
			}

			body, err := json.Marshal(payload)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(body))
			r.Header.Set("Content-Length", strconv.Itoa(len(body)))

			h.ServeHTTP(w, r)
		})
	}
}
