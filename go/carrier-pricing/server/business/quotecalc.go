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
	PriceList  PriceList      `json:"price_list"`
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

	dist, err := q.findDistanceBetween()
	if err != nil {
		return nil, err
	}

	prices := make(PriceList, len(servs))
	for i, v := range servs {
		prices[i] = q.produceCarrierDataOutput(dist, v)
	}

	sort.Slice(prices, func(i, j int) bool {
		return prices[i].Price < prices[j].Price
	})

	return prices, nil
}

var ErrWrongfulIntConversion = errors.New("Couldn't rightfully convert to general Integer type")

const POSTCODE_BASE = 36

func parsePostcode(pc string) (int, error) {
	ppc, err := strconv.ParseInt(pc, POSTCODE_BASE, 64)
	tppc := int(ppc)
	if err != nil || int64(tppc) != ppc {
		err = ErrWrongfulIntConversion
	}
	if tppc < 0 {
		tppc = -tppc
	}
	return tppc, err
}

func (q QuoteInputModel) findDistanceBetween() (int, error) {
	dpc, err := parsePostcode(q.DeliveryPC)
	if err != nil {
		return 0, err
	}
	ppc, err := parsePostcode(q.PickupPC)
	dist := findCostOfDifference(dpc, ppc)
	return dist, err
}

const PC_DIV = 100_000_000

func findCostOfDifference(a, b int) int {
	diff := (a - b) / PC_DIV
	if diff < 0 {
		diff = -diff
	}
	return diff
}

func (q QuoteInputModel) produceCarrierDataOutput(dist int, cs db.CarrierServiceRow) CarrierDataOutput {
	cdo := CarrierDataOutput{
		cs.Service, 0, cs.Delivery,
	}
	servPrice := q.deduceChargePrice(dist + cs.Price + cs.Markup)
	cdo.Price = servPrice
	return cdo
}

func (q QuoteInputModel) deduceChargePrice(basePrice int) int {
	price := basePrice
	f := q.getVehiclePriceFactor()
	price = applyPriceFactor(price, f)
	return price
}

func applyPriceFactor(p int, f float64) int {
	return int(math.Round(float64(p) * f))
}

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
