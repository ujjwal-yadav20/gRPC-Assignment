package main

import (
	"context"
	"fmt"
	"io"
	 "log"
	 //"net"
	"time"
	"github.com/ujjwal-yadav20/calculator/calculatorpb"
	"google.golang.org/grpc"

)
func main() {

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	//Sum(c)
	//Prime(c)
	//Average(c)
	MaxSoFar(c)

	
}
func Sum(c calculatorpb.CalculatorServiceClient) {

	fmt.Println("Starting to do Addition")

	req := calculatorpb.SumRequest{
		SumNumbers: &calculatorpb.Data{
			Num1: 15,
			Num2: 20,
		},
	}

	resp, err := c.Sum(context.Background(), &req)
	if err != nil {
		log.Fatalf("error while calling Sum grpc unary call: %v", err)
	}

	log.Printf("Response from Sum Unary Call : %v", resp.Result)

}
func Prime(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Staring ServerSide GRPC streaming ....")

	req := calculatorpb.PrimeRequest{

		NumberM:40,
	}

	respStream, err := c.Prime(context.Background(), &req)
	if err != nil {
		log.Fatalf("error while calling Prime server-side streaming grpc : %v", err)
	}

	for {
		msg, err := respStream.Recv()
		if err == io.EOF {
			//we have reached to the end of the file
			break
		}

		if err != nil {
			log.Fatalf("error while receving server stream : %v", err)
		}

		fmt.Println("Response From Prime Server : ", msg.NumberPrime)
	}
}
func Average(c calculatorpb.CalculatorServiceClient) {

	fmt.Println("Starting Client Side Streaming over GRPC ....")

	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalf("error occured while performing client-side streaming : %v", err)
	}

	requests := []*calculatorpb.AvgRequest{
		&calculatorpb.AvgRequest{
			Number:104,
		},
		&calculatorpb.AvgRequest{
			Number:194,
		},
		&calculatorpb.AvgRequest{
			Number:900,
		},
		&calculatorpb.AvgRequest{
			Number:7894,
		},
		&calculatorpb.AvgRequest{
			Number:1004,
		},
	}

	for _, req := range requests {
		fmt.Println("\nSending Request.... : ", req)
		stream.Send(req)
		time.Sleep(1 * time.Second)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from server : %v", err)
	}
	fmt.Println("\n****Response From Server : ", resp.GetResult())
}
func MaxSoFar(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting Bi-directional stream by calling MaxSoFar over GRPC......")

	requests := []*calculatorpb.MaxRequest{
		&calculatorpb.MaxRequest{
			Number:3,
		},
		&calculatorpb.MaxRequest{
			Number:2,
		},
		&calculatorpb.MaxRequest{
			Number:5,
		},
		&calculatorpb.MaxRequest{
			Number:1,
		},
		&calculatorpb.MaxRequest{
			Number:9,
		},
	}

	stream, err := c.MaxSoFar(context.Background())
	if err != nil {
		log.Fatalf("error occured while performing client side streaming : %v", err)
	}

	//wait channel to block receiver
	waitchan := make(chan struct{})

	go func(requests []*calculatorpb.MaxRequest) {
		for _, req := range requests {

			fmt.Println("\nSending Request..... : ", req.Number)
			err := stream.Send(req)
			if err != nil {
				log.Fatalf("error while sending request to MaxSoFarservice : %v", err)
			}
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}(requests)

	go func() {
		for {

			resp, err := stream.Recv()
			if err == io.EOF {
				close(waitchan)
				return
			}

			if err != nil {
				log.Fatalf("error receiving response from server : %v", err)
			}

			fmt.Printf("\nResponse From Server : %v", resp.GetNumber())
		}
	}()

	//block until everything is finished
	<-waitchan
}

