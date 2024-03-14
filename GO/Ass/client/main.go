package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	pb "Ass/AllServices" // Import the generated package

	"google.golang.org/grpc"
)

func main() {
	//later:assume that the client connection can be one upload or one download only

	// Read input from user
	fmt.Print("For uploading enter '1' , for download enter '2': ")
	var upORdown int
	fmt.Scanln(&upORdown)

	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		fmt.Println("did not connect:", err)
		return
	}
	defer conn.Close()
	c := pb.NewClientServiceClient(conn)

	if upORdown == 1 {
		//---------  upload file request to master  ---------//
		// Call the RPC method
		resp, err := c.Upload(context.Background(), &pb.UpdateRequest{})
		if err != nil {
			fmt.Println("Error calling UploadFile:", err)
			return
		}
		fmt.Print("Enter file name: ")
		var filename string
		fmt.Scanln(&filename)

		// Print the result
		fmt.Println("PortNum :", resp.GetPortNum())
		fmt.Println("DataNodeIp :", resp.GetDataNodeIp())

		//part2:for uploading: connect to the machine keeper with PortNum and DataNodeIp

		// later: uncomment this line
		conn2, err := grpc.Dial(resp.GetDataNodeIp()+":"+strconv.Itoa(int(resp.GetPortNum())), grpc.WithInsecure())
		// conn2, err := grpc.Dial("localhost:3000", grpc.WithInsecure())
		if err != nil {
			fmt.Println("did not connect:", err)
			return
		}
		defer conn2.Close()
		c2 := pb.NewUploadFileServiceClient(conn2)
		// Open the file to upload
		fileContent, err := ioutil.ReadFile("D:/CUFE24/4th year/second term/Wireless Networks/test.mp4")
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}

		// Send the file content to the server
		_, err = c2.UploadFile(context.Background(), &pb.UploadFileRequest{
			File:       fileContent,
			FileName:   filename,
			ClientIp:   "ClientIp",
			PortNum:    resp.GetPortNum(),
			DataNodeIp: resp.GetDataNodeIp(),
		})
		if err != nil {
			log.Fatalf("Failed to upload file: %v", err)
		}

		fmt.Println("File uploaded successfully")

		// // later: take the file from client :...
		// //NOTE: lazm t2ol ll machine esm el file ely el client 3ays y3mlo save
		// //y3ny e3ml soket tb3t el file name , then eb3 el file nfso chunk by chunck
		// //parse el file 34an ytb3t 1000 bytes for ex
		// //NOTE: you have to mark the end of the file chunks transfaring , you

		// // Call the RPC method
		// resp, err := c.UploadFile(context.Background(), &pb.UploadFileRequest{})
		// if err != nil {
		// 	fmt.Println("Error calling UploadFile:", err)
		// 	return
		// }

		// //transfare file throuh sockets
		// //sokets connection :
		// // conn2, err := ned.tcp.Dial("localhost:...", grpc.WithInsecure())

	} else { //later: change this to "else"
		//---------  download file request to master  ---------//
		// Call the RPC method
		fmt.Print("Enter file name you want to download: ")
		var filename string
		fmt.Scanln(&filename)
		resp, err := c.Download(context.Background(), &pb.DownloadRequest{FileName: filename})
		if err != nil {
			fmt.Println("Error calling UploadFile:", err)
			return
		}

		// Print the result
		fmt.Println("PortNum :", resp.GetPortNum())
		fmt.Println("DataNodeIp :", resp.GetDataNodeIp())

		// Connect to the machine to download the file
		// later: uncomment this line
		conn2, err := grpc.Dial(fmt.Sprintf("%s:%d", resp.GetDataNodeIp(), resp.GetPortNum()), grpc.WithInsecure())
		// conn2, err := grpc.Dial("localhost:3000", grpc.WithInsecure())
		if err != nil {
			fmt.Println("did not connect:", err)
			return
		}
		defer conn2.Close()

		c2 := pb.NewDownloadFileServiceClient(conn2)

		// Call the DownloadFile RPC method
		downloadResp, err := c2.DownloadFile(context.Background(), &pb.DownloadFileRequest{})
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

}

// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"os"
// 	"io/ioutil"
// )

// func main() {
// 	// Read input from user
// 	reader := bufio.NewReader(os.Stdin)
// 	fmt.Print("Enter the file path: ")
// 	filePath, _ := reader.ReadString('\n')

// 	// Remove newline character from the file path
// 	filePath = filePath[:len(filePath)-1]

// 	// Open the file
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		fmt.Println("Error opening file:", err)
// 		return
// 	}
// 	defer file.Close()

// 	// Read the file contents
// 	content, err := ioutil.ReadAll(file)
// 	if err != nil {
// 		fmt.Println("Error reading file:", err)
// 		return
// 	}

// 	// Print the file contents
// 	fmt.Println("File contents:")
// 	fmt.Println(string(content))
// }
