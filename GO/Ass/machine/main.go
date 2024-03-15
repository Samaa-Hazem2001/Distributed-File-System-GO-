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

type UploadDownloadServer struct {
	pb.UnimplementedUploadDownloadFileServiceServer
}


type NotifyMachineDataTransferServer struct {
	pb.UnimplementedNotifyMachineDataTransferServiceServer
}

var callKeeperDone func(
	filename string,
	fileSize int32,
	clientIp string,
	clientPort int32,
	uploadIP string,
	uploadPortNum int32,
)

// func (s *UploadDownloadServer) Upload(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
// 	// text := req.GetText()
// 	port := int32(8080) //later: change it to be an unbusy port
// 	ip_string := "ip_here" //later: change it to be the IP with the an unbusy machine
// 	return &pb.UpdateResponse{PortNum: port,DataNodeIp: ip_string}, nil
// }

func (s *UploadDownloadServer) UploadFile(ctx context.Context, req *pb.UploadFileRequest) (*pb.UploadFileResponse, error) {
	// isUpload = true

	err := ioutil.WriteFile(req.FileName, req.File, 0644)
	if err != nil {
		log.Printf("Failed to write file: %v", err)
		return nil, err
	}

	//before returning from the function , call "callKeeperDone"
	callKeeperDone(
		req.FileName,
		int32(len(req.File)),
		req.ClientIp,
		req.ClientPortNum,
		req.DataNodeIp,
		req.PortNum,
	)

	return &pb.UploadFileResponse{}, nil
}

func (s *UploadDownloadServer) DownloadFile(ctx context.Context, req *pb.DownloadFileRequest) (*pb.DownloadFileResponse, error) {
	// isUpload = false

	// Read the file content from the disk
	//later: change filename
	print(req.FileName)
	fileContent, err := ioutil.ReadFile(req.FileName)
	if err != nil {
		log.Fatalf("%s Failed to read file: %v", req.FileName, err)
		return nil, err
	}

	// Return the file content in the response
	return &pb.DownloadFileResponse{
		File: fileContent,
	}, nil
}

func (s *NotifyMachineDataTransferServer) NotifyMachineDataTransfer(ctx context.Context, req *pb.NotifyMachineDataTransferRequest) (*pb.NotifyMachineDataTransferResponse, error) {
	sourceID := req.GetSourceId()
	fmt.Printf("Notification received to transfer data from source ID: %d\n", sourceID)
	return &pb.NotifyMachineDataTransferResponse{}, nil
}

