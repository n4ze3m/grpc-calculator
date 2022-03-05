package main

import (
	"context"
	"net"

	cl "github.com/n4ze3m/grpc-calculator/calculator"
	"google.golang.org/grpc"
)

type server struct {
	cl.UnimplementedCalculatorServiceServer
}

func (s *server) Calculate(ctx context.Context, in *cl.CalculatorRequest) (*cl.CalculatorResponse, error) {
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
