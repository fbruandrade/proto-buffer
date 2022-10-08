package main

import (
	"context"
	"github.com/fbruandrade/proto-buffer/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"time"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to gRPC server: %v", err)
	}
	defer func(connection *grpc.ClientConn) {
		err := connection.Close()
		if err != nil {

		}
	}(connection)

	client := pb.NewUserServiceClient(connection)

	//AddUser(client) AddUserVerbose(client) AddUsers(client)
	AddUserStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:   "123456",
		Name: "Bruno",
		Age:  25,
	}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	log.Printf("Response from server: %v", res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:   "123456",
		Name: "Bruno",
		Age:  25,
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	for {
		msg, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not receive the message: %v", err)
		}

		log.Printf("Response from server: %v %v", msg.Status, msg.User)
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		{
			Id:   "123456",
			Name: "Felipe",
			Age:  40,
		},
		{
			Id:   "123457",
			Name: "Bruno",
			Age:  25,
		},
		{
			Id:   "123458",
			Name: "Charles",
			Age:  64,
		},
		{
			Id:   "123459",
			Name: "Charlene",
			Age:  36,
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	for _, req := range reqs {
		err := stream.Send(req)
		if err != nil {
			return
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Could not receive the message: %v", err)
	}

	log.Printf("Response from server: %v", res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {
	reqs := []*pb.User{
		{
			Id:   "123456",
			Name: "Felipe",
			Age:  40,
		},
		{
			Id:   "123457",
			Name: "Bruno",
			Age:  25,
		},
		{
			Id:   "123458",
			Name: "Charles",
			Age:  64,
		},
		{
			Id:   "123459",
			Name: "Charlene",
			Age:  36,
		},
	}

	stream, err := client.AddUserStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	wait := make(chan struct{})

	go func() {
		for _, req := range reqs {
			log.Printf("Sending user: %v", req)
			err := stream.Send(req)
			if err != nil {
				return
			}
			time.Sleep(time.Second * 2)
		}
		err := stream.CloseSend()
		if err != nil {
			return
		}
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Could not receive the message: %v", err)
			}

			log.Printf("Response from server: %v with status: %v", res.GetUser(), res.Status)
		}
		close(wait)
	}()

	<-wait
}