func main() {
	// isUpload = false

	//------- act as client (client to master)  ------//
	conn, err := grpc.Dial("localhost:8082", grpc.WithInsecure()) //<=later //to asmaa : replace with final IP and Port of the master
	if err != nil {
		fmt.Println("did not connect to the master:", err)
		return
	}
	defer conn.Close()
	c := pb.NewKeepersServiceClient(conn)
	// KeeperId := 0 //later: remove it from proto file and here ?
	
	
	//------- act as sever (server to client or other keeper -for replication-) ------//

	// //later: how to listen to multiple ports ? and call the same function for any connection of them?
	// //listen to client connection or other keeper connection
	s := grpc.NewServer()
	pb.RegisterUploadDownloadFileServiceServer(s, &UploadDownloadServer{})

	// listener 1
	lis1, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println("failed to listen:", err)
		return
	}
	fmt.Println("Server started. Listening on port 8000...")

	// listener 2
	lis2, err := net.Listen("tcp", ":8001")
	if err != nil {
		fmt.Println("failed to listen:", err)
		return
	}
	fmt.Println("Server started. Listening on port 8001...")
		
	// listener 3
	lis3, err := net.Listen("tcp", ":8002")
	if err != nil {
		fmt.Println("failed to listen:", err)
		return
	}
	fmt.Println("Server started. Listening on port 8002...")
		

	go func() {	
		fmt.Println("inside go " )
		if err := s.Serve(lis1); err != nil {
			fmt.Println("failed to serve:", err)
		}
		// else if isUpload == true {
		// 	resp, err := c.KeeperDone(context.Background(), &pb.KeeperDoneRequest{FileName: filename, FileSize: int32(fileSize), PortNum: uploadPortNum, DataNodeIp: uploadIP, KeeperId: int32(KeeperId)})
		// 	if err != nil {
		// 		fmt.Println("Error calling KeeperDone:", err, resp)
		// 	}
		// }
		// fmt.Println("isUpload = ",isUpload )

	}()
	go func() {
		fmt.Println("inside go " )
		if err := s.Serve(lis2); err != nil {
			fmt.Println("failed to serve:", err)
		}
		// else if isUpload == true  {
		// 	resp, err := c.KeeperDone(context.Background(), &pb.KeeperDoneRequest{FileName: filename, FileSize: int32(fileSize), PortNum: uploadPortNum, DataNodeIp: uploadIP, KeeperId: int32(KeeperId)})
		// 	if err != nil {
		// 		fmt.Println("Error calling KeeperDone:", err, resp)
		// 	}
		// }
		// fmt.Println("isUpload = ",isUpload )
	}()
	go func() {
		fmt.Println("inside go " )
		if err := s.Serve(lis3); err != nil {
			fmt.Println("failed to serve:", err)
		}
		// else if isUpload == true  {
		// 	resp, err := c.KeeperDone(context.Background(), &pb.KeeperDoneRequest{FileName: filename, FileSize: int32(fileSize), PortNum: uploadPortNum, DataNodeIp: uploadIP, KeeperId: int32(KeeperId)})
		// 	if err != nil {
		// 		fmt.Println("Error calling KeeperDone:", err, resp)
		// 	}
		// }
		// fmt.Println("isUpload = ",isUpload )
	}()

	//we assume we will not call "keeperDone" after finishing downloading like in uploading 

	// Register the gRPC service implementation
	s22 := grpc.NewServer()
	pb.RegisterNotifyMachineDataTransferServiceServer(s22, &NotifyMachineDataTransferServer{}) //<=later:kant bt3ml eh de?

	lisMtoM, err := net.Listen("tcp", ":3000") //<=later
	if err != nil {
		fmt.Println("failed to listen:", err)
		return
	}
	fmt.Println("Server started. Listening on port 3000...")
		
	go func() {
		if err := s22.Serve(lisMtoM); err != nil {
			fmt.Println("failed to serve:", err)
		}
	}()

	callKeeperDone  = func(
		filename string,
		fileSize int32,
		clientIp string,
		clientPort int32,
		uploadIP string,
		uploadPortNum int32,
		) {
		fmt.Println("This is a nested function")
		resp, err := c.KeeperDone(context.Background(),
								 &pb.KeeperDoneRequest{
									FileName: filename, 
									FileSize: int32(fileSize),
									ClientIp: clientIp,
									ClientPortNum: clientPort,
									DataNodeIp: uploadIP,
									PortNum: uploadPortNum,
									//KeeperId: int32(KeeperId)
									})
		if err != nil {
			fmt.Println("Error calling KeeperDone:", err, resp)
		}
	}

	

	// // Concurrently send KeepersService requests to master
	// go func() {
	// 	for {
	// 		// later: change file name...
	// 		resp, err := c.KeeperDone(context.Background(), &pb.KeeperDoneRequest{FileName: filename, FileSize: int32(fileSize), PortNum: uploadPortNum, DataNodeIp: uploadIP, KeeperId: int32(KeeperId)})
	// 		if err != nil {
	// 			fmt.Println("Error calling KeeperDone:", err, resp)
	// 		}
	// 		time.Sleep(time.Second) // Adjust the frequency of sending requests
	// 	}
	// }()

	//--- Alive ---//
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop() // Stop the ticker when the function returns

		for {
			select {
			case <-ticker.C:
				// fmt.Println("Alive Ping!!")
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
