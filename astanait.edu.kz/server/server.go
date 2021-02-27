package server

import (
	greetpb "astanait.edu.kz/Protos"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"time"
)

type Server struct {
	greetpb.UnimplementedCalculatorServiceServer
}

func (*Server) ComputeAverage(stream greetpb.CalculatorService_AverageServer) error {
	log.Println("Avg called")
	var total float32
	var count int
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			resp := &greetpb.AvgRespond{
				Avg: total / float32(count),
			}
			return stream.SendAndClose(resp)
		}
		if err != nil {
			log.Fatalf("err while Recv Avg %v", err)
			return err
		}
		log.Printf("receive req %v", req)
		total += req.GetNum()
		count++
	}
}

func (*Server) PrimeNumberDecomposition(req *greetpb.DecomposeRequest, stream greetpb.CalculatorService_DecompositionServer) error {
	fmt.Printf("Decompose function was invoked with %v \n", req)
	number := req.GetNum()
	var divide int64 = 2
	for number > 1 {
		if number%divide == 0 {
			res := &greetpb.DecomposeResponse{Decompose: divide}
			err := stream.Send(res)
			if err != nil {
				log.Fatalf("error: %v", err.Error())
			}
			number = number / divide
			time.Sleep(time.Second)
		} else {
			divide++
		}
	}
	return nil
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:8081")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	greetpb.RegisterCalculatorServiceServer(s, &Server{})
	log.Println("Server is running on port:8081")
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
