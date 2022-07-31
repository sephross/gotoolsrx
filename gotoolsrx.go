package gotoolsrx

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type Tools struct {
}

type JSONResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Read, Decode and Valiidate JSON helper function
func (t *Tools) ReadJSON(w http.ResponseWriter, r *http.Request, d any) error {

	maxBytes := 1048576 // one megabyte
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(d)
	if err != nil {
		return err
	}
	err = dec.Decode(&struct{}{}) // check struct only contains single JSON object, not two or more.
	if err != io.EOF {
		return errors.New("body must only contain single JSON value")
	}
	return nil
}

// Validate and Write JSON helper function.
func (t *Tools) WriteJSON(w http.ResponseWriter, status int, d any, h ...http.Header) error { // variadic parameter allows it to be optional, ie the function may not receive it
	out, err := json.MarshalIndent(d, "", "\t")
	if err != nil {
		return err
	}
	if len(h) > 0 {
		for k, v := range h[0] {
			w.Header()[k] = v
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

// JSON error reponse writer helper function.
// Default value of 400, http.StatusBadRequest if no status is provided.
func (t *Tools) ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}
	var payload JSONResponse
	payload.Error = true
	payload.Message = err.Error()

	return t.WriteJSON(w, statusCode, payload)
}
