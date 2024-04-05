package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"

	pb "Ass/AllServices" // Import the generated package

	"google.golang.org/grpc"
)

type DoneUpServer struct {
	pb.UnimplementedDoneUpServiceServer
}

func (s *DoneUpServer) DoneUp(ctx context.Context, req *pb.DoneUpRequest) (*pb.DoneUpResponse, error) {
	return &pb.DoneUpResponse{}, nil
}

func main() {

	var clientPort int32 = 8005
	var clientIp string = "localhost"
	var masterPortToClient int32 = 8081
	var masterIp string = "localhost"

	//---------------------------------------------------------------------------//

	conn, err := grpc.Dial(masterIp+":"+strconv.Itoa(int(masterPortToClient)), grpc.WithInsecure())

	if err != nil {
		fmt.Println("did not connect:", err)
		return
	}
	defer conn.Close()
	c := pb.NewClientServiceClient(conn)

	//---------------------------------------------------------------------------//

	//initialize the listener of the "success" request from the master
	lisDone, err := net.Listen("tcp", ":"+strconv.Itoa(int(clientPort)))
	if err != nil {
		fmt.Println("failed to listen in the client port (client side):", err)
		return
	}
	sDone := grpc.NewServer()
	pb.RegisterDoneUpServiceServer(sDone, &DoneUpServer{})
	fmt.Println("Client started. Listening on port = ", clientPort)

	//wait for the requst with success that will be send for the uploaded file from the master
	go func() {
		fmt.Println("before serve")
		if err := sDone.Serve(lisDone); err != nil {
			fmt.Println("failed to serve:", err)
		} else {
			fmt.Println("File uploaded successfully , confirmed from master")
		}
		fmt.Println("after serve")
	}()

	//---------------------------------------------------------------------------//
	for {
		// Read input from user
		fmt.Print("For uploading enter '1' , for download enter '2': ")
		var upORdown int
		fmt.Scanln(&upORdown)

		if upORdown == 1 {
			//---------  upload file request to master  ---------//
			// Call the RPC method
			resp, err := c.Upload(context.Background(), &pb.UpdateRequest{})
			if err != nil {
				fmt.Println("Error calling Upload:", err)
				return
			}
			// Print the result
			fmt.Println("PortNum :", resp.GetPortNum())
			fmt.Println("DataNodeIp :", resp.GetDataNodeIp())

			fmt.Print("Enter file path you want to upload: ")
			var filepath string
			fmt.Scanln(&filepath)

			fmt.Print("Enter file name: ")
			var filename string
			fmt.Scanln(&filename)

			//part2:for uploading: connect to the machine keeper with PortNum and DataNodeIp

			conn2, err := grpc.Dial(resp.GetDataNodeIp()+":"+strconv.Itoa(int(resp.GetPortNum())), grpc.WithInsecure())
			if err != nil {
				fmt.Println("did not connect:", err)
				return
			}
			defer conn2.Close()
			c2 := pb.NewFileServiceClient(conn2)
			fileContent, err := ioutil.ReadFile(filepath + "/" + filename)
			if err != nil {
				log.Fatalf("Failed to read file: %v", err)
			}

			// Send the file content to the server
			fmt.Println(resp.GetPortNum())
			_, err = c2.UploadFile(context.Background(), &pb.UploadFileRequest{
				File:          fileContent,
				FileName:      filename,
				ClientIp:      clientIp,
				PortNum:       resp.GetPortNum(),
				ClientPortNum: clientPort,
			})
			if err != nil {
				log.Fatalf("Failed to upload file: %v", err)
			}

		} else {
			//---------  download file request to master  ---------//
			// Call the RPC method
			fmt.Print("Enter file name you want to download: ")
			var filename string
			fmt.Scanln(&filename)
			resp, err := c.Download(context.Background(), &pb.DownloadRequest{FileName: filename})
			if err != nil {
				fmt.Println("Error calling Download:", err)
				return
			}

			// Print the result
			fmt.Println("PortNum :", resp.GetPortNum())
			fmt.Println("DataNodeIp :", resp.GetDataNodeIp())

			// Connect to the machine to download the file
			conn2, err := grpc.Dial(fmt.Sprintf("%s:%d", resp.GetDataNodeIp(), resp.GetPortNum()), grpc.WithInsecure())

			if err != nil {
				fmt.Println("did not connect:", err)
				return
			}
			defer conn2.Close()

			c2 := pb.NewFileServiceClient(conn2)

			// Call the DownloadFile RPC method
			downloadResp, err := c2.DownloadFile(context.Background(), &pb.DownloadFileRequest{FileName: filename})
			if err != nil {
				fmt.Println("Error calling DownloadFile:", err)
				return
			}

			// Save the downloaded file
			fmt.Print("Enter file path with file name you want to save file in: ")
			fmt.Scanln(&filename)
			err = ioutil.WriteFile(filename, downloadResp.GetFile(), 0644)
			if err != nil {
				fmt.Println("Failed to write file:", err)
				return
			}

			fmt.Println("File downloaded successfully")

		}

		//  conn.Close()
	}

}
