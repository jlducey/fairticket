package main

import (
	//"github.com/cmu440/airline/rpc/"
	"github.com/cmu440/airline/rpc/storagerpc"
	"net/rpc"
	"os"
	"fmt"
)
//args "host:port1 host:port2 ...."
//for example: "localhost:8080 localhost:8081"
func main(){
	arg_num := len(os.Args) 
	//args start from 1, the 0th argument is the program name
	for i := 1 ; i < arg_num ;i++ {
		hostport := os.Args[i]
		
		//Call RPC to stop storage server
		client, err := rpc.DialHTTP("tcp", hostport)
		if err != nil {
			fmt.Println("Err dialing to server:", err)
			continue
		}
		args := &storagerpc.ResumeArgs{}
		var reply storagerpc.ResumeReply


		err = client.Call("StorageServer.ResumeServer", args, &reply)
		if err != nil ||  reply.Status != storagerpc.OK {
			fmt.Println("Err stoping storage server:", err)
		} else {
			fmt.Println("Resume server " + hostport + " success!")
		}
	}
}

