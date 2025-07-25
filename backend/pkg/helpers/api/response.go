package api

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
)

type ApiResponse struct {
	Json      func(w http.ResponseWriter, logger *zerolog.Logger, t any, statusCode ...int)
	JsonError func(w http.ResponseWriter, logger *zerolog.Logger, message string, statusCode int)
}

func responseJson(w http.ResponseWriter, logger *zerolog.Logger, t any, statusCode ...int) {
	data, err := json.Marshal(t)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to marshal response data")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if len(statusCode) > 0 {
		w.WriteHeader(statusCode[0])
	} else {
		w.WriteHeader(http.StatusOK)
	}

	if _, err := w.Write(data); err != nil {
		logger.Error().Err(err).Msg("Failed to write response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func responseJsonError(w http.ResponseWriter, logger *zerolog.Logger, message string, statusCode int) {
	responseJson(w, logger, map[string]string{"error": message}, statusCode)
}
