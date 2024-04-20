package main

import (
	"context"
	"fmt"
	"log"
	"net"
	pb "proto"
	"strconv"

	"google.golang.org/grpc"
)

func (s *Server) WorkerStateRPC(
	ctx context.Context,
	in *pb.WorkerState,
) (*pb.WorkerStateResp, error) {
	if len(workers) == len(in.Status) {
		for i := 0; i < len(in.Status); i++ {
			if in.Status[i] == '0' {
				workers[i].Status = "Free"
			} else {
				workers[i].Status = "Busy"
			}
		}
	} else {
		workers = []worker_info{}
		for i := 0; i < len(in.Status); i++ {
			if in.Status[i] == '0' {
				workers = append(workers, worker_info{Worker_number: strconv.Itoa(i + 1), Status: "Free"})
			} else {
				workers = append(workers, worker_info{Worker_number: strconv.Itoa(i + 1), Status: "Busy"})
			}
		}
	}
	return nil, nil
}

func startRPC() {
	host := "localhost"
	port := "5000"
	addr := fmt.Sprintf("%s:%s", host, port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		ErrorLogger.Fatal("error starting tcp listener: ", err)
	}
	log.Println("tcp listener started at port: ", port)
	grpcServer := grpc.NewServer()
	geomServiceServer := NewServer()
	pb.RegisterWorkerServiceServer(grpcServer, geomServiceServer)
	if err := grpcServer.Serve(lis); err != nil {
		ErrorLogger.Fatal("error serving grpc: ", err)
	}
}
