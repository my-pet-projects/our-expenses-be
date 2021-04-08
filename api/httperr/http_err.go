package httperr

import (
	"encoding/json"
	"net/http"
	"our-expenses-server/entity"
)

type errorResponse struct {
	Message interface{} `json:"message"`
}

func RespondWithBadRequest(msg string, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(errorResponse{Message: msg})
	return
}

func RespondWithError(err error, w http.ResponseWriter, r *http.Request) {
	appError, ok := err.(entity.AppError)
	if ok {
		switch appError.Type {
		case entity.ErrorTypeIncorrectInput:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{Message: appError.Msg})
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	return
}
