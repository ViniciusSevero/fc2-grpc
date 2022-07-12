package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/ViniciusSevero/fc2-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Coult not connect to GRPC Server %v", err)
	}
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	// AddUser(client)
	// AddUserVerbose(client)
	// AddUsers(client)
	AddUsersVerbose(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "1",
		Name:  "teste",
		Email: "teste@teste.com.br",
	}

	resp, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Coult not make GRPC request %v", err)
	}

	fmt.Printf("Response: %v \n", resp)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "1",
		Name:  "teste",
		Email: "teste@teste.com.br",
	}

	stream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Coult not make GRPC request %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error chunking response %v", err)
		} else {
			log.Println(resp.Result, " - ", resp.User)
		}
	}
}

func AddUsers(client pb.UserServiceClient) {
	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Could not adquire stream to send %v", err)
	}
	for i := 1; i < 4; i++ {
		fmt.Println("Sending ", i, "...")
		stream.Send(&pb.User{
			Id:    fmt.Sprintf("%v", i),
			Name:  fmt.Sprintf("Name %v", i),
			Email: fmt.Sprintf("Email %v", i),
		})

		time.Sleep(time.Second * 2)
	}

	respFromServer, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error from server %v", err)
	}

	fmt.Println(respFromServer)

}

func AddUsersVerbose(client pb.UserServiceClient) {
	stream, err := client.AddUsersVerbose(context.Background())
	if err != nil {
		log.Fatalf("Could not adquire stream to send %v", err)
	}

	wait := make(chan int)

	go func() {
		for i := 1; i < 4; i++ {
			fmt.Println("Sending ", i, "...")
			stream.Send(&pb.User{
				Id:    "",
				Name:  fmt.Sprintf("Name %v", i),
				Email: fmt.Sprintf("Email %v", i),
			})

			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("End of stream")
				break
			}
			if err != nil {
				log.Fatalf("Could not receive msg from server %v", err)
			}
			fmt.Println(req.Result, " - ", req.User.Id, ":", req.User.Name)
		}

		close(wait)
	}()

	<-wait
}
