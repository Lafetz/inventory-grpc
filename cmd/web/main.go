package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/lafetz/inventory-grpc/proto"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type createProductBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(grpc_retry.WithCodes(codes.Internal), grpc_retry.WithMax(5), grpc_retry.WithBackoff(grpc_retry.BackoffLinear(time.Second)))))
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	//
	conn, err := grpc.Dial("localhost:8080", opts...)
	//
	if err != nil {
		log.Fatalln("Failed to dial:", err)
	}
	defer conn.Close()
	s := &server{client: proto.NewInventoryServiceClient(conn)}
	g := gin.Default()
	g.GET("/products/:id", s.getHandler)
	g.POST("/products", s.postHandler)
	g.Run(":3000")
}

func newCb() *gobreaker.CircuitBreaker {
	return gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "inventory-call",
		MaxRequests: 1,
		Timeout:     10,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Printf("Circuit Breaker: %s, changed from %v, to %v", name, from, to)
		},
	})
}

type server struct {
	client proto.InventoryServiceClient
}

func (s *server) getHandler(ctx *gin.Context) {
	//Circuit Breaker
	cb := newCb()
	//
	id := ctx.Param("id")
	req := &proto.GetProductReq{Id: id}
	res, err := cb.Execute(func() (interface{}, error) {
		return s.client.GetProduct(ctx, req)
	}) // client.GetProduct(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		ctx.JSON(400, gin.H{
			"Error": st.Message(),
		})
		return
	}
	ctx.String(200, fmt.Sprint(res.(*proto.GetProductRes).Product))

}
func (s *server) postHandler(ctx *gin.Context) {
	var productReq createProductBody
	if err := ctx.ShouldBindJSON(&productReq); err != nil {
		_, ok := err.(validator.ValidationErrors)
		if ok {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"Errors": err,
			})
			return

		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": "Error processing request body",
		})
		return
	}
	req := &proto.AddProductReq{Title: productReq.Title, Description: productReq.Description}
	res, err := s.client.AddProduct(ctx, req)
	if err != nil {
		ctx.JSON(400, gin.H{
			"Error": err,
		})

		return
	}

	ctx.String(200, fmt.Sprint(res.Product))

}
