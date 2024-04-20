package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func handleDatabaseTime(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		InfoLogger.Println("http get database time")
		getDatabaseTime(w, r)
	case http.MethodPost:
		InfoLogger.Println("http post database time")
		changeDatabaseTime(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		ErrorLogger.Println("error: method not allowed (database time)")
		fmt.Fprintf(w, "Method not allowed")
	}
}

func getDatabaseTime(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < len(sl); i++ {
		if currentUser == sl[i].Username && len(sl) > 0 {
			json.NewEncoder(w).Encode(sl[i].time)
			InfoLogger.Println("get database time OK")
			return
		}
	}
	ErrorLogger.Println("get database time ERROR")
}

func changeDatabaseTime(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Query().Get("time")
	sl_, err := strconv.Atoi(t)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ErrorLogger.Println("error (post databse time)")
		fmt.Fprint(w, "Error")
		return
	}
	InfoLogger.Println("post database time OK")
	w.WriteHeader(http.StatusOK)
	for i := 0; i < len(sl); i++ {
		if sl[i].Username == currentUser {
			sl[i].time = sl_
			return
		}
	}
	sl = append(sl, sleep{Username: currentUser, time: sl_})
}
