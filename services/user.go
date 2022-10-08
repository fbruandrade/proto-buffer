package services

import (
	"context"
	"github.com/fbruandrade/proto-buffer/pb"
	"io"
	"log"
	"time"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func (s *UserService) AddUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	// Insert database
	println(req.Name)

	return &pb.User{
		Id:   req.GetId(),
		Name: req.GetName(),
		Age:  req.GetAge(),
	}, nil
}

func (s *UserService) AddUserVerbose(req *pb.User, stream pb.UserService_AddUserVerboseServer) error {

	init := stream.Send(&pb.UserResultStream{
		Status: "Init",
		User:   &pb.User{},
	})
	if init != nil {
		return init
	}

	time.Sleep(time.Second * 3)

	ins := stream.Send(&pb.UserResultStream{
		Status: "Inserting",
		User:   &pb.User{},
	})
	if ins != nil {
		return ins
	}

	time.Sleep(time.Second * 3)

	insd := stream.Send(&pb.UserResultStream{
		Status: "Inserted",
		User: &pb.User{
			Id:   req.GetId(),
			Name: req.GetName(),
			Age:  req.GetAge(),
		},
	})
	if insd != nil {
		return insd
	}

	time.Sleep(time.Second * 3)

	comp := stream.Send(&pb.UserResultStream{
		Status: "Completed",
		User: &pb.User{
			Id:   req.GetId(),
			Name: req.GetName(),
			Age:  req.GetAge(),
		},
	})
	if comp != nil {
		return comp
	}

	time.Sleep(time.Second * 3)

	return nil
}

func (s *UserService) AddUsers(stream pb.UserService_AddUsersServer) error {
	users := []*pb.User{}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Users{
				User: users,
			})
		}
		if err != nil {
			log.Fatalf("Error receiving stream: %v", err)
		}

		users = append(users, &pb.User{
			Id:   req.GetId(),
			Name: req.GetName(),
			Age:  req.GetAge(),
		})

		// Insert database
		log.Printf("Adding: %v", req.GetName())
	}
}

func (s *UserService) AddUserStreamBoth(stream pb.UserService_AddUserStreamBothServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error receiving stream: %v", err)
		}

		// Insert database
		log.Printf("Adding: %v", req.GetName())

		err = stream.Send(&pb.UserResultStream{
			Status: "Added",
			User:   req,
		})

		if err != nil {
			log.Fatalf("Error sending stream: %v", err)
		}
	}
}
