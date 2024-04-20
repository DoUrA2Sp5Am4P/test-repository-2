package main

import (
	"net/http"
	"time"
)

var cook = 0

func delCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cook == 0 {
			c := &http.Cookie{
				Name:    "token",
				Value:   "",
				Path:    "/",
				Expires: time.Unix(0, 0),
			}
			http.SetCookie(w, c)
			cook = 1
		}
		next.ServeHTTP(w, r)
	})
}
