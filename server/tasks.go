package main

import (
	"context"
	pb "proto"
	"strconv"
)

func (s *Server) GetTaskRPC(
	ctx context.Context,
	_ *pb.GetTask,
) (*pb.GetTaskResponse, error) {
	if len(tasks) > 0 {
		task := tasks[0]
		tasks = tasks[1:]
		InfoLogger.Println("get completed")
		return &pb.GetTaskResponse{Oper1: strconv.Itoa(operations[0].ExecutionTime), Oper2: strconv.Itoa(operations[1].ExecutionTime), Oper3: strconv.Itoa(operations[2].ExecutionTime), Oper4: strconv.Itoa(operations[3].ExecutionTime), ID: task.ID, Expression: task.Expression}, nil
	}
	return &pb.GetTaskResponse{Oper1: "-1"}, nil
}
