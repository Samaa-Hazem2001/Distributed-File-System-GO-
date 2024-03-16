package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"sync"
	"time"
	"strconv"

	//"strings"

	pb "Ass/AllServices" // Import the generated package

	"google.golang.org/grpc"
)

// //// global variables //////
var (
	aliveCount  map[int32]int // Define aliveCount as a global variable
	lock        sync.RWMutex
	machineMap  map[int]Machine
	filenameMap map[string][]int
)

type PortInfo struct {
	Port int  `json:"port"`
	Busy bool `json:"busy"`
}

// later: should we add filepath
type Machine struct {
	ID        int        `json:"id"`
	IP        string     `json:"ip"`
	Ports     []PortInfo `json:"ports"`
	FileNames []string
	IsAlive   bool
}

type Config struct {
	Machines []Machine `json:"machines"`
}

func loadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// ///////////-------------  client services (from master)  -------------///////////////
type ClientServer struct {
	pb.UnimplementedClientServiceServer
}

// ----------  Update  -----------//
func (s *ClientServer) Upload(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	PortNum, DataNodeIp, err := findNonBusyPort()
	if err != nil {
		return nil, err
	}
	return &pb.UpdateResponse{PortNum: int32(PortNum), DataNodeIp: DataNodeIp}, nil
}
func findNonBusyPort() (int, string, error) {
	for _, machine := range machineMap {
		if machine.IsAlive {
			for i, port := range machine.Ports {
				if !port.Busy {
					machineMap[machine.ID].Ports[i].Busy = true
					return port.Port, machine.IP, nil
				}
			}
		}
	}
	return 0, "", errors.New("no available non-busy port found")
}

// ----------  Download  -----------//
func (s *ClientServer) Download(ctx context.Context, req *pb.DownloadRequest) (*pb.DownloadResponse, error) {
	//get fileName from the client
	fileName := req.GetFileName()
	//for debuging:-
	fmt.Println("fileName to be downloaded:", fileName)

	//later: search which mahine have this file
	PortNum, DataNodeIp, err := findNonBusyPort()
	if err != nil {
		return nil, err
	}
	return &pb.DownloadResponse{PortNum: int32(PortNum), DataNodeIp: DataNodeIp}, nil
}

// ///////////-------------  client services (from master)  -------------///////////////
type KeepersServer struct {
	pb.UnimplementedKeepersServiceServer
}

func (s *KeepersServer) KeeperDone(ctx context.Context, req *pb.KeeperDoneRequest) (*pb.KeeperDoneResponse, error) {
	fileName := req.GetFileName()
	// fileSize := req.GetFileSize()
	freePortNum := req.GetPortNum()
	// keeperId := req.GetKeeperId()
	DataNodeIp := req.GetDataNodeIp()

	clientPort := req.GetClientPortNum()
	clientIp := req.GetClientIp()

	//ConfirmClient(clientIp, clientPort) 

	// later: what about ip? is it the same as Id
	err := setPortStatus(DataNodeIp, int(freePortNum), false)
	if err != nil {
		fmt.Println("Error:", err)
	}

	//for debuging:-
	// Print the result
	fmt.Println("fileName :", fileName)
	// fmt.Println("fileSize :", fileSize)
	fmt.Println("freePortNum :", freePortNum)
	fmt.Println("DataNodeIp :", DataNodeIp)
	fmt.Println("clientPort :", clientPort)
	fmt.Println("clientIp :", clientIp)

	// later: 5-The master tracker then adds the file record to the main look-up table.

	lock.Lock()
	defer lock.Unlock()
	for _, machine := range machineMap {
		if machine.IP == DataNodeIp {
			machine.FileNames = append(machine.FileNames, fileName)
			machineMap[machine.ID] = machine
		}
	}

	// later: 6-The master will notify the client with a successful message.
	return &pb.KeeperDoneResponse{}, nil
}
func ConfirmClient(ip string, port int32) {
	conn, err := grpc.Dial(ip+":"+strconv.Itoa(int(port)), grpc.WithInsecure())
	if err != nil {
		fmt.Println("Failed to connect to client with IP:", ip, ":", err)
		return
	}

	defer conn.Close()
	c := pb.NewDoneUpServiceClient(conn)

	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// _, err = c.DoneUpload(ctx, &pb.EmptyMessage{})
	// cancel()
	// if err != nil {
	// 	return errors.New("RPC failed: " + err.Error())
	// }
	// return nil

	resp, err := c.DoneUp(context.Background(), &pb.DoneUpRequest{})
	if err != nil {
		fmt.Println("Error calling DoneUp:", err,resp)
		return
	}

}
func setPortStatus(DataNodeIp string, portNumber int, isBusy bool) error {
	for _, machine := range machineMap {
		if machine.IP == DataNodeIp {
			for i, port := range machine.Ports {
				if port.Port == portNumber {
					machine.Ports[i].Busy = isBusy
					machineMap[machine.ID] = machine
					return nil
				}
			}
		}
	}

	return fmt.Errorf("port %d not found in machine with IP %s", portNumber, DataNodeIp)
}

