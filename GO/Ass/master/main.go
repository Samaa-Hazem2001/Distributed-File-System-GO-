package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"sync"
	"time"

	//"strings"

	pb "Ass/AllServices" // Import the generated package

	"google.golang.org/grpc"
)

// //// global variables //////
var (
	aliveCount     map[int]int // Define aliveCount as a global variable
	lock           sync.RWMutex
	machineMap     map[int]Machine
	filenameMap    map[string][]string
	replicationMap map[string]map[string]bool
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
	print("inside master upload")
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
					//later: look on the machineMap here?
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
	// defer lock.Unlock()
	for _, machine := range machineMap {
		if machine.IP == DataNodeIp {
			machine.FileNames = append(machine.FileNames, fileName)
			machineMap[machine.ID] = machine
		}
	}
	lock.Unlock()

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
		fmt.Println("Error calling DoneUp:", err, resp)
		return
	}

}
func setPortStatus(DataNodeIp string, portNumber int, isBusy bool) error {
	for _, machine := range machineMap {
		if machine.IP == DataNodeIp {
			for i, port := range machine.Ports {
				if port.Port == portNumber {
					//later: look on the machineMap here?
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
	keeperIP := req.GetDataNodeIp()

	//from ip get id
	var keeperId int
	for _, machine := range machineMap {
		if machine.IP == keeperIP {
			keeperId = machine.ID
			break
		}
	}

	aliveCount[keeperId] += 1

	//for debuging:-
	// Print the result
	fmt.Println("keeperIP :", keeperIP, " has come alive.")

	return &pb.AliveResponse{}, nil
}

func AliveChecker(numKeepers int) {
	ticker := time.NewTicker(8 * time.Second) // Create a ticker that ticks every 8 seconds
	defer ticker.Stop()                       // Stop the ticker when the function returns

	for {
		select {
		case <-ticker.C:

			//for debuging
			// Perform the task you want to do every 8 seconds
			// fmt.Println("Alive Tracker is here")
			// fmt.Println("aliveCount[i]", aliveCount[0])

			for i := int(0); i < numKeepers; i++ { // assuming you want to initialize values for keys 0 to 9

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

func (s *KeepersServer) ReplicationDone(ctx context.Context, req *pb.ReplicationDoneRequest) (*pb.ReplicationDoneResponse, error) {
	keeperIp := req.GetDataNodeIp()
	portNum := req.GetPortNum()
	fileName := req.GetFileName()

	//for debuging:-
	// Print the result
	fmt.Println("keeperIp :", keeperIp, " finished the replication for file : ", fileName)

	replicationMap[fileName][keeperIp] = true

	//mark this "portNum" as an avialble port to the machine with ip = keeperIp
	err := setPortStatus(keeperIp, int(portNum), false)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return &pb.ReplicationDoneResponse{}, nil
}

func replicationFinishChecker() {
	ticker := time.NewTicker(28 * time.Second) //NOTE: Create a ticker that ticks every 28 seconds to break the tie with 10 seconds of the replicationchecker
	defer ticker.Stop()                        // Stop the ticker when the function returns

	for {
		select {
		case <-ticker.C:

			for fileName, machine_lists := range replicationMap { //iterate over files

				// for debuging:
				fmt.Printf("fileName: %s, machine_lists: %d\n", fileName, machine_lists)

				for currentIp, done := range machine_lists { //iterate over machines

					if done {
						continue
					}

					//later: if not done, then go to the lookup table and delete that machine for this file name
					if ips, ok := filenameMap[fileName]; ok {
						indexToRemove := -1
						for i, ip := range ips {
							if ip == currentIp {
								indexToRemove = i
								break
							}
						}
						if indexToRemove != -1 {
							filenameMap[fileName] = append(ips[:indexToRemove], ips[indexToRemove+1:]...)
						}
					}

				}

			}

			//reset replicationMap
			replicationMap = make(map[string]map[string]bool)
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

			for filename, machineIps := range filenameMap {
				fmt.Printf("%s: %v\n", filename, machineIps)
				sourceMachineIp := machineIps[0]
				machineIpsLen := len(machineIps)

				for machineIpsLen < 3 {
					//later: for testing in group of laptops , uncomment this
					// destinationMachineIp, destinationMachineId, destMachinePort, err := selectMachineToCopyTo(filename) later un comment this
					_, destinationMachineId, destMachinePort, err := selectMachineToCopyTo(filename)
					destinationMachineIp := "localhost"
					if err != nil {
						fmt.Println("Error: ", err)
					}
					notifyMachineDataTransfer(sourceMachineIp, destinationMachineIp, destMachinePort, filename)
					machineIpsLen++
					filenameMap[filename] = append(filenameMap[filename], destinationMachineIp)

					//samaa:
					if replicationMap[filename] == nil {
						replicationMap[filename] = make(map[string]bool)
					}
					replicationMap[filename][destinationMachineIp] = false
					machine := machineMap[destinationMachineId]
					machine.FileNames = append(machine.FileNames, filename)
					machineMap[destinationMachineId] = machine
				}
			}
		}
	}
}
func selectMachineToCopyTo(filename string) (string, int, int, error) {
	machineIps := filenameMap[filename]
	for _, machine := range machineMap {
		if machine.IsAlive {
			found := false
			for _, ip := range machineIps {
				if ip == machine.IP {
					found = true
					break
				}
			}
			if !found {
				//later:asmaa change the dheck for a machine to have unbusy port
				//later:asmaa change the port num
				return machine.IP, machine.ID, 3000, nil
			}
		}
	}
	return "", 0, 0, fmt.Errorf("failed to find machine")
}
func notifyMachineDataTransfer(sourceMachineIp string, destinationMachineIp string, destMachinePort int, filename string) error {
	//later: from the ip get the id
	var machineID int
	for _, machine := range machineMap {
		if machine.IP == sourceMachineIp {
			machineID = machine.ID
			break
		}
	}

	var nonBusyPort int
	machine := machineMap[machineID]
	//later: look on the machineMap here?
	if machine.IsAlive { //later: do we have to delete this unnecessary condition?
		for i, port := range machine.Ports {
			if !port.Busy {
				machineMap[machineID].Ports[i].Busy = true
				nonBusyPort = port.Port
				break
			}
		}
	}

	conn, err := grpc.Dial(sourceMachineIp+":"+strconv.Itoa(int(nonBusyPort)), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewNotifyMachineDataTransferServiceClient(conn)

	_, err = client.NotifyMachineDataTransfer(context.Background(), &pb.NotifyMachineDataTransferRequest{
		SourceIp: sourceMachineIp,
		DistIp:   destinationMachineIp,
		PortNum:  int32(destMachinePort),
		FileName: filename,
	})
	if err != nil {
		return fmt.Errorf("failed to notify machine for data transfer: %v", err)
	}

	//remark the nonBusyPort of the source as non busy again
	err = setPortStatus(sourceMachineIp, int(nonBusyPort), false)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("Done")

	return nil
}
func generateFilenameMap() map[string][]string {
	filenameMap := make(map[string][]string)

	for _, machine := range machineMap {
		if machine.IsAlive {
			machineIps := make(map[int]bool)
			for _, filename := range machine.FileNames {
				if !machineIps[machine.ID] {
					filenameMap[filename] = append(filenameMap[filename], machine.IP)
					machineIps[machine.ID] = true
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

	//later: look on the machineMap here?
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
	aliveCount = make(map[int]int)
	numKeepers := int(4) //?+later:change it manually or according to what?

	for i := int(0); i < numKeepers; i++ { // assuming you want to initialize values for keys 0 to 9
		aliveCount[i] = 0
	}

	//-------------  client services (from master)  ------------- //
	lisUp, err := net.Listen("tcp", ":8081")
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
	go replicationFinishChecker()

	//fmt.Println("Server started. again")

	// Wait for the server to finish (optional)
	<-done
}
