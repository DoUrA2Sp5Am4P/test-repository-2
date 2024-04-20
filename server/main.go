package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	pb "proto"
)

type Server struct {
	pb.WorkerServiceServer
}

func NewServer() *Server {
	return &Server{}
}

type worker_info struct {
	Worker_number string `json:"worker_number"`
	Status        string `json:"status"`
}

var workers []worker_info

type Expression struct {
	Username    string    `json:"username"`
	ID          string    `json:"id"`
	Expression  string    `json:"expression"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at"`
}
type Expression_with_result struct {
	Username string `json:"username"`
	ID       string `json:"id"`
	Result   string `json:"value"`
}
type Operation struct {
	Username      string `json:"username"`
	Name          string `json:"name"`
	ExecutionTime int    `json:"execution_time"`
}
type sleep struct {
	Username string
	time     int
}

var sl []sleep
var expressions []Expression
var tasks []Expression
var results []Expression_with_result
var operations = []Operation{
	// {Name: "Сложение", ExecutionTime: 1},
	// {Name: "Вычитание", ExecutionTime: 2},
	// {Name: "Умножение", ExecutionTime: 3},
	// {Name: "Деление", ExecutionTime: 4},
}
var InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
var ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
var WarningLogger = log.New(os.Stderr, "WARNING: ", log.Ldate|log.Ltime)

func main() {
	InfoLogger.Println("start program")
	go startRPC()
	ReadAll()
	go saveAll_r()
	InfoLogger.Println("started save opertions")
	InfoLogger.Println("started backend")
	mux := http.NewServeMux()
	mux.Handle("/", AuthorizedUsers(http.FileServer(http.Dir("./html/"))))
	mux.Handle("/expressions", AuthorizedUsers(http.HandlerFunc(handleExpressions)))
	mux.Handle("/expressions/", AuthorizedUsers(http.HandlerFunc(handleExpressionByID)))
	mux.Handle("/operations", AuthorizedUsers(http.HandlerFunc(handleOperations)))
	mux.Handle("/operations/addition", AuthorizedUsers(http.HandlerFunc(handleAddition)))
	mux.Handle("/operations/subtraction", AuthorizedUsers(http.HandlerFunc(handleSubtraction)))
	mux.Handle("/operations/multiplication", AuthorizedUsers(http.HandlerFunc(handleMultiplication)))
	mux.Handle("/operations/division", AuthorizedUsers(http.HandlerFunc(handleDivision)))
	mux.Handle("/workers", AuthorizedUsers(http.HandlerFunc(handleWorkers)))
	mux.Handle("/results", AuthorizedUsers(http.HandlerFunc(handleResults)))
	mux.Handle("/results/", AuthorizedUsers(http.HandlerFunc(handleResultsByID)))
	mux.Handle("/database", AuthorizedUsers(http.HandlerFunc(handleDatabaseTime)))
	mux.Handle("/exit", AuthorizedUsers(http.HandlerFunc(http.HandlerFunc(handleExit))))
	auth(mux)
	ErrorLogger.Println(http.ListenAndServe(":8080", cache(delCookie(mux))))
}
func handleWorkers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		InfoLogger.Println("get workers info")
		getWorkers(w, r)
	default:
		ErrorLogger.Println("error: method not allowed (workers info)")
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not allowed")
	}
}

func getWorkers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workers)
	InfoLogger.Println("get workers info OK")
}

type Task struct {
	ID         string `json:"id"`
	Expression string `json:"expression"`
}

type Result struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}
