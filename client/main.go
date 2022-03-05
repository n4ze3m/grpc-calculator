package main

import (
	"context"
	"fmt"
	cl "github.com/n4ze3m/grpc-calculator/calculator"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Println(err)
	}

	defer conn.Close()

	c := cl.NewCalculatorServiceClient(conn)
	calculator(c, 1, 2, "+")
	calculator(c, 1, 2, "-")
	calculator(c, 5, 2, "*")
	calculator(c, 10, 2, "/")
}

func calculator(c cl.CalculatorServiceClient, lhs int64, rhs int64, operator string) {

	req := &cl.CalculatorRequest{
		Calculator: &cl.Calculator{
			Lhs:      lhs,
			Rhs:      rhs,
			Operator: operator,
		},
	}

	res, err := c.Calculate(context.Background(), req)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res.Result)
}
