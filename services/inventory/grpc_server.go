package inventory

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
	"github.com/lafetz/inventory-grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServiceServer
	srv serviceApi
}

func (s *InventoryServer) AddProduct(ctx context.Context, req *proto.AddProductReq) (*proto.AddProductRes, error) {
	product, err := s.srv.AddProduct(ctx, req.Title, req.Description)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	res := &proto.Product{Id: product.Id.String(), Title: product.Title, Description: product.Description}
	return &proto.AddProductRes{Product: res}, nil
}
func (s *InventoryServer) GetProduct(ctx context.Context, req *proto.GetProductReq) (*proto.GetProductRes, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.New(codes.NotFound, ErrInvalidId.Error()).Err()
	}
	product, err := s.srv.GetProduct(ctx, id)
	if err != nil {
		return nil, status.New(codes.NotFound, err.Error()).Err()

	}
	res := &proto.Product{Id: product.Id.String(), Title: product.Title, Description: product.Description}
	return &proto.GetProductRes{Product: res}, nil
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
	//
	reflection.Register(grpcServer)
	//
	proto.RegisterInventoryServiceServer(grpcServer, s)
	fmt.Println("Starting Server on port ", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Grpc failed :", err)
	}

}
