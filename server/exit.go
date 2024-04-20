package main

import (
	"net/http"
	"os"
	"time"
)

var Exit = 0

func handleExit(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/exit.html", http.StatusSeeOther)
	Exit = 1
	go func() {
		time.Sleep(10 * time.Second)
		InfoLogger.Println("Program exits with code 0")
		os.Exit(0)
	}()
}