func (s *KeepersServer) Alive(ctx context.Context, req *pb.AliveRequest) (*pb.AliveResponse, error) {
	keeperId := req.GetKeeperId()

	aliveCount[keeperId] += 1

	//for debuging:-
	// Print the result
	// fmt.Println("keeperId :", keeperId)

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
			// fmt.Println("Alive Tracker is here")
			// fmt.Println("aliveCount[i]", aliveCount[0])

			for i := int32(0); i < numKeepers; i++ { // assuming you want to initialize values for keys 0 to 9

				if aliveCount[i] == 0 {
					//edit the main lookup table --asmaa
					// mark as dead
					lock.Lock()
					machine := machineMap[int(i)]
					machine.IsAlive = false
					machineMap[int(i)] = machine
					lock.Unlock()
					//for debuging
					// fmt.Println("aliveCount with id = ", i, " is out of service now")
				} else {
					// else mark as alive
					lock.Lock()
					machine := machineMap[int(i)]
					machine.IsAlive = true
					machineMap[int(i)] = machine
					lock.Unlock()
					//reset the aliveCount for this keeper
					aliveCount[i] = 0
				}

			}
			// Call your function or do any operation here

			// Add other cases if you need to handle other channels
		}
	}
}

func replicationChecker() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:

			filenameMap = generateFilenameMap()

			for filename, machineIDs := range filenameMap {
				fmt.Printf("%s: %v\n", filename, machineIDs)
				sourceMachineId := machineIDs[0]
				machineIDsLen := len(machineIDs)

				for machineIDsLen < 3 {
					destinationMachineId, err := selectMachineToCopyTo(filename)
					if err != nil {
						fmt.Println("Error: ", err)
					}
					notifyMachineDataTransfer(sourceMachineId, destinationMachineId)
					machineIDsLen++
					filenameMap[filename] = append(filenameMap[filename], destinationMachineId)
				}
			}
		}
	}
}
func selectMachineToCopyTo(filename string) (int, error) {
	machineIDs := filenameMap[filename]
	for _, machine := range machineMap {
		if machine.IsAlive {
			found := false
			for _, id := range machineIDs {
				if id == machine.ID {
					found = true
					break
				}
			}
			if !found {
				return machine.ID, nil
			}
		}
	}
	return -1, fmt.Errorf("failed to find machine")
}
func notifyMachineDataTransfer(sourceMachineId int, destinationMachineId int) error {
	conn, err := grpc.Dial("localhost:3000", grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewNotifyMachineDataTransferServiceClient(conn)

	_, err = client.NotifyMachineDataTransfer(context.Background(), &pb.NotifyMachineDataTransferRequest{
		SourceId: int32(sourceMachineId),
		DistId:   int32(destinationMachineId),
	})
	if err != nil {
		return fmt.Errorf("failed to notify machine for data transfer: %v", err)
	}
	fmt.Println("Done")

	return nil
}
func generateFilenameMap() map[string][]int {
	filenameMap := make(map[string][]int)

	for _, machine := range machineMap {
		if machine.IsAlive {
			machineIDs := make(map[int]bool)
			for _, filename := range machine.FileNames {
				if !machineIDs[machine.ID] {
					filenameMap[filename] = append(filenameMap[filename], machine.ID)
					machineIDs[machine.ID] = true
				}
			}
		}
	}

	return filenameMap
}

/////////////////////////////// main ///////////////////////////////

// later: is there is one client at a time to the master? wla el master laz ykon 3ndha multiple ports 34an ykon fe kza client?
// ?: hwa el upload request and download request from the clients ,each one have to be in a sepearte ports?(the current assumption is yes)
func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}
	machineMap = make(map[int]Machine)
	for _, machine := range config.Machines {
		machineMap[machine.ID] = machine
		machine.FileNames = make([]string, 0)
		// later change it depend for what?
		machine.IsAlive = true
	}
	// printing
	// fmt.Println("Machines:")
	// for id, machine := range machineMap {
	// 	fmt.Printf("ID: %d\n", id)
	// 	fmt.Printf("  IP: %s\n", machine.IP)
	// 	fmt.Printf("  FileNames: %s\n", machine.FileNames)
	// 	fmt.Printf("  IsAlive: %t\n", machine.IsAlive)
	// 	fmt.Println("  Ports:")
	// 	for _, port := range machine.Ports {
	// 		fmt.Printf("    %d (%t)\n", port.Port, port.Busy)
	// 	}
	// 	fmt.Println()
	// }

	//NOTE: map[KeyType]ValueType
	// var aliveCount map[int]int
	aliveCount = make(map[int32]int)
	numKeepers := int32(4) //?+later:change it manually or according to what?

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
		fmt.Println("Client server started. Listening on port 8080...")
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
	go replicationChecker()

	//fmt.Println("Server started. again")

	// Wait for the server to finish (optional)
	<-done
}
