package main

import (
	"github.com/lafetz/inventory-grpc/services/inventory"
)

var PORT = ":8080"

func main() {

	repo := inventory.NewDb()
	srv := inventory.NewService(repo)
	grpcServer := inventory.NewInventoryServer(srv)

	grpcServer.Run(PORT)

}
