package app

import (
	"log"
	"net"

	"github.com/bcchicr/cinemaroom/services/authn/config"
)

func Run(cfg *config.Config) {
	container, err := NewContainer(cfg)
	if err != nil {
		log.Fatalf("Failed to build container: %v", err)
	}
	defer container.Close()

	lis, err := net.Listen("tcp", ":"+cfg.Grpc.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Server running on :%s", cfg.Grpc.Port)
	if err := container.GrpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
