package main

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/lalathealter/kensu/server/business"
	"github.com/lalathealter/kensu/server/db"
)

func TestHasNoHTTPScheme(t *testing.T) {
	cases := map[string]bool{
		"http://":            false,
		"https://":           false,
		"https:123///dasd":   true,
		"dgqwe:https://":     true,
		"ttps://":            true,
		"http:////eqweda":    false,
		"https://google.com": false,
	}

	for got, want := range cases {
		if want != hasNoHTTPScheme(got) {
			t.Fail()
		}
	}
}

const errorResp = `{"error":"%v"}`

func TestPostQuoteRequest(t *testing.T) {
	addr := clampAddress("")
	var PostQuoteRequest = BindPostQuoteRequest(addr)

	errorCases := [][4]string{
		{"DA", "DA", "",
			fmt.Sprintf(errorResp, db.ErrVehicleTypeNotSupported.Error())},
		{"DA", "DA", "bycicle",
			fmt.Sprintf(errorResp, db.ErrVehicleTypeNotSupported.Error())},
		{"qDA", "/d8c", string(db.Bicycle),
			fmt.Sprintf(errorResp, business.ErrWrongfulIntConversion.Error())},
		{"ZZZZZZZZZZZZZZ", "0", string(db.Basic),
			fmt.Sprintf(errorResp, business.ErrWrongfulIntConversion.Error())},
		{"qA  DA", "dac", string(db.Bicycle),
			fmt.Sprintf(errorResp, business.ErrWrongfulIntConversion.Error())},
	}

	for i, inp := range errorCases {
		qb := QuoteBody{}
		qb[PickPS] = inp[0]
		qb[DelPS] = inp[1]
		qb[Vehicle] = inp[2]
		res, err := PostQuoteRequest(qb)
		if err != nil {
			t.Fatal(err)
		}

		got, _ := ioutil.ReadAll(res.Body)
		want := inp[3] + "\n"
		sgot := string(got)
		if sgot != want {
			fmt.Println("case", i+1)
			fmt.Print(sgot)
			fmt.Println(want)
			t.Fail()
		}
	}

}
