package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func HandleQuotes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handlePostQuotes(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

type QuoteInputModel struct {
	PickupPC   string `json:"pickup_postcode"`
	DeliveryPC string `json:"delivery_postcode"`
	Price      int    `json:"price"`
}

func handlePostQuotes(w http.ResponseWriter, r *http.Request) {

	qinput := QuoteInputModel{}

	if err := json.NewDecoder(r.Body).Decode(&qinput); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	price, err := qinput.DeduceChargePrice()
	if err != nil {
		if err == ErrWrongfulIntConversion {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		panic(err)
	}

	qinput.Price = price
	if err := json.NewEncoder(w).Encode(qinput); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
}

const PC_DIV = 100_000_000

func (q QuoteInputModel) DeduceChargePrice() (int, error) {
	dpc, ppc, err := q.ParsePostcodesToIntegers()
	diff := (dpc - ppc) / PC_DIV
	if diff < 0 {
		diff = -diff
	}
	return diff, err
}

const POSTCODE_BASE = 36

var ErrWrongfulIntConversion = errors.New("Couldn't rightfully convert to general Integer type")

func (qinp QuoteInputModel) ParsePostcodesToIntegers() (int, int, error) {
	ppc, err := strconv.ParseInt(qinp.PickupPC, POSTCODE_BASE, 64)
	dpc, err := strconv.ParseInt(qinp.DeliveryPC, POSTCODE_BASE, 64)
	tppc, tdpc := int(ppc), int(dpc)
	if err != nil || int64(tppc) != ppc || int64(tdpc) != dpc {
		err = ErrWrongfulIntConversion
	}
	return tppc, tdpc, err
}
