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
	ip_string := "ip_here" //later: change it to be the IP with the an unbusy machine
	return &pb.UpdateResponse{PortNum: port,DataNodeIp: ip_string}, nil
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
func main() {
	lis, err := net.Listen("tcp", ":8080")
if err != nil {
    fmt.Println("failed to listen:", err)
    return
}

s := grpc.NewServer()
pb.RegisterUpdateServiceServer(s, &UploadServer{})
fmt.Println("Server started. Listening on port 8080...")
if err := s.Serve(lis); err != nil {
    fmt.Println("failed to serve:", err)
}

fmt.Println("Server started. again")
// lisKeeper, err := net.Listen("tcp", ":8081")
// if err != nil {
//     fmt.Println("failed to listen:", err)
//     return
// }

// sKeeper := grpc.NewServer()
// pb.RegisterKeeperDoneServiceServer(sKeeper, &KeeperDoneServer{})
// fmt.Println("Keeper server started. Listening on port 8081...")
// if err := sKeeper.Serve(lisKeeper); err != nil {
//     fmt.Println("failed to serve:", err)
// }
}
