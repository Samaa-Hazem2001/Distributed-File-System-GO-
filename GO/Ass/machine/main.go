package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	//"net"
	//"strings"

	pb "Ass/AllServices" // Import the generated package

	"google.golang.org/grpc"
)

// type UploadDownloadServer struct {
// 	pb.UnimplementedUploadDownloadFileServiceServer
// }

// type NotifyMachineDataTransferServer struct {
// 	pb.UnimplementedNotifyMachineDataTransferServiceServer
// }

//	type TransferFileServiceServer struct {
//		pb.UnimplementedTransferFileServiceServer
//	}
type FileServer struct {
	pb.UnimplementedFileServiceServer
}

var callKeeperDone func(
	filename string,
	fileSize int32,
	clientIp string,
	clientPort int32,
	// uploadIP string,
	uploadPortNum int32,
)

var callKeeperDoneDownload func(
	// uploadIP string,
	downloadPortNum int32,
)

var callReplicationDone func(
	filename string,
	destPortNum int32,
)

var myIp string = "localhost"

// func (s *UploadDownloadServer) Upload(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
// 	// text := req.GetText()
// 	port := int32(8080)
// 	ip_string := "ip_here"
// 	return &pb.UpdateResponse{PortNum: port,DataNodeIp: ip_string}, nil
// }

func (s *FileServer) UploadFile(ctx context.Context, req *pb.UploadFileRequest) (*pb.UploadFileResponse, error) {
	// isUpload = true
	if _, err := os.Stat(myIp); os.IsNotExist(err) {
		err := os.Mkdir(myIp, 0777)
		if err != nil {
			fmt.Println("Error creating folder:", err)
		}

		fmt.Println("Folder created successfully.")
	} else if err != nil {
		fmt.Println("Error checking folder existence:", err)
	} else {
		fmt.Println("Folder already exists.")
	}
	err := ioutil.WriteFile("./"+myIp+"/"+req.FileName, req.File, 0644)
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
		// req.DataNodeIp,
		req.PortNum,
	)

	return &pb.UploadFileResponse{}, nil
}

func (s *FileServer) DownloadFile(ctx context.Context, req *pb.DownloadFileRequest) (*pb.DownloadFileResponse, error) {
	// isUpload = false

	// Read the file content from the disk
	fileContent, err := ioutil.ReadFile("./" + myIp + "/" + req.FileName)
	if err != nil {
		log.Fatalf("%s Failed to read file: %v", req.FileName, err)
		return nil, err
	}

	callKeeperDoneDownload(
		// req.DataNodeIp,
		req.PortNum,
	)

	// Return the file content in the response
	return &pb.DownloadFileResponse{
		File: fileContent,
	}, nil
}

func (s *FileServer) NotifyMachineDataTransfer(ctx context.Context, req *pb.NotifyMachineDataTransferRequest) (*pb.NotifyMachineDataTransferResponse, error) {
	filename := req.GetFileName()
	sourceMachineIp := req.GetSourceIp()
	destinationMachineIp := req.GetDistIp()
	portNum := req.GetPortNum()

	fmt.Printf("Received Notification to upload file: %s from machine %s to machine %s\n", filename, sourceMachineIp, destinationMachineIp)
	fileContent, err := ioutil.ReadFile("./" + sourceMachineIp + "/" + filename)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	err = sendFileData(destinationMachineIp, portNum, filename, fileContent)
	if err != nil {
		log.Fatalf("Failed to send file: %v", err)
	}
	return &pb.NotifyMachineDataTransferResponse{}, nil
}

