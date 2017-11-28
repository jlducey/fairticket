package main

import (
	//"github.com/cmu440/airline/rpc/"
	//"github.com/cmu440/airline/rpc/storagerpc"
	"github.com/cmu440/airline/storageserver"
	"os"
	"fmt"
	"strconv"
)
//args "numOfNodes selfport nodeID copyfrom-host:port1 copyfrom-host:port2 ...."
//for example: "3 8080 localhost:8081 localhost:8082"
func main(){
	
	numNodes, err := strconv.Atoi(os.Args[1])
	if err !=nil {
		fmt.Println("Err atoi:", err)
	}
	port, err := strconv.Atoi(os.Args[2])
	if err !=nil {
		fmt.Println("Err atoi:", err)
	}

	nodeID := 0 // or os.Args[3]  //Can we do it this way?
	copyFromAddr := os.Args[3]
	
	_, err = storageserver.RecoverStorageServer(numNodes, port, nodeID, copyFromAddr)
	if err != nil {
		fmt.Println("Err recovering server:", err)
		return
	} else {
		fmt.Println("Recover server port: " + os.Args[2] + "success!")
	}
	
	//keep server alive forever
	select {}
}

