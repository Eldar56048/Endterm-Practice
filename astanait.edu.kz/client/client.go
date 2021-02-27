package client

import (
	greetpb "astanait.edu.kz/Protos"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	_ "net/http"
	"time"
	_ "time"
)

func doManyTimesFromServer(c greetpb.CalculatorServiceClient) {
	ctx := context.Background()
	req := &greetpb.DecomposeRequest{Num: 120}
	stream, err := c.Decomposition(ctx, req)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer stream.CloseSend()

LOOP:
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break LOOP
		}
		if err != nil {
			log.Fatalf("Err %v", err)
		}
		log.Printf("resp:%v \n", res.GetDecompose())
	}
}

func doClientStreaming(c greetpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Average server streaming RPC...")
	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	numbers := []float32{3, 8, 12}
	for _, number := range numbers {
		fmt.Printf("Sending number: %v\n", number)
		stream.Send(&greetpb.AvgRequest{
			Num: number,
		})
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response: %v", err)
	}
	fmt.Printf("The Average is: %v\n", res.GetAvg())
}

func doCalculateAverage(c greetpb.CalculatorServiceClient) {
	ctx := context.Background()
	stream, err := c.Average(ctx)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	requests := greetpb.AvgRequest{}
	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	fmt.Printf("Avg resp: %v\n", res)
}

func main() {
	fmt.Println("Hello I'm a client")
	conn, err := grpc.Dial("localhost:8081", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("couldnt connect: %v", err)
	}
	defer conn.Close()
	con := greetpb.NewCalculatorServiceClient(conn)
	doCalculateAverage(con)
}
