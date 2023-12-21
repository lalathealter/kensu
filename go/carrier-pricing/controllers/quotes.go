package controllers

import (
	"encoding/json"
	"errors"
	"math"
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
	PickupPC   string      `json:"pickup_postcode"`
	DeliveryPC string      `json:"delivery_postcode"`
	Vehicle    VehicleType `json:"vehicle"`
	Price      int         `json:"price"`
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
	f, err := q.GetVehiclePriceFactor()
	diff = applyPriceFactor(diff, f)
	return diff, err
}

func applyPriceFactor(p int, f float64) int {
	return int(math.Round(float64(p) * f))
}

const POSTCODE_BASE = 36

type VehicleType string

const (
	Basic     VehicleType = ""
	Bicycle   VehicleType = "bicycle"
	Motorbike VehicleType = "motorbike"
	ParcelCar VehicleType = "parcel_car"
	SmallVan  VehicleType = "small_van"
	LargeVan  VehicleType = "large_van"
)

var ErrVehicleTypeNotSupported = errors.New("Provided vehicle type isn't supported by this API;")

func (qinp QuoteInputModel) GetVehiclePriceFactor() (factor float64, err error) {
	factor = 1.00
	switch qinp.Vehicle {
	case Basic:
	case Bicycle:
		factor += .10
	case Motorbike:
		factor += .15
	case ParcelCar:
		factor += .20
	case SmallVan:
		factor += .30
	case LargeVan:
		factor += .40
	default:
		err = ErrVehicleTypeNotSupported
	}

	return factor, err
}

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
