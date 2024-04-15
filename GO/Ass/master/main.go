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
	aliveCount map[int]int // Define aliveCount as a global variable
	lock       sync.RWMutex
	// lockReplication sync.RWMutex
	lockFilename sync.RWMutex
	machineMap   map[int]Machine
	filenameMap  map[string][]string
	// replicationMap map[string]map[string]bool
	replicationMap_1    map[string]map[string]bool
	replicationMap_2    map[string]map[string]bool
	replicationMap_3    map[string]map[string]bool
	IPReplicationMapNum map[int]int
)

type PortInfo struct {
	Port int  `json:"port"`
	Busy bool `json:"busy"`
}

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
					lock.Lock()
					machineMap[machine.ID].Ports[i].Busy = true
					lock.Unlock()
					return port.Port, machine.IP, nil
				}
			}
		}
	}
	return 0, "", errors.New("no available non-busy port found")
}
func removeFilenameFromMachine(targetIP string, filenameToRemove string) {
	for _, machine := range machineMap {
		if machine.IP == targetIP {
			for i, filename := range machine.FileNames {
				if filename == filenameToRemove {
					// Remove the filename from the slice
					machine.FileNames = append(machine.FileNames[:i], machine.FileNames[i+1:]...)
					fmt.Printf("Filename '%s' removed from machine with IP '%s'\n", filenameToRemove, targetIP)
					return
				}
			}
		}
	}
	updateFilenameMap()
	fmt.Printf("No machine found with IP '%s' or filename '%s' not found\n", targetIP, filenameToRemove)
}
func getMachineByIP(machineMap map[int]Machine, ip string) (Machine, bool) {
	for _, machine := range machineMap {
		if machine.IP == ip {
			return machine, true
		}
	}
	return Machine{}, false
}
func findNonBusyPortForFilename(filename string) (int, string, error) {
	ips, ok := filenameMap[filename]
	if !ok {
		return 0, "", errors.New("no available machine found that has the file")
	}
	for _, ip := range ips {
		if machine, ok := getMachineByIP(machineMap, ip); ok {
			if machine.IsAlive {
				for _, port := range machine.Ports {
					if !port.Busy {
						return port.Port, ip, nil
					}
				}
			}
		}
	}
	return 0, "", errors.New("no available non-busy port found that has the file")

}

// ----------  Download  -----------//
func (s *ClientServer) Download(ctx context.Context, req *pb.DownloadRequest) (*pb.DownloadResponse, error) {
	//get fileName from the client
	fileName := req.GetFileName()
	//for debuging:-
	fmt.Println("fileName to be downloaded:", fileName)

	PortNum, DataNodeIp, err := findNonBusyPortForFilename(fileName)
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

	ConfirmClient(clientIp, clientPort)

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

	for _, machine := range machineMap {
		if machine.IP == DataNodeIp {
			lock.Lock()
			machine.FileNames = append(machine.FileNames, fileName)
			machineMap[machine.ID] = machine
			lock.Unlock()
			break
		}
	}

	filenameMap[fileName] = append(filenameMap[fileName],DataNodeIp)
	fmt.Println("inside keeper Done :",fileName,filenameMap[fileName])
	
	/*
	if _, ok := filenameMap[fileName]; ok {
		filenameMap[fileName] = append(filenameMap[fileName],DataNodeIp)
		fmt.Println("inside keeper Done :",filenameMap[fileName])
	} else {
		filenameMap[fileName] = DataNodeIp
	}
	*/
	//updateFilenameMap()

	return &pb.KeeperDoneResponse{}, nil
}

