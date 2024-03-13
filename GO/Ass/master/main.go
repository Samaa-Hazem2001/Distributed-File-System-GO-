package main

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	//"strings"

	pb "Ass/AllServices" // Import the generated package

	"google.golang.org/grpc"
)

// //// global variables //////
var (
	aliveCount  map[int32]int        // Define aliveCount as a global variable
	lookupTable map[string]FileEntry // Asmaa
	lock        sync.RWMutex         // Asmaa
)

// Asmaa
type FileEntry struct {
	DataNode    string
	FilePath    string
	IsAlive     bool
	ReplicaNode []string
}

// ///////////-------------  client services (from master)  -------------///////////////
type ClientServer struct {
	pb.UnimplementedClientServiceServer
}

// ----------  Update  -----------//
func (s *ClientServer) Upload(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	// text := req.GetText()
	port := int32(8080)   //later: change it to be an unbusy port
	ipString := "ip_here" //later: change it to be the IP with the an unbusy machine
	return &pb.UpdateResponse{PortNum: port, DataNodeIp: ipString}, nil
}

// ----------  Download  -----------//
func (s *ClientServer) Download(ctx context.Context, req *pb.DownloadRequest) (*pb.DownloadResponse, error) {
	//get fileName from the client
	fileName := req.GetFileName()
	//for debuging:-
	fmt.Println("fileName to be downloaded:", fileName)

	//later: search which mahine have this file

	//send the port and ip to this machine to the client
	port := int32(3000)   //later: change it to be an unbusy port
	ipString := "ip_down" //later: change it to be the IP with the an unbusy machine
	return &pb.DownloadResponse{PortNum: port, DataNodeIp: ipString}, nil
}

// ///////////-------------  client services (from master)  -------------///////////////
type KeepersServer struct {
	pb.UnimplementedKeepersServiceServer
}

func (s *KeepersServer) KeeperDone(ctx context.Context, req *pb.KeeperDoneRequest) (*pb.KeeperDoneResponse, error) {
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

func (s *KeepersServer) Alive(ctx context.Context, req *pb.AliveRequest) (*pb.AliveResponse, error) {
	keeperId := req.GetKeeperId()

	aliveCount[keeperId] += 1

	//for debuging:-
	// Print the result
	fmt.Println("keeperId :", keeperId)

	return &pb.AliveResponse{}, nil
}

func AliveChecker(numKeepers int32) {
	ticker := time.NewTicker(8 * time.Second) // Create a ticker that ticks every 8 seconds
	defer ticker.Stop()                       // Stop the ticker when the function returns

	for {
		select {
		case <-ticker.C:

			//for debuging
			// Perform the task you want to do every 8 seconds
			fmt.Println("Alive Tracker is here")
			fmt.Println("aliveCount[i]", aliveCount[0])

			for i := int32(0); i < numKeepers; i++ { // assuming you want to initialize values for keys 0 to 9

				if aliveCount[i] == 0 {
					//edit the main lookup table --asmaa
					lock.Lock()
					// later: which file should be deleted
					delete(lookupTable, "file_name") // Delete the entry
					lock.Unlock()
					//for debuging
					fmt.Println("aliveCount with id = ", i, " is out of service now")
				}

				//reset the aliveCount for this keeper
				aliveCount[i] = 0

			}
			// Call your function or do any operation here

			// Add other cases if you need to handle other channels
		}
	}
}

/////////////////////////////// main ///////////////////////////////

// later: is there is one client at a time to the master? wla el master laz ykon 3ndha multiple ports 34an ykon fe kza client?
// ?: hwa el upload request and download request from the clients ,each one have to be in a sepearte ports?(the current assumption is yes)
func main() {
	//NOTE: map[KeyType]ValueType
	// var aliveCount map[int]int
	aliveCount = make(map[int32]int)
	numKeepers := int32(4) //?+later:change it manually or according to what?

	lookupTable = make(map[string]FileEntry) // Asmaa

	for i := int32(0); i < numKeepers; i++ { // assuming you want to initialize values for keys 0 to 9
		aliveCount[i] = 0
	}

	//-------------  client services (from master)  ------------- //
	lisUp, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("failed to listen:", err)
		return
	}

	sUp := grpc.NewServer()
	pb.RegisterClientServiceServer(sUp, &ClientServer{})

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

	//-------------  Keeper Done (step5) and Alive (from master)  ------------- //
	lisKeeper, err := net.Listen("tcp", ":8082")
	if err != nil {
		fmt.Println("failed to listen:", err)
		return
	}

	sKeeper := grpc.NewServer()
	pb.RegisterKeepersServiceServer(sKeeper, &KeepersServer{})
	// Start the gRPC server in a separate Goroutine
	go func() {
		fmt.Println("Keeper server started. Listening on port 8082...")
		if err := sKeeper.Serve(lisKeeper); err != nil {
			fmt.Println("failed to serve:", err)
		}
		done <- true // Signal that the server is done
	}()

	//----- Alive check-----//
	go AliveChecker(numKeepers)

	//fmt.Println("Server started. again")

	// Wait for the server to finish (optional)
	<-done
}
