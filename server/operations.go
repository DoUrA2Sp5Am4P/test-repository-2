package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func getOperations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	operations_by_current_user := []Operation{}
	for i := 0; i < len(operations); i++ {
		if operations[i].Username == currentUser {
			operations_by_current_user = append(operations_by_current_user, Operation{Name: operations[i].Name, ExecutionTime: operations[i].ExecutionTime, Username: currentUser})
		}
	}
	json.NewEncoder(w).Encode(operations_by_current_user)
	InfoLogger.Println("get all operations OK")
}
func handleOperations(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		InfoLogger.Println("http get all operations")
		getOperations(w, r)
	default:
		ErrorLogger.Println("http get all operations error: method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not allowed")
	}
}

func handleAddition(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		InfoLogger.Println("http post operation")
		add, err1 := strconv.Atoi(r.URL.Query().Get("time"))
		if err1 != nil {
			ErrorLogger.Println("invalid operation")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid operations")
			return
		}
		InfoLogger.Println("Operation: addiction")
		InfoLogger.Println("Changed OK")
		for i := 0; i < len(operations); i++ {
			if operations[i].Name == "Сложение" && operations[i].Username == currentUser {
				operations[i].ExecutionTime = add
			}
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Changed OK")
	default:
		ErrorLogger.Println("http not allowed method (operation addiction)")
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not allowed")
	}
}

func handleSubtraction(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		InfoLogger.Println("http post operation")
		sub, err1 := strconv.Atoi(r.URL.Query().Get("time"))
		if err1 != nil {
			ErrorLogger.Println("invalid operation")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid operations")
			return
		}
		InfoLogger.Println("Operation: subtraction")
		for i := 0; i < len(operations); i++ {
			if operations[i].Name == "Вычитание" && operations[i].Username == currentUser {
				operations[i].ExecutionTime = sub
			}
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Changed OK")
	default:
		ErrorLogger.Println("http not allowed method (operation subtraction)")
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not allowed")
	}
}

func handleMultiplication(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		InfoLogger.Println("http post operation")
		mult, err1 := strconv.Atoi(r.URL.Query().Get("time"))
		if err1 != nil {
			ErrorLogger.Println("invalid operation")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid operations")
			return
		}
		InfoLogger.Println("Operation: multiplication")
		for i := 0; i < len(operations); i++ {
			if operations[i].Name == "Умножение" && operations[i].Username == currentUser {
				operations[i].ExecutionTime = mult
			}
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Changed OK")
	default:
		ErrorLogger.Println("http not allowed method (operation multiplication)")
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not allowed")
	}
}

func handleDivision(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		InfoLogger.Println("http post operation")
		div, err1 := strconv.Atoi(r.URL.Query().Get("time"))
		if err1 != nil {
			ErrorLogger.Println("invalid operation")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid operations")
			return
		}
		InfoLogger.Println("Operation: division")
		for i := 0; i < len(operations); i++ {
			if operations[i].Name == "Деление" && operations[i].Username == currentUser {
				operations[i].ExecutionTime = div
			}
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Changed OK")
	default:
		ErrorLogger.Println("http not allowed method (operation division)")
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not allowed")
	}
}
