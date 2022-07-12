package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/ViniciusSevero/fc2-grpc/pb"
)

// criando minha implementação da service stub
type UserServiceImpl struct {
	pb.UnimplementedUserServiceServer
}

// Atribuindo função a minha implementação
// classe que a função pertence - Nome da Func(parametros) (retorno)

// AddUser(context.Context, *User) (*User, error)
func (*UserServiceImpl) AddUser(ctx context.Context, req *pb.User) (*pb.User, error) {

	// Insert Database
	fmt.Println(req.Name)

	return &pb.User{
		Id:    "123",
		Name:  req.GetName(),
		Email: req.GetEmail(),
	}, nil
}

// AddUserVerbose(*User, UserService_AddUserVerboseServer) error
func (*UserServiceImpl) AddUserVerbose(req *pb.User, stream pb.UserService_AddUserVerboseServer) error {
	stream.Send(&pb.UserResultStream{
		Result: "INIT",
		User:   &pb.User{},
	})

	time.Sleep(time.Second * 2)

	stream.Send(&pb.UserResultStream{
		Result: "INSERTING",
		User:   &pb.User{},
	})

	time.Sleep(time.Second * 2)

	stream.Send(&pb.UserResultStream{
		Result: "INSERTED",
		User: &pb.User{
			Id:    "123",
			Name:  "TESTE",
			Email: "teste@teste.com.br",
		},
	})

	time.Sleep(time.Second * 2)

	stream.Send(&pb.UserResultStream{
		Result: "COMPLETED",
		User: &pb.User{
			Id:    "123",
			Name:  "TESTE",
			Email: "teste@teste.com.br",
		},
	})

	time.Sleep(time.Second * 2)

	return nil
}

// AddUsers(UserService_AddUsersServer) error
func (*UserServiceImpl) AddUsers(stream pb.UserService_AddUsersServer) error {
	users := []*pb.User{}

	for {
		recvdUser, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Users{
				Users: users,
			})
		}
		if err != nil {
			log.Fatalf("Error chunking request %v", err)
		} else {
			fmt.Println("Received: ", recvdUser)
			users = append(users, recvdUser)
		}
	}
}

// AddUsersVerbose(UserService_AddUsersVerboseServer) error
func (*UserServiceImpl) AddUsersVerbose(stream pb.UserService_AddUsersVerboseServer) error {
	for {
		recvdUser, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("End of stream")
			return nil
		}
		if err != nil {
			log.Fatalf("Error chunking request %v", err)
		} else {
			fmt.Println("Received: ", recvdUser)
			stream.Send(&pb.UserResultStream{
				Result: "RECEIVED",
				User: &pb.User{
					Id:    "",
					Name:  recvdUser.Name,
					Email: recvdUser.Email,
				},
			})
			time.Sleep(time.Second * 2)

			stream.Send(&pb.UserResultStream{
				Result: "INSERTING DB",
				User: &pb.User{
					Id:    "",
					Name:  recvdUser.Name,
					Email: recvdUser.Email,
				},
			})
			time.Sleep(time.Second * 2)

			stream.Send(&pb.UserResultStream{
				Result: "INSERTED",
				User: &pb.User{
					Id:    fmt.Sprintf("%v", rand.Int()),
					Name:  recvdUser.Name,
					Email: recvdUser.Email,
				},
			})
			time.Sleep(time.Second * 2)
		}
	}
}
