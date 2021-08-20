package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go-chi-ddd/interface/handler"
	"go-chi-ddd/pkg/xerrors"
)

type handlerFunc func(w http.ResponseWriter, r *http.Request) error

func get(r chi.Router, relativePath string, handlerFunc handlerFunc) {
	r.Get(relativePath, hf(handlerFunc))
}

func post(r chi.Router, relativePath string, handlerFunc handlerFunc) {
	r.Post(relativePath, hf(handlerFunc))
}

func put(r chi.Router, relativePath string, handlerFunc handlerFunc) {
	r.Put(relativePath, hf(handlerFunc))
}

func patch(r chi.Router, relativePath string, handlerFunc handlerFunc) {
	r.Patch(relativePath, hf(handlerFunc))
}

func delete(r chi.Router, relativePath string, handlerFunc handlerFunc) {
	r.Delete(relativePath, hf(handlerFunc))
}

func options(r chi.Router, relativePath string, handlerFunc handlerFunc) {
	r.Options(relativePath, hf(handlerFunc))
}

func head(r chi.Router, relativePath string, handlerFunc handlerFunc) {
	r.Head(relativePath, hf(handlerFunc))
}

func hf(handlerFunc handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if w.Header().Get("Content-Type") == "" {
			// Per spec, UTF-8 is the default, and the charset parameter should not
			// be necessary. But some clients (eg: Chrome) think otherwise.
			// Since json.Marshal produces UTF-8, setting the charset parameter is a
			// safe option.
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
		}
		err := handlerFunc(w, r)

		if err != nil {
			switch v := err.(type) {
			case *xerrors.Expected:
				if v.StatusOk() {
					return
				} else {
					err := handler.WriteJson(w, v.StatusCode(), v.Message())
					if err != nil {
						return
					}
				}
			case *xerrors.Validation:
				err := handler.WriteJson(w, http.StatusBadRequest, v)
				if err != nil {
					return
				}
			default:
				err := handler.WriteJson(w, http.StatusInternalServerError, v)
				if err != nil {
					return
				}
			}

			// _ = err.Error(errors.Errorf("%+v", err))
		}
	}
}