func (s *KeepersServer) KeeperDoneDown(ctx context.Context, req *pb.KeeperDoneDownRequest) (*pb.KeeperDoneDownResponse, error) {
	freePortNum := req.GetPortNum()
	DataNodeIp := req.GetDataNodeIp()

	err := setPortStatus(DataNodeIp, int(freePortNum), false)
	if err != nil {
		fmt.Println("Error:", err)
	}

	//for debuging:-
	// Print the result
	fmt.Println("freePortNum :", freePortNum)
	fmt.Println("DataNodeIp :", DataNodeIp)

	return &pb.KeeperDoneDownResponse{}, nil
}
func ConfirmClient(ip string, port int32) {
	conn, err := grpc.Dial(ip+":"+strconv.Itoa(int(port)), grpc.WithInsecure())
	if err != nil {
		fmt.Println("Failed to connect to client with IP:", ip, ":", err)
		return
	}
	fmt.Println("Conncted to Client %s",ip+":"+strconv.Itoa(int(port)))

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
					lock.Lock()
					machine.Ports[i].Busy = isBusy
					machineMap[machine.ID] = machine
					lock.Unlock()
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
	// fmt.Println("keeperIP :", keeperIP, " has come alive.")

	return &pb.AliveResponse{}, nil
}
func removeElementByValue(arr []string, value string) []string {
	var result []string
	for _, v := range arr {
		if v != value {
			result = append(result, v)
		}
	}
	return result
}
func AliveChecker(numKeepers int) {
	ticker := time.NewTicker(8 * time.Second) // Create a ticker that ticks every 8 seconds
	defer ticker.Stop()                       // Stop the ticker when the function returns

	for range ticker.C {

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

				machineDeadIp := machineMap[i].IP
				for key,_ := range filenameMap {
					for j := int(0); j < len(filenameMap[key]); j++ {
						if machineDeadIp ==  filenameMap[key][j] {
							filenameMap[key] = removeElementByValue(filenameMap[key],machineDeadIp)
							break
						}
					}

					
				}
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

func (s *KeepersServer) ReplicationDone(ctx context.Context, req *pb.ReplicationDoneRequest) (*pb.ReplicationDoneResponse, error) {
	keeperIp := req.GetDataNodeIp()
	portNum := req.GetPortNum()
	fileName := req.GetFileName()

	//for debuging:-
	// Print the result
	fmt.Println("keeperIp :", keeperIp, " finished the replication for file : ", fileName)

	mapNum := IPReplicationMapNum[getIdFromIp(keeperIp)]
	// mapNum := 1
	var replicationMap *map[string]map[string]bool

	// lockReplication.Lock()
	if mapNum == 1 {
		replicationMap = &replicationMap_1
	} else if mapNum == 2 {
		replicationMap = &replicationMap_2
	} else if mapNum == 0 { //NOTE 1,2,0 not 1,2,3
		replicationMap = &replicationMap_3
	} else {
		fmt.Println("error !! undefined mapNum")
	}

	if innerMap, ok := (*replicationMap)[fileName]; ok {
		// Check if the inner map exists for the given fileName
		if _, exists := innerMap[keeperIp]; exists {
			(*replicationMap)[fileName][keeperIp] = true
		}
	}
	// lockReplication.Unlock()

	//mark this "portNum" as an avialble port to the machine with ip = keeperIp
	err := setPortStatus(keeperIp, int(portNum), false)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return &pb.ReplicationDoneResponse{}, nil
}

func replicationFinishfunc(replicationMap map[string]map[string]bool) {
	// replicationMap  = *replicationMap_pass
	for fileName, machine_lists := range replicationMap { //iterate over files

		// for debuging:
		fmt.Printf("fileName: %s, machine_lists: %v\n", fileName, machine_lists)

		for currentIp, done := range machine_lists { //iterate over machines

			if done {
				continue
			}

			if ips, ok := filenameMap[fileName]; ok {
				indexToRemove := -1
				for i, ip := range ips {
					if ip == currentIp {
						indexToRemove = i
						break
					}
				}
				if indexToRemove != -1 {
					lockFilename.Lock()
					filenameMap[fileName] = append(ips[:indexToRemove], ips[indexToRemove+1:]...)
					lockFilename.Unlock()
				}
			}
			removeFilenameFromMachine(currentIp, fileName)
			fmt.Println("Machines:")
			for id, machine := range machineMap {
				fmt.Printf("ID: %d\n", id)
				fmt.Printf("  IP: %s\n", machine.IP)
				fmt.Printf("  FileNames: %s\n", machine.FileNames)
				fmt.Printf("  IsAlive: %t\n", machine.IsAlive)
				fmt.Println("  Ports:")
				for _, port := range machine.Ports {
					fmt.Printf("    %d (%t)\n", port.Port, port.Busy)
				}
				fmt.Println()
			}

		}

	}

	//reset replicationMap //<= no need to it as we are working on a copy
	//replicationMap = make(map[string]map[string]bool)
	// Add other cases if you need to handle other channels
}

func replicationChecker() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	// Get the start time of the application
	startTime := time.Now()

	for range ticker.C {

		currentTime := time.Now()
		elapsedTime := currentTime.Sub(startTime)
		elapsedTimeSeconds := int(elapsedTime.Seconds())

		var replicationMap *map[string]map[string]bool
		ticker_int := int(elapsedTimeSeconds) / 10
		if ticker_int%3 == 1 {
			replicationMap = &replicationMap_1 //by reference not a copy
			if ticker_int != 1 {
				replicationFinishfunc(replicationMap_1) //just take a copy
			}
		} else if ticker_int%3 == 2 {
			replicationMap = &replicationMap_2 //by reference not a copy
			if ticker_int != 2 {
				replicationFinishfunc(replicationMap_2) //just take a copy
			}
		} else if ticker_int%3 == 0 {
			replicationMap = &replicationMap_3 //by reference not a copy
			if ticker_int != 3 {               //NOTE: != 3 not !=0 "samaa"
				replicationFinishfunc(replicationMap_3) //just take a copy
			}
		} else {
			fmt.Println("error !! undefined ticker_int")
			fmt.Println("Elapsed time (seconds):", elapsedTimeSeconds)
		}
		//reset replicationMap
		// lockReplication.Lock()
		(*replicationMap) = make(map[string]map[string]bool)
		// lockReplication.Unlock()

		//for debug:
		// fmt.Println("replicationChecker :", (*replicationMap), "with ticker_int%3 = ", ticker_int%3)
		//fmt.Println("inside replicationChecker")

		for filename, machineIps := range filenameMap {
			// fmt.Println("inside filenameMap loop")

			fmt.Printf("%s: %v\n", filename, machineIps)
			sourceMachineIp := machineIps[0]
			machineIpsLen := len(machineIps)
			//fmt.Println("machineIpsLen = ", machineIpsLen)

			for machineIpsLen < 2 {
				//fmt.Println("inside machineIpsLen < 3 loop")
				destinationMachineIp, destinationMachineId, destMachinePort, err := selectMachineToCopyTo(filename)
				// _, destinationMachineId, _, err := selectMachineToCopyTo(filename)
				// destinationMachineIp := "localhost"
				// destMachinePort := 8003
				// destinationMachineId := 1
				if err != nil {
					fmt.Println("Error: ", err)
				} else {
				
					notifyMachineDataTransfer(sourceMachineIp, destinationMachineIp, destMachinePort, filename)
					machineIpsLen++
					lockFilename.Lock()
					filenameMap[filename] = append(filenameMap[filename], destinationMachineIp)
					lockFilename.Unlock()
					fmt.Println("filenameMap[filename]: ", filenameMap[filename])
					//samaa:
					// lockReplication.Lock()
					if (*replicationMap)[filename] == nil {
						(*replicationMap)[filename] = make(map[string]bool)
					}
					(*replicationMap)[filename][destinationMachineIp] = false
					// lockReplication.Unlock() //NOTE: ReplicationDone service may be called at this time, so we will look it as both it and ReplicationDone write on the 3 of replicationMap(s)

					lock.Lock()
					machine := machineMap[destinationMachineId]
					machine.FileNames = append(machine.FileNames, filename)
					machineMap[destinationMachineId] = machine
					fmt.Println("machineMap[destinationMachineId].FileNames = ", machineMap[destinationMachineId].FileNames)
					lock.Unlock()

					IPReplicationMapNum[destinationMachineId] = ticker_int % 3 //NOTE 1,2,0 not 1,2,3
					fmt.Println("tttttttttttttt")
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
				for i, port := range machine.Ports {
					if !port.Busy {
						lock.Lock()
						machineMap[machine.ID].Ports[i].Busy = true
						lock.Unlock()
						return machine.IP, machine.ID, port.Port, nil
					}
				}
			}
		}
	}
	return "", 0, 0, fmt.Errorf("failed to find machine")
}

func getIdFromIp(Ip string) int {
	for _, machine := range machineMap {
		if machine.IP == Ip {
			return machine.ID
		}
	}
	return -1
}
func notifyMachineDataTransfer(sourceMachineIp string, destinationMachineIp string, destMachinePort int, filename string) error {
	fmt.Println("notifyMachineDataTransfernotifyMachineDataTransfernotifyMachineDataTransfer")

	machineID := getIdFromIp(sourceMachineIp)

	var nonBusyPort int
	machine := machineMap[machineID]
	if machine.IsAlive {
		for i, port := range machine.Ports {
			if !port.Busy {
				lock.Lock()
				machineMap[machineID].Ports[i].Busy = true
				lock.Unlock()
				nonBusyPort = port.Port
				break
			}
		}
	}
	fmt.Println(sourceMachineIp + ":" + strconv.Itoa(int(nonBusyPort)))
	conn, err := grpc.Dial(sourceMachineIp+":"+strconv.Itoa(int(nonBusyPort)), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewFileServiceClient(conn)

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
func updateFilenameMap() {

	for _, machine := range machineMap {
		if machine.IsAlive {
			machineIps := make(map[int]bool)
			for _, filename := range machine.FileNames {
				fmt.Println("filename = " + filename)
				if !machineIps[machine.ID] {
					lockFilename.Lock()
					filenameMap[filename] = append(filenameMap[filename], machine.IP)
					lockFilename.Unlock()
					machineIps[machine.ID] = true
				}
			}
		}
	}
	fmt.Println("filename map:")
	for filename, machines := range filenameMap {
		fmt.Printf("filename: %s\n", filename)
		fmt.Printf("  machines: %v\n", machines)
		fmt.Println()
	}
}

/////////////////////////////// main ///////////////////////////////

func main() {
	filenameMap = make(map[string][]string)
	config, err := loadConfig("config.json")
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}
	machineMap = make(map[int]Machine)
	IPReplicationMapNum = make(map[int]int)

	for _, machine := range config.Machines {
		machineMap[machine.ID] = machine
		machine.FileNames = make([]string, 0)
		// later change it depend for what?
		machine.IsAlive = true
	}

	//NOTE: map[KeyType]ValueType
	// var aliveCount map[int]int
	aliveCount = make(map[int]int)
	numKeepers := int(4)

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
	fmt.Println("Client server started. Listening on port 8081...")

	//-------------  Keeper Done (step5) and Alive (from master)  ------------- //
	lisKeeper, err := net.Listen("tcp", ":8082")
	if err != nil {
		fmt.Println("failed to listen:", err)
		return
	}

	sKeeper := grpc.NewServer()
	pb.RegisterKeepersServiceServer(sKeeper, &KeepersServer{})
	fmt.Println("Keeper server started. Listening on port 8082...")

	go sUp.Serve(lisUp)
	go sKeeper.Serve(lisKeeper)

	go AliveChecker(numKeepers)
	go replicationChecker()

	select {}
}
