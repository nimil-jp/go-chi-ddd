package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"

	"go-chi-ddd/pkg/xerrors"
)

type handlerFunc func(w http.ResponseWriter, r *http.Request) error

func Get(r chi.Router, relativePath string, handlerFunc handlerFunc) {
	r.Get(relativePath, hf(handlerFunc))
}

func Post(r chi.Router, relativePath string, handlerFunc handlerFunc) {
	r.Post(relativePath, hf(handlerFunc))
}

func Put(r chi.Router, relativePath string, handlerFunc handlerFunc) {
	r.Put(relativePath, hf(handlerFunc))
}

func Patch(r chi.Router, relativePath string, handlerFunc handlerFunc) {
	r.Patch(relativePath, hf(handlerFunc))
}

func Delete(r chi.Router, relativePath string, handlerFunc handlerFunc) {
	r.Delete(relativePath, hf(handlerFunc))
}

func Options(r chi.Router, relativePath string, handlerFunc handlerFunc) {
	r.Options(relativePath, hf(handlerFunc))
}

func Head(r chi.Router, relativePath string, handlerFunc handlerFunc) {
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
					err := writeJson(w, v.StatusCode(), v.Message())
					if err != nil {
						return
					}
				}
			case *xerrors.Validation:
				err := writeJson(w, http.StatusBadRequest, v)
				if err != nil {
					return
				}
			default:
				err := writeJson(w, http.StatusInternalServerError, v)
				if err != nil {
					return
				}
			}

			//_ = err.Error(errors.Errorf("%+v", err))
		}
	}
}

func writeJson(w http.ResponseWriter, code int, response interface{}) error {
	w.WriteHeader(code)
	b, err := json.Marshal(response)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	if err != nil {
		return err
	}
	return nil
}
