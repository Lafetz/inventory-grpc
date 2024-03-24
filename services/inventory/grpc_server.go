package inventory

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"github.com/lafetz/inventory-grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const (
	TTL     = time.Second * 8
	CheckID = "inventory_health"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServiceServer
	srv          serviceApi
	consulClient *api.Client
	port         int
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
func NewInventoryServer(srv serviceApi, port int) *InventoryServer {
	client, err := api.NewClient(&api.Config{})
	if err != nil {
		log.Fatal("couldn't create consul client", err)
	}
	return &InventoryServer{srv: srv, port: port, consulClient: client}
}
func (s *InventoryServer) Run() {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(s.port))
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}
	grpcServer := grpc.NewServer()
	//

	reflection.Register(grpcServer)
	//
	proto.RegisterInventoryServiceServer(grpcServer, s)
	fmt.Println("Starting Server on port ", s.port)
	s.Register()
	go s.updateHealthCheck()
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Grpc failed :", err)
	}

}
func (s *InventoryServer) Register() {
	fmt.Println("Registering service...")
	address, err := os.Hostname()
	if err != nil {
		log.Fatal("couldn't get address", err)
	}
	check := &api.AgentServiceCheck{
		DeregisterCriticalServiceAfter: TTL.String(),
		TTL:                            TTL.String(),
		CheckID:                        CheckID,
	}
	registeration := &api.AgentServiceRegistration{
		ID:      "grpcserver1",
		Name:    "inventoryGrpc",
		Port:    s.port,
		Address: address,
		Check:   check,
	}
	regiErr := s.consulClient.Agent().ServiceRegister(registeration)
	if regiErr != nil {
		log.Fatal("couldn't register service", regiErr)
	}

}
func (s *InventoryServer) updateHealthCheck() {
	ticker := time.NewTicker(time.Second * 5)
	for {
		err := s.consulClient.Agent().UpdateTTL(CheckID, "first_time", api.HealthPassing)
		if err != nil {
			log.Fatal("health check startup failed")
		}
		<-ticker.C
	}
}
