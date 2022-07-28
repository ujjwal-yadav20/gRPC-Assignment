package main

import (
	"context"
	"fmt"
	 "io"
	 "log"
	 "net"
	 "time"
	"github.com/ujjwal-yadav20/calculator/calculatorpb"
	"google.golang.org/grpc"

)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (resp *calculatorpb.SumResponse, err error) {
	fmt.Println("Sum Function was invoked to demonstrate unary streaming")

	num1 := req.GetSumNumbers().GetNum1()
	num2 := req.GetSumNumbers().GetNum2()

	res := num1 + num2

	resp = &calculatorpb.SumResponse{
		Result: res,
	}
	return resp, nil
}
func (*server) Prime(req *calculatorpb.PrimeRequest, resp calculatorpb.CalculatorService_PrimeServer) error {

	fmt.Println("Prime function invoked for server side streaming")

	numberM := req.GetNumberM()
	var fact int64;
	fact = 1;
	var k int64;
    for  k =2;k<numberM;k++ {
        fact = fact * (k - 1);
        if ((fact + 1) % k == 0) {
            res:= calculatorpb.PrimeResponse{
				NumberPrime : k,
			}
			time.Sleep(1000 * time.Millisecond)
			resp.Send(&res)
		}	
    }
	return nil
}
func (*server) Average(stream calculatorpb.CalculatorService_AverageServer) error {

	fmt.Println("Average Function is invoked to demonstrate client side streaming")

	var result float32
	var total float32
	var sum float32
	var temp float32
	total =0
	sum=0

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			result = sum/total
			//we have finished reading client stream
			return stream.SendAndClose(&calculatorpb.AvgResponse{

				Result: result,
			})
		}

		if err != nil {
			log.Fatalf("Error while reading client stream : %v", err)
		}
		temp = float32(msg.GetNumber())
		sum = sum + temp
		total = total + 1

	}
}
func (*server) MaxSoFar(stream calculatorpb.CalculatorService_MaxSoFarServer) error {
	fmt.Println("MaxSoFar Function is invoked to demonstrate Bi-directional streaming")
	flag := 0
	var maxReceived int64
	for {

		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("error while receiving data from MaxSoFar client : %v", err)
			return err
		}
		currNum := req.GetNumber()
		if flag==0{
			maxReceived=currNum
			flag=1
			sendErr := stream.Send(&calculatorpb.MaxResponse{
				Number: maxReceived,
			})
	
			if sendErr != nil {
				log.Fatalf("error while sending response to GreetEveryone Client : %v", err)
				return err
			}

		}
		if currNum > maxReceived {
			maxReceived =currNum
			sendErr := stream.Send(&calculatorpb.MaxResponse{
				Number: maxReceived,
			})
	
			if sendErr != nil {
				log.Fatalf("error while sending response to GreetEveryone Client : %v", err)
				return err
			}
		}
		
		// firstName := req.GetGreeting().GetFirstName()
		// lastName := req.GetGreeting().GetLastName()

		//result := "Hello " + firstName + " " + lastName + " !"

		// sendErr := stream.Send(&greetpb.GreetEveryoneResponse{
		// 	Result: result,
		// })

		// if sendErr != nil {
		// 	log.Fatalf("error while sending response to GreetEveryone Client : %v", err)
		// 	return err
		// }
	}
}

func main() {
	fmt.Println("vim-go")

	listen, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("Failed to Listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	// //Register reflection service on gRPC server
	// reflection.Register(s)

	if err = s.Serve(listen); err != nil {
		log.Fatal("failed to serve : %v", err)
	}
}


