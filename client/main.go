package main

import (
	"context"
	"fmt"
	"io"
	"time"

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
	// unary
	calculator(c, 1, 2, "+")
	calculator(c, 1, 2, "-")
	calculator(c, 5, 2, "*")
	calculator(c, 10, 2, "/")
	// server stream
	multiply(c, 5)
	multiply(c, 7)
	// client stream
	average(c)
	// bi directional
	double(c)
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

func multiply(c cl.CalculatorServiceClient, num int64) {
	req := &cl.MultiplicationRequest{
		Number: num,
	}

	stream, err := c.Multiply(context.Background(), req)

	if err != nil {
		fmt.Println(err)
	}

	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(msg.Result)
	}
}

func average(c cl.CalculatorServiceClient) {
	stream, err := c.Average(context.Background())

	if err != nil {
		fmt.Println(err)
	}

	for i := 1; i < 11; i++ {
		msg := &cl.AverageRequest{
			Number: int64(i),
		}

		stream.Send(msg)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res.Result)
}

func double(c cl.CalculatorServiceClient) {

	stream, err := c.Double(context.Background())

	if err != nil {
		fmt.Println(err)
	}

	// channel
	ch := make(chan struct{})

	// send

	go func() {
		for i := 1; i < 11; i++ {
			msg := &cl.DoubleRequest{
				Number: int64(i),
			}

			stream.Send(msg)
			time.Sleep( 3000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	// recieve
	go func() {
		for {
			msg, err := stream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(msg.Result)
		}
		close(ch)
	}()

	<-ch
}
