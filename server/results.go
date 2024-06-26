package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	pb "proto"
	"time"
)

func handleResults(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		//InfoLogger.Println("http post results")
		//processResult(w, r)
	default:
		InfoLogger.Println("http get results")
		getResult(w, r)
	}
}
func getResult(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	results_ := []Expression_with_result{}
	for i, res := range results {
		if res.Username == currentUser {
			results_ = append(results_, results[i])
		}
	}
	json.NewEncoder(w).Encode(results_)
	InfoLogger.Println("http get results OK")
}
func (s *Server) SendResultRPC(
	ctx context.Context,
	in *pb.SendResult,
) (*pb.SendResultResponse, error) {
	x := 0
	for i, res := range results {
		if res.ID == in.ID {
			results[i] = Expression_with_result{ID: in.ID, Result: in.Result, Username: currentUser}
			x = 1
		}
	}
	if x == 0 {
		results = append(results, Expression_with_result{ID: in.ID, Result: in.Result, Username: currentUser})
	}
	for i := 0; i < len(expressions); i++ {
		if expressions[i].ID == in.ID {
			expressions[i].CompletedAt = time.Now()
			expressions[i].Status = "Completed"
		}
	}
	return nil, nil
}

func processResult(w http.ResponseWriter, r *http.Request) {
	var result Expression_with_result
	err := json.NewDecoder(r.Body).Decode(&result)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid result")
		ErrorLogger.Println("Invalid result")
		return
	}

	for i, expression := range expressions {
		if expression.ID == result.ID {
			expressions[i].Status = "Completed"
			expressions[i].CompletedAt = time.Now()
			InfoLogger.Printf("Expression ID: %s, Value: %s\n", result.ID, result.Result)
			break
		}
	}
	for i, expression := range results {
		if expression.ID == result.ID {
			results[i] = result
		}
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Result processed")
	InfoLogger.Println("Result processed")
}

func handleResultsByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/results/"):]
	switch r.Method {
	case http.MethodGet:
		InfoLogger.Println("get result by id")
		getResultsByID(w, r, id)
	default:
		InfoLogger.Println("get results by id error: method not aloowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not allowed")
	}
}

func getResultsByID(w http.ResponseWriter, r *http.Request, id string) {
	for _, expression := range results {
		if expression.ID == id && expression.Username == currentUser {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(expression)
			InfoLogger.Println("result by id get OK")
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	ErrorLogger.Println("result not found")
	fmt.Fprintf(w, "result not found")
}
