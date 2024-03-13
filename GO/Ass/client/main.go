package main

import (
	"context"
	"fmt"

	pb "Ass/AllServices" // Import the generated package

	"google.golang.org/grpc"
)

func main() {

	//---------  upload file request to master  ---------//
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		fmt.Println("did not connect:", err)
		return
	}
	defer conn.Close()
	c := pb.NewUpdateServiceClient(conn)

	// Call the RPC method
	resp, err := c.Upload(context.Background(), &pb.UpdateRequest{})
	if err != nil {
		fmt.Println("Error calling UploadFile:", err)
		return
	}

	// Print the result
	fmt.Println("PortNum :", resp.GetPortNum())
	fmt.Println("DataNodeIp :", resp.GetDataNodeIp())



	//part2:for uploading: connect to the machine keeper with PortNum and DataNodeIp
	// conn2, err := grpc.Dial("localhost:...", grpc.WithInsecure())
	// if err != nil {
	// 	fmt.Println("did not connect:", err)
	// 	return
	// }
	// defer conn2.Close()
	// c2 := pb.NewUploadFileServiceClient(conn)

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

	// //---------  download file request to master  ---------//
	// connDownload, err := grpc.Dial("localhost:8081", grpc.WithInsecure())
	// if err != nil {
	// 	fmt.Println("did not connect:", err)
	// 	return
	// }
	// defer connDownload.Close()
	// c := pb.NewDownloadServiceClient(connDownload)

	// // Call the RPC method
	// resp, err := c.Download(context.Background(), &pb.DownloadRequest{FileName: "testfile_for_downloading"})
	// if err != nil {
	// 	fmt.Println("Error calling UploadFile:", err)
	// 	return
	// }

	// // Print the result
	// fmt.Println("PortNum :", resp.GetPortNum())
	// fmt.Println("DataNodeIp :", resp.GetDataNodeIp())


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