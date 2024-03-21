package inventory

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/lafetz/inventory-grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServiceServer
	srv serviceApi
}

func (s *InventoryServer) AddProduct(ctx context.Context, req *proto.AddProductReq) (*proto.AddProductRes, error) {
	product, err := s.srv.AddProduct(req.Title, req.Description)
	if err != nil {
		return nil, err
	}
	res := &proto.Product{Id: product.Id.String(), Title: product.Title, Description: product.Description}
	return &proto.AddProductRes{Product: res}, nil
}
func NewInventoryServer(srv serviceApi) *InventoryServer {
	return &InventoryServer{srv: srv}
}
func (s *InventoryServer) Run(port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	proto.RegisterInventoryServiceServer(grpcServer, s)
	fmt.Println("Starting Server on port ", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Grpc faile :", err)
	}

}
