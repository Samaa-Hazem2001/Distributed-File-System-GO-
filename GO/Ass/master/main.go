package main

import (
	"context"
	"fmt"
	"net"
	//"strings"

	pb "Ass/AllServices" // Import the generated package

	"google.golang.org/grpc"
)

//----------  Update  -----------//
type UploadServer struct {
	pb.UnimplementedUpdateServiceServer
}

func (s *UploadServer) Upload(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	// text := req.GetText()
	port := int32(8080) //later: change it to be an unbusy port
	ipString := "ip_here" //later: change it to be the IP with the an unbusy machine
	return &pb.UpdateResponse{PortNum: port,DataNodeIp: ipString}, nil
}

//----------  Download  -----------//
type DownloadServer struct {
	pb.UnimplementedDownloadServiceServer
}

func (s *DownloadServer) Download(ctx context.Context, req *pb.DownloadRequest) (*pb.DownloadResponse, error) {
	//get fileName from the client
	fileName := req.GetFileName()
	//for debuging:-
	fmt.Println("fileName to be downloaded:", fileName)


	//later: search which mahine have this file

	//send the port and ip to this machine to the client
	port := int32(8080) //later: change it to be an unbusy port
	ipString := "ip_here" //later: change it to be the IP with the an unbusy machine
	return &pb.DownloadResponse{PortNum: port,DataNodeIp: ipString}, nil
}

//----------  Keeper Done (step5)  -----------//
type KeeperDoneServer struct {
	pb.UnimplementedKeeperDoneServiceServer  
}

func (s *KeeperDoneServer) KeeperDone(ctx context.Context, req *pb.KeeperDoneRequest) (*pb.KeeperDoneResponse, error) {
	// text := req.GetText()
	fileName := req.GetFileName()
	fileSize := req.GetFileSize()
	freePortNum := req.GetPortNum() 
	keeperId := req.GetKeeperId()

	//for debuging:-
	// Print the result
	fmt.Println("fileName :", fileName)
	fmt.Println("fileSize :", fileSize)
	fmt.Println("freePortNum :", freePortNum)
	fmt.Println("keeperId :", keeperId)
 
	return &pb.KeeperDoneResponse{}, nil
}	


//----------------  main  -----------//

// later: is there is one client at a time to the master? wla el master laz ykon 3ndha multiple ports 34an ykon fe kza client?
//?: hwa el upload request and download request from the clients ,each one have to be in a sepearte ports?(the current assumption is yes)
func main() {

	//----------------  upload file   -------------//
	lisUp, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("failed to listen:", err)
		return
	}

	sUp := grpc.NewServer()
	pb.RegisterUpdateServiceServer(sUp, &UploadServer{})

	// Create a channel to signal when the server is done
	done := make(chan bool)

	// Start the gRPC server in a separate Goroutine
	go func() {
		fmt.Println("Keeper server started. Listening on port 8080...")
		if err := sUp.Serve(lisUp); err != nil {
			fmt.Println("failed to serve:", err)
		}
		done <- true // Signal that the server is done
	}()

	// fmt.Println("Server started. again")

	//----------------  download file   -------------//
	lisDown, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println("failed to listen:", err)
		return
	}

	sDown := grpc.NewServer()
	pb.RegisterDownloadServiceServer(sDown, &DownloadServer{})
	// Start the gRPC server in a separate Goroutine
	go func() {
		fmt.Println("Keeper server started. Listening on port 8081...")
		if err := sDown.Serve(lisDown); err != nil {
			fmt.Println("failed to serve:", err)
		}
		done <- true // Signal that the server is done
	}()


	//----------------  machine done uploading file   -------------//
	lisKeeper, err := net.Listen("tcp", ":8082")
	if err != nil {
	    fmt.Println("failed to listen:", err)
	    return
	}

	sKeeper := grpc.NewServer()
	pb.RegisterKeeperDoneServiceServer(sKeeper, &KeeperDoneServer{})
	// Start the gRPC server in a separate Goroutine
	go func() {
		fmt.Println("Keeper server started. Listening on port 8082...")
		if err := sKeeper.Serve(lisKeeper); err != nil {
			fmt.Println("failed to serve:", err)
		}
		done <- true // Signal that the server is done
	}()

	






	
	// Wait for the server to finish (optional)
	<-done
}
