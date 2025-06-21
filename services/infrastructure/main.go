package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"infrastructure/config"
	"infrastructure/server"
	pb "wildfire-risk-platform/api/proto/generated"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println("Starting Infrastructure Service...")

	// Create the infrastructure server
	infrastructureServer := server.NewInfrastructureServer(cfg)

	// Create gRPC server with interceptors
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
	)

	// Register the infrastructure service
	pb.RegisterInfrastructureServiceServer(grpcServer, infrastructureServer)

	// Register health check service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("infrastructure", grpc_health_v1.HealthCheckResponse_SERVING)

	// Register reflection service for debugging
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.GRPCPort, err)
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down gracefully...")
		grpcServer.GracefulStop()
	}()

	// Start serving
	log.Printf("Infrastructure Service listening on port %s", cfg.GRPCPort)
	log.Printf("Using Overpass API at: %s", cfg.OverpassAPIURL)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("gRPC call: %s", info.FullMethod)
	resp, err := handler(ctx, req)
	if err != nil {
		log.Printf("gRPC call failed: %s - %v", info.FullMethod, err)
	}
	return resp, err
}
