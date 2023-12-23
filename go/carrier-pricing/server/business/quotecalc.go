package business

import (
	"errors"
	"math"
	"sort"
	"strconv"

	"github.com/lalathealter/kensu/server/db"
)

type QuoteInputModel struct {
	PickupPC   string         `json:"pickup_postcode"`
	DeliveryPC string         `json:"delivery_postcode"`
	Vehicle    db.VehicleType `json:"vehicle"`
}

type CarrierDataOutput struct {
	Service  string `json:"service"`
	Price    int    `json:"price"`
	Delivery int    `json:"delivery_time"`
}

type PriceList []CarrierDataOutput

type QuoteOutputModel struct {
	PickupPC   string         `json:"pickup_postcode"`
	DeliveryPC string         `json:"delivery_postcode"`
	Vehicle    db.VehicleType `json:"vehicle"`
	PriceList PriceList `json:"price_list"`
}

func (q QuoteInputModel) ProduceOutput() (QuoteOutputModel, error) {
	pl, err := q.compilePriceList()
	out := QuoteOutputModel{
    q.PickupPC, q.DeliveryPC, q.Vehicle, pl,
  }
	return out, err
}

func (q QuoteInputModel) compilePriceList() (PriceList, error) {
	servs, err := q.Vehicle.FindServices()
	if err != nil {
		return nil, err
	}

	prices := make(PriceList, len(servs))
	for i, v := range servs {
		prices[i], err = q.produceCarrierDataOutput(v)
		if err != nil {
			return nil, err
		}
	}

	sort.Slice(prices, func(i, j int) bool {
		return prices[i].Price < prices[j].Price
	})

	return prices, nil
}

func (q QuoteInputModel) produceCarrierDataOutput(cs db.CarrierServiceRow) (CarrierDataOutput, error) {
	cdo := CarrierDataOutput{
		cs.Service, 0, cs.Delivery,
	}
	servPrice, err := q.deduceChargePrice(cs.Price + cs.Markup)
	cdo.Price = servPrice
	return cdo, err
}

const PC_DIV = 100_000_000

func (q QuoteInputModel) deduceChargePrice(basePrice int) (int, error) {
	dpc, ppc, err := q.parsePostcodesToIntegers()
  if err != nil {
    return 0, err
  }
	diff := (dpc - ppc) / PC_DIV
	if diff < 0 {
		diff = -diff
	}

	diff += basePrice

	f := q.getVehiclePriceFactor()
	diff = applyPriceFactor(diff, f)
	return diff, nil
}

func applyPriceFactor(p int, f float64) int {
	return int(math.Round(float64(p) * f))
}

const POSTCODE_BASE = 36


func (qinp QuoteInputModel) getVehiclePriceFactor() (factor float64) {
	factor = 1.00
	switch qinp.Vehicle {
	case db.Basic:
	case db.Bicycle:
		factor += .10
	case db.Motorbike:
		factor += .15
	case db.ParcelCar:
		factor += .20
	case db.SmallVan:
		factor += .30
	case db.LargeVan:
		factor += .40
	}

	return factor
}

var ErrWrongfulIntConversion = errors.New("Couldn't rightfully convert to general Integer type")

func (qinp QuoteInputModel) parsePostcodesToIntegers() (int, int, error) {
	ppc, err := strconv.ParseInt(qinp.PickupPC, POSTCODE_BASE, 64)
	dpc, err := strconv.ParseInt(qinp.DeliveryPC, POSTCODE_BASE, 64)
	tppc, tdpc := int(ppc), int(dpc)
	if err != nil || int64(tppc) != ppc || int64(tdpc) != dpc {
		err = ErrWrongfulIntConversion
	}
	return tppc, tdpc, err
}
