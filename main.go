package main

import (
	"context"
	"os/signal"

	"fmt"

	cvnet2 "github.com/DatanoiseTV/cvnet2-proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"syscall"
)

type server struct {
	cvnet2.UnimplementedCVServer
}

func (s *server ) PinMode(ctx context.Context, in *cvnet2.ConfigMessage) (*cvnet2.ConfigMessage, error){
	return nil, nil
}


func (s *server ) ReadCV(ctx context.Context, in *cvnet2.CVMessage) (*cvnet2.CVMessage, error){
	return nil, nil
}

func (s *server ) WriteCV(ctx context.Context, in *cvnet2.CVMessage) (*cvnet2.CVMessage, error){
	return nil, nil
}

func (s *server ) ReadGate(ctx context.Context, in *cvnet2.GateMessage) (*cvnet2.GateMessage, error){
	return nil, nil
}

func (s *server ) WriteGate(ctx context.Context, in *cvnet2.GateMessage) (*cvnet2.GateMessage, error){
	return nil, nil
}

func (s *server ) ReadCVStream(in *cvnet2.CVMessage, src cvnet2.CV_ReadCVStreamServer) (error){
	return nil
}

func (s *server ) WriteCVStream(src cvnet2.CV_WriteCVStreamServer) (error){
	return nil
}

func (s *server ) ReadGateStream(in *cvnet2.GateMessage, src cvnet2.CV_ReadGateStreamServer) (error){
	return nil
}

func (s *server ) WriteGateStream(src cvnet2.CV_WriteGateStreamServer) (error){
	return nil
}

func main() {
	// We need to be root to access GPIO
	if os.Getuid() != 0 {
		fmt.Println("Sorry, root required.")
		os.Exit(1)
	}

	// Set nice priority to -20 to allow low-latency output on PREEMPT_RT Linux Kernel
	syscall.Setpgid(0, 0); syscall.Setpriority(syscall.PRIO_PGRP, 0, -20)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		<-sigs
		fmt.Printf("You pressed ctrl + C. User interrupted infinite loop.")
		os.Exit(0)
	}()

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	cvnet2.RegisterCVServer(grpcServer, &server{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
