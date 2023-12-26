package business

import (
	"fmt"
	"testing"

	"github.com/lalathealter/kensu/server/db"
)

func TestParsePostcode(t *testing.T) {
	validCases := map[string]int{
		"DDD":           17329,
		"0":             0,
		"11":            37,
		"1289":          49545,
		"FF":            555,
		"fFfF":          719835,
		"MKQW38":        1365102836,
		"-332":          3998,
		"FFFFFFFFFF":    1566925045741275,
		"1Y2P0IJ32E8E7": 9223372036854775807,
	}

	for input, want := range validCases {
		got, err := parsePostcode(input)
		if err != nil || got != want {
			fmt.Println("wanted", want, "but got", got)
			t.Fail()
		}
	}

	errCases := []string{
		"//qwe", "", ",.<AD", "AD.,eqw", "-",
		"---323", "CCBQWUDBQUWCZWQ", "1Y2P0IJ32E8E8",
	}

	for _, errC := range errCases {
		got, err := parsePostcode(errC)
		if err != ErrWrongfulIntConversion {
			fmt.Println("for", errC, "wanted error, but got", got)
			t.Fail()
		}
	}
}

func TestDeduceChargePrice(t *testing.T) {
	cases := map[QuoteInputModel]map[int]int{
		{"ADD", "DDD", db.Bicycle}:   {10: 11, 111: 122, 100: 110},
		{"ADD", "DDD", db.SmallVan}:  {10: 13, 33: 43},
		{"ADD", "DDD", db.ParcelCar}: {20: 24, 100: 120, 0: 0},
		{"ADD", "DDD", db.LargeVan}:  {20: 28, 0: 0, 100: 140, 28: 39},
		{"ADD", "DDD", db.Motorbike}: {20: 23, 100: 115, 900: 1035},
		{"ADD", "DDD", db.Basic}:     {1: 1, 0: 0, 999: 999, 22: 22},
	}

	for qinp, paires := range cases {
		for input, want := range paires {
			got := qinp.deduceChargePrice(input)
			if got != want {
				fmt.Println(qinp.Vehicle)
				fmt.Println(want, "is wanted, but got", got)
				t.Fail()
			}
		}
	}
}
