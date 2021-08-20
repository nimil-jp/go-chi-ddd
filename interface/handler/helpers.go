package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

func newCtx() context.Context {
	return context.Background()
}

func bind(w http.ResponseWriter, r *http.Request, request interface{}) (ok bool) {
	length, err := strconv.Atoi(r.Header.Get("Content-Length"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}

	// Read body data to parse json
	body := make([]byte, length)
	length, err = r.Body.Read(body)
	if err != nil && err != io.EOF {
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}

	if err := json.Unmarshal(body[:length], &request); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}

	return true
}

func WriteJson(w http.ResponseWriter, code int, response interface{}) error {
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
