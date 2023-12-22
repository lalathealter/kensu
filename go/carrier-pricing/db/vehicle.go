package db

import "errors"

type VehicleType string

func (v VehicleType) IsSupported() bool {
	switch v {
	case Basic:
	case Bicycle:
	case Motorbike:
	case ParcelCar:
	case SmallVan:
	case LargeVan:
	default:
		return false
	}
	return true
}

const (
	Basic     VehicleType = "basic"
	Bicycle   VehicleType = "bicycle"
	Motorbike VehicleType = "motorbike"
	ParcelCar VehicleType = "parcel_car"
	SmallVan  VehicleType = "small_van"
	LargeVan  VehicleType = "large_van"
)

var getServicesSQLTemp = `
  SELECT name, price, time, markup DISTINCT
  FROM carriers
  WHERE $1=ANY(vehicles)
;`

var ErrVehicleTypeNotSupported = errors.New("Provided vehicle type isn't supported by this API;")

func (vh VehicleType) FindServices() ([]CarrierServiceRow, error) {
	if !vh.IsSupported() {
		return nil, ErrVehicleTypeNotSupported
	}
	db := Get()
	rows, err := db.Query(getServicesSQLTemp, vh)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]CarrierServiceRow, 0)
	for rows.Next() {
		curr := CarrierServiceRow{}
		err := rows.Scan(
			&curr.Service, &curr.Price, &curr.Delivery, &curr.Markup,
		)
		if err != nil {
			return nil, err
		}
		out = append(out, curr)
	}

	return out, nil
}

type CarrierServiceRow struct {
	Service  string
	Price    int
	Markup   int
	Delivery int
}
