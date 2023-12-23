package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/lalathealter/kensu/server/business"
)

func HandleQuotes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handlePostQuotes(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handlePostQuotes(w http.ResponseWriter, r *http.Request) {

	qinput := business.QuoteInputModel{}

	if err := json.NewDecoder(r.Body).Decode(&qinput); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	out, err := qinput.ProduceOutput()
	if err != nil {
		if err == business.ErrWrongfulIntConversion {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(out); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
}
