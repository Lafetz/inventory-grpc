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
	conn, err := grpc.Dial("localhost:8080", opts...)
	if err != nil {
		log.Fatalln("Failed to dial:", err)
	}
	defer conn.Close()

	client := proto.NewInventoryServiceClient(conn)
	g := gin.Default()
	g.GET("/products/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		req := &proto.GetProductReq{Id: id}
		res, err := client.GetProduct(ctx, req)
		if err != nil {
			st, _ := status.FromError(err)
			ctx.JSON(400, gin.H{
				"Error": st.Message(),
			})
			return
		}

		ctx.String(200, fmt.Sprint(res.Product))

	})
	g.POST("/products", func(ctx *gin.Context) {
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
		res, err := client.AddProduct(ctx, req)
		if err != nil {
			ctx.JSON(400, gin.H{
				"Error": err,
			})

			return
		}

		ctx.String(200, fmt.Sprint(res.Product))

	})
	g.Run(":3000")
}
