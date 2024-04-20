package main

import (
	"context"
	"fmt"
	"go-calculator"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	pb "client/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
var ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
var WarningLogger = log.New(os.Stderr, "WARNING: ", log.Ldate|log.Ltime)

type worker_info struct {
	Worker_number string `json:"worker_number"`
	Status        string `json:"status"`
}

var workers []worker_info

type Task struct {
	ID         string `json:"id"`
	Expression string `json:"expression"`
}
type Result struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}
type Operation struct {
	Name          string `json:"name"`
	ExecutionTime int    `json:"execution_time"`
}

var operations = []Operation{
	{Name: "Сложение", ExecutionTime: 1},
	{Name: "Вычитание", ExecutionTime: 2},
	{Name: "Умножение", ExecutionTime: 3},
	{Name: "Деление", ExecutionTime: 4},
}

func main() {
	host := "localhost"
	port := "5000"
	addr := fmt.Sprintf("%s:%s", host, port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("could not connect to grpc server: ")
		os.Exit(1)
	}
	defer conn.Close()
	grpcClient = pb.NewWorkerServiceClient(conn)
	concurrency := getConcurrency()
	for i := 0; i < concurrency; i++ {
		go startWorker()
		InfoLogger.Println("started worker with number" + strconv.Itoa(i+1))
	}
	InfoLogger.Println("Agent started with concurrency: ", concurrency)
	for {
		time.Sleep(1000)
	}
}

var grpcClient pb.WorkerServiceClient

func startWorker() {
	numb := len(workers) + 1
	numb_str := strconv.Itoa(numb)
	workers = append(workers, worker_info{Worker_number: numb_str, Status: "Free"})
	for {
		workers[numb-1].Status = "Free"
		err := WorkerState_(workers)
		if err != nil {
			ErrorLogger.Fatal("Нет подключения к серверу")
		}
		task, err := getTask_()
		if err != nil {
			WarningLogger.Println(err)
			time.Sleep(time.Second)
			continue
		}
		workers[numb-1].Status = "Busy"
		WorkerState_(workers)
		result := processTask(task)
		err = sendResult(result)
		if err != nil {
			ErrorLogger.Println(err)
		}
	}
}

func getTask_() (*Task, error) {
	task, err := grpcClient.GetTaskRPC(context.TODO(), &pb.GetTask{})
	if err != nil {
		return &Task{}, err
	}
	if task.Oper1 == "-1" {
		return &Task{}, fmt.Errorf("No available tasks")
	}
	a, _ := strconv.Atoi(task.Oper1)
	if operations[0].ExecutionTime != a {
		operations[0].ExecutionTime = a
	}
	a, _ = strconv.Atoi(task.Oper2)
	if operations[1].ExecutionTime != a {
		operations[1].ExecutionTime = a
	}
	a, _ = strconv.Atoi(task.Oper3)
	if operations[2].ExecutionTime != a {
		operations[2].ExecutionTime = a
	}
	a, _ = strconv.Atoi(task.Oper4)
	if operations[3].ExecutionTime != a {
		operations[3].ExecutionTime = a
	}
	task_res := Task{ID: task.ID, Expression: task.Expression}
	return &task_res, nil
}

func processTask(task *Task) *Result {
	result := &Result{
		ID:    task.ID,
		Value: evaluateExpression(task.Expression),
	}

	return result
}

func calculate(input string) (float64, error) {
	containsSymbols := func(str string, symbols string) bool {
		for _, symbol := range symbols {
			if strings.Contains(str, string(symbol)) {
				return true
			}
		}
		return false
	}
	if containsSymbols(input, "0123456789/*-+ ,.()") {
		a, err := calculator.Calculate(input)
		if err != nil {
			return 0, err
		}
		return a, nil
	} else {
		return 0, fmt.Errorf("Недопустимые символы")
	}
}

func evaluateExpression(expression string) string {
	input := strings.ReplaceAll(expression, " ", "")
	var t int = 0
	var q int = 0
	for _, char := range input {
		if char == 42 {
			t += operations[2].ExecutionTime
			q = 1
		}
		if char == 47 {
			t += operations[3].ExecutionTime
			q = 1
		}
		if char == 43 {
			t += operations[0].ExecutionTime
			q = 1
		}
		if char == 45 {
			t += operations[1].ExecutionTime
			q = 1
		}
	}
	if q == 0 {
		return "error, when calculating"
	}
	time.Sleep(time.Duration(t) * time.Second)
	res, err := calculate(input)
	if err != nil {
		return "error, when calculating"
	}
	fmt.Println(res)
	return fmt.Sprintf("%f", res)
}

func sendResult(result *Result) error {
	_, err := grpcClient.SendResultRPC(context.TODO(), &pb.SendResult{
		ID:     result.ID,
		Result: result.Value,
	})
	if err != nil {
		return err
	}
	return nil
}

func WorkerState_(workers2 []worker_info) error {
	str := ""
	for i := 0; i < len(workers2); i++ {
		if workers2[i].Status == "Free" {
			str = str + "0"
		} else {
			str = str + "1"
		}
	}
	_, err := grpcClient.WorkerStateRPC(context.TODO(), &pb.WorkerState{
		Status: str,
	})
	if err != nil {
		return err
	}
	return nil
}
