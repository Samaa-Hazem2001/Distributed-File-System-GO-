package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"

	//"net"
	//"strings"

	pb "Ass/AllServices" // Import the generated package

	"google.golang.org/grpc"
)

type UploadServer struct {
	pb.UnimplementedUploadFileServiceServer
}

// func (s *UploadServer) Upload(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
// 	// text := req.GetText()
// 	port := int32(8080) //later: change it to be an unbusy port
// 	ip_string := "ip_here" //later: change it to be the IP with the an unbusy machine
// 	return &pb.UpdateResponse{PortNum: port,DataNodeIp: ip_string}, nil
// }

// UploadFile implements the UploadFileService.UploadFile method
func (s *UploadServer) UploadFile(ctx context.Context, req *pb.UploadFileRequest) (*pb.UploadFileResponse, error) {
	// later: what should I name the file
	err := ioutil.WriteFile("uploaded_file.mp4", req.File, 0644)
	if err != nil {
		log.Printf("Failed to write file: %v", err)
		return nil, err
	}
	return &pb.UploadFileResponse{}, nil
}

func main() {

	//------- act as sever (server to client or other keeper -for replication-) ------//

	// //later: how to listen to multiple ports ? and call the same function for any connection of them?
	// //listen to client connection or other keeper connection
	s := grpc.NewServer()
	pb.RegisterUploadFileServiceServer(s, &UploadServer{})

	go func() {
		lis, err := net.Listen("tcp", ":3000")
		if err != nil {
			fmt.Println("failed to listen:", err)
			return
		}
		fmt.Println("Server started. Listening on port 8080...")
		if err := s.Serve(lis); err != nil {
			fmt.Println("failed to serve:", err)
		}
	}()

	//------- act as client (client to master)  ------//
	conn, err := grpc.Dial("localhost:8082", grpc.WithInsecure())
	if err != nil {
		fmt.Println("did not connect:", err)
		return
	}
	defer conn.Close()
	c := pb.NewKeepersServiceClient(conn)

	// Concurrently send KeepersService requests to master
	go func() {
		for {
			// later: change file name...
			resp, err := c.KeeperDone(context.Background(), &pb.KeeperDoneRequest{FileName: "uploaded_file.mp4", FileSize: int32(2000), PortNum: int32(9999), KeeperId: int32(0)})
			if err != nil {
				fmt.Println("Error calling KeeperDone:", err, resp)
			}
			time.Sleep(time.Second) // Adjust the frequency of sending requests
		}
	}()

	//--- Alive ---//
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop() // Stop the ticker when the function returns

		for {
			select {
			case <-ticker.C:
				fmt.Println("Alive Ping!!")
				resp, err := c.Alive(context.Background(), &pb.AliveRequest{KeeperId: int32(0)})
				if err != nil {
					fmt.Println("Error calling KeeperDone:", err, resp)
					return
				}
			}
		}
	}()

	select {}
}

//for heartbeat feature, i want each keeper to send the alive signal without waiting to the respone (without waiting ,m4 btklm en el responce hykon fady)
