package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorModel struct {
	Error string `json:"error"`
}

func LogErrors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				errV, ok := err.(error)
				if !ok {
					return
				}
				fmt.Println(errV)
				errObj := ErrorModel{errV.Error()}
				json.NewEncoder(w).Encode(errObj)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