func (s *FileServer) TransferFile(ctx context.Context, req *pb.TransferFileUploadRequest) (*pb.TransferFileUploadResponse, error) {
	filename := req.GetFileName()
	fileData := req.GetFile()
	portNum := req.GetPortNum()
	if _, err := os.Stat(myIp); os.IsNotExist(err) {
		err := os.Mkdir(myIp, 0777)
		if err != nil {
			fmt.Println("Error creating folder:", err)
		}
		fmt.Println("Folder created successfully.")
	} else if err != nil {
		fmt.Println("Error checking folder existence:", err)
	} else {
		fmt.Println("Folder already exists.")
	}
	err := ioutil.WriteFile("./"+myIp+"/"+filename, fileData, 0644)
	if err != nil {
		return nil, fmt.Errorf("error writing file: %v", err)
	}
	callReplicationDone(
		filename,
		portNum,
	)

	return &pb.TransferFileUploadResponse{Success: true, Message: "File transferred successfully"}, nil
}
func sendFileData(destinationMachineIp string, destPortNum int32, filename string, fileData []byte) error {
	fmt.Println(destinationMachineIp + ":" + strconv.Itoa(int(destPortNum)))
	conn, err := grpc.Dial(destinationMachineIp+":"+strconv.Itoa(int(destPortNum)), grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewFileServiceClient(conn)

	_, err = client.TransferFile(context.Background(), &pb.TransferFileUploadRequest{
		FileName: filename,
		File:     fileData,
		PortNum:  destPortNum,
	})
	if err != nil {
		return err
	}
	return nil
}
func serve(port int) {
	// Create a listener
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(int(port)))
	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
		return
	}

	// Create a gRPC server
	s := grpc.NewServer()

	// Register your gRPC services with the server
	pb.RegisterFileServiceServer(s, &FileServer{})

	// Serve gRPC requests
	if err := s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v\n", err)
		return
	}
}
func main() {
	// isUpload = false
	//some definations:
	for i, arg := range os.Args[1:] {
		fmt.Printf("Argument %d: %s\n", i+1, arg)
	}
	// var keeperPort1 int32 = 8000
	// var keeperPort2 int32 = 8001
	// var keeperPort3 int32 = 8002
	keeperPort1, _ := strconv.Atoi(os.Args[1])
	keeperPort2, _ := strconv.Atoi(os.Args[2])
	keeperPort3, _ := strconv.Atoi(os.Args[3])
	var masterPortToKeeper int32 = 8082
	// var masterIp string = "172.28.177.163"
	var masterIp string = "localhost"

	//------- act as client (client to master)  ------//
	conn, err := grpc.Dial(masterIp+":"+strconv.Itoa(int(masterPortToKeeper)), grpc.WithInsecure()) //<=later //to asmaa : replace with final IP and Port of the master
	if err != nil {
		fmt.Println("did not connect to the master:", err)
		return
	}
	defer conn.Close()
	c := pb.NewKeepersServiceClient(conn)
	// KeeperId := 0

	//------- act as sever (server to client or other keeper -for replication-) ------//

	//we assume we will not call "keeperDone" after finishing downloading like in uploading
	go serve(keeperPort1)
	go serve(keeperPort2)
	go serve(keeperPort3)

	fmt.Println("Server started on ports:", keeperPort1, keeperPort2, keeperPort3)

	callKeeperDone = func(
		filename string,
		fileSize int32,
		clientIp string,
		clientPort int32,
		// uploadIP string,
		uploadPortNum int32,
	) {
		fmt.Println("This is a nested function")
		resp, err := c.KeeperDone(context.Background(),
			&pb.KeeperDoneRequest{
				FileName:      filename,
				FileSize:      int32(fileSize),
				ClientIp:      clientIp,
				ClientPortNum: clientPort,
				DataNodeIp:    myIp,
				PortNum:       uploadPortNum,
				//KeeperId: int32(KeeperId)
			})
		if err != nil {
			fmt.Println("Error calling KeeperDone:", err, resp)
		}
	}

	callReplicationDone = func(
		filename string,
		destPortNum int32,
	) {
		fmt.Println("This is a nested function for replication")
		resp, err := c.ReplicationDone(context.Background(),
			&pb.ReplicationDoneRequest{
				FileName:   filename,
				DataNodeIp: myIp,
				PortNum:    destPortNum,
				//KeeperId: int32(KeeperId)
			})
		if err != nil {
			fmt.Println("Error calling ReplicationDone:", err, resp)
		}
	}

	callKeeperDoneDownload = func(
		downloadPortNum int32,
	) {
		fmt.Println("This is a nested function")
		resp, err := c.KeeperDoneDown(context.Background(),
			&pb.KeeperDoneDownRequest{
				DataNodeIp: myIp,
				PortNum:    downloadPortNum,
			})
		if err != nil {
			fmt.Println("Error calling KeeperDoneDown :", err, resp)
		}
	}

	//--- Alive ---//
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop() // Stop the ticker when the function returns

		for range ticker.C {

			fmt.Println("Alive Ping!!")
			resp, err := c.Alive(context.Background(), &pb.AliveRequest{DataNodeIp: myIp})
			if err != nil {
				fmt.Println("Error calling KeeperDone:", err, resp)
				return
			}

		}
	}()

	select {}
}

//for heartbeat feature, i want each keeper to send the alive signal without waiting to the respone (without waiting ,m4 btklm en el responce hykon fady)
