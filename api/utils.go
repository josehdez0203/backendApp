package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
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

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1024 * 1024 // 1 MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}

func (app *application) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}
	var payload JSONResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJSON(w, statusCode, payload)
}

// func (app *application) LogInfo(m string) {
// 	infoColor := color.New(color.FgHiBlue).SprintFunc()
// 	fechaColor := color.New(color.FgYellow).SprintFunc()
// 	fecha := time.Now()
// 	d := fecha.Format(time.DateOnly)
// 	t := fecha.Format(time.TimeOnly)
// 	fmt.Printf("%s %s INFO: %s\n", fechaColor(d), fechaColor(t), infoColor(m))
// }
//
// func (app *application) LogError(m string) {
// 	errorColor := color.New(color.FgRed).SprintFunc()
// 	fechaColor := color.New(color.FgHiRed).SprintFunc()
// 	fecha := time.Now()
// 	d := fecha.Format(time.DateOnly)
// 	t := fecha.Format(time.TimeOnly)
// 	res := errorColor(" ERROR ")
// 	fmt.Printf("%s %s %s %s\n", fechaColor(d), fechaColor(t), res, errorColor(m))
// }
