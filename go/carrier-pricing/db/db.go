package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

func BindDB(db *sql.DB) func() *sql.DB {
	return func() *sql.DB {
		return db
	}
}

var Get func() *sql.DB

func init() {
	dbInstance := initDB()
	Get = BindDB(dbInstance)
}

const PRICE_LIST_JSON = "./carrier-data.json"

func initDB() *sql.DB {

	db, err := sql.Open("postgres", makeConnString())
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	err = uploadPriceList(db, PRICE_LIST_JSON)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

const (
	DBkey     = "dbname"
	DBHOSTkey = "dbhost"
	PORTkey   = "dbport"
	PASSkey   = "password"
	USERkey   = "user"
)

func makeConnString() string {
	if err := godotenv.Load(); err != nil {
		return ""
	}
	db := os.Getenv(DBkey)
	dbhost := os.Getenv(DBHOSTkey)
	port := os.Getenv(PORTkey)
	pass := os.Getenv(PASSkey)
	user := os.Getenv(USERkey)
	connS := fmt.Sprintf("host=%v port=%v dbname=%v user=%v password=%v sslmode=disable",
		dbhost, port, db, user, pass,
	)
	return connS
}

type CarrierServiceInput struct {
	DeliveryTime int           `json:"delivery_time"`
	Markup       int           `json:"markup"`
	Vehicles     []VehicleType `json:"vehicles"`
}

type CarrierDataInput struct {
	Carrier  string `json:"carrier_name"`
	Price    int    `json:"base_price"`
	Services []CarrierServiceInput
}

const uploadCarrierSQLTemp = `
  INSERT INTO carriers(name, price, time, markup, vehicles)
  VALUES($1, $2, $3, $4, $5)
  ON CONFLICT DO NOTHING;
  `

func uploadPriceList(db *sql.DB, p string) error {
	pList, err := parsePriceFile(p)
	if err != nil {
		return err
	}

	for _, data := range pList {

		for _, s := range data.Services {
			_, err = db.Exec(uploadCarrierSQLTemp,
				data.Carrier, data.Price, s.DeliveryTime, s.Markup, pq.Array(s.Vehicles),
			)
			if err != nil {
				return err
			}
		}
	}

	return err
}

func parsePriceFile(p string) ([]CarrierDataInput, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}

	var pList []CarrierDataInput
	err = json.NewDecoder(f).Decode(&pList)
	return pList, err
}
