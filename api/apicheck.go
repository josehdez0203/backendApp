package main

import (
	"net/http"

	"github.com/josehdez0203/realstate/logger"
)

func (app *application) Hello(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "Go movies up and running",
		Version: "1.0.0",
	}
	logger.L_Info("Checando API 🌟")
	_ = app.writeJSON(w, http.StatusOK, payload)
}
