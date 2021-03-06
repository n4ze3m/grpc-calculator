package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"

	cl "github.com/n4ze3m/grpc-calculator/calculator"
	"google.golang.org/grpc"
)

type server struct {
	cl.UnimplementedCalculatorServiceServer
}

func (*server) Calculate(ctx context.Context, in *cl.CalculatorRequest) (*cl.CalculatorResponse, error) {
	lhs := in.GetCalculator().Lhs
	rhs := in.GetCalculator().Rhs

	var result int64

	switch in.GetCalculator().Operator {
	case "+":
		result = lhs + rhs
	case "-":
		result = lhs - rhs
	case "*":
		result = lhs * rhs
	case "/":
		result = lhs / rhs
	default:
		result = 0
	}

	return &cl.CalculatorResponse{
		Result: result,
	}, nil
}


func (*server) Multiply( in *cl.MultiplicationRequest,stream cl.CalculatorService_MultiplyServer) error {
	num := in.Number

	for i := 1; i < 11; i++ {
		result := fmt.Sprintf("%d * %d = %d", num, i , int(num) * i)
		res := &cl.MultiplicationResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(3000 * time.Millisecond)
	}

	return nil
}

func (*server) Average(stream cl.CalculatorService_AverageServer) error {
	var sum int64
	var count int64
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&cl.AverageResponse{
				Result: sum / count,
			})
		}
		if err != nil {
			return err
		}
		sum += in.GetNumber()
		count++
	}
}

func (*server) Double(stream cl.CalculatorService_DoubleServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		res := &cl.DoubleResponse{
			Result: in.GetNumber() * 2,
		}
		stream.Send(res)
	}
}

func main() {
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()

	cl.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(listen); err != nil {
		panic(err)
	}
}
