package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

const contentTypeJson = "application/json"

func main() {
	fmt.Println("cli-client v1.0.0")

	address := GetTargetAddress()
	var PostQuoteRequest = BindPostQuoteRequest(address)

	for {

		qb, err := FormQuoteBody()
		if err != nil {
			fmt.Println("failed to form a request body:", err)
			continue
		}

		res, err := PostQuoteRequest(qb)
		if err != nil {
			fmt.Println("failed to send a request:", err.Error())
			continue
		}
		defer res.Body.Close()

		var ans any
		if err := json.NewDecoder(res.Body).Decode(&ans); err != nil {
			fmt.Println("failed to decode a response:", err)
			continue
		}
		pretty, _ := json.MarshalIndent(ans, "", "\t")
		fmt.Println(string(pretty))
	}
}

var GetInputWith = BindGetInputWith()

func BindGetInputWith() func(*string, string) {
	scanner := bufio.NewScanner(os.Stdin)
	return func(dest *string, msg string) {
		fmt.Println(msg)
		scanner.Scan()
		inp := scanner.Text()
		(*dest) = inp
	}
}

type QuoteBody map[string]string

var (
	PickPS  = `pickup_postcode`
	DelPS   = `delivery_postcode`
	Vehicle = `vehicle`
)

var inputTokens = [3]string{PickPS, DelPS, Vehicle}

func FormQuoteBody() (QuoteBody, error) {
	qb := make(map[string]string, len(inputTokens))

	for _, token := range inputTokens {
		var inpVal string
		GetInputWith(&inpVal, "ENTER "+token)
		qb[token] = inpVal
	}

	return qb, nil
}

const LOCAL_TEST_ADDR = "localhost:8080/quotes"

func GetTargetAddress() string {
	var addr string
	GetInputWith(&addr, "Enter an address to connect to:")
	if addr == "" {
		addr = LOCAL_TEST_ADDR
	}

	if hasNoHTTPScheme(addr) {
		addr = "http://" + addr
	}

	pingAddress(addr)

	return addr
}

func hasNoHTTPScheme(addr string) bool {
	return !strings.HasPrefix(addr, "https://") || !strings.HasPrefix(addr, "http://")
}

func pingAddress(addr string) {
	res, e := http.Get(addr)
	if e != nil {
		log.Fatalf("couldn't reach out to the specified address: %v", e)
	}
	defer res.Body.Close()
}

func BindPostQuoteRequest(addr string) func(QuoteBody) (*http.Response, error) {
	return func(qb QuoteBody) (*http.Response, error) {
		reqBodyBytes, err := json.Marshal(qb)
		if err != nil {
			return nil, err
		}
		bodyBuf := bytes.NewBuffer(reqBodyBytes)
		res, err := http.Post(addr, contentTypeJson, bodyBuf)
		return res, err
	}
}
