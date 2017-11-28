package storageserver

import (
	"errors"
	"container/list"
	"github.com/cmu440/airline/rpc/storagerpc"
	"github.com/cmu440/airline/paxos"
	"github.com/cmu440/airline/util"
	"sync"
	"net/http"
	"net/rpc"
	"fmt"
	"net"
	"time"
)

type storageServer struct {
	port   int
	nodeID int

	//var for starting up the storage servers
	numNodes    int
	nodeMap     map[int]storagerpc.Node
	serverReady bool
	
	//functional var of storage servers
	nodeList    []storagerpc.Node
	
	//storage server
	lockMap map[string]*sync.Mutex //lock map for keys in flightMap and reserveMap, a lock for each key
	mutex *sync.Mutex //lock when adding to lockMap
	
	paxos paxos.Paxos
	
	//stores data
	flightMap map[string] util.FlightInfo //flight number <-> flightInfo
	reserveMap map[string] *list.List	//flight number <-> clientList (type: util/client)
}


func NewStorageServer(masterServerHostPort string, numNodes, port int, nodeID int) (StorageServer, error) {
	//hostport string for current storage server
	portStr := fmt.Sprintf("%d", port)
	portStr = "localhost:" + portStr

	paxosPort := fmt.Sprintf("%d", port+10)
	paxosPort = "localhost:" + paxosPort

	//initial storage server struct
	newServer := new(storageServer)
	newServer.port = port
	newServer.nodeID = nodeID
	newServer.numNodes = numNodes
	newServer.nodeMap = make(map[int]storagerpc.Node)
	newServer.serverReady = false
	newServer.nodeList = make([]storagerpc.Node, numNodes)
	
	newServer.flightMap = make(map[string]util.FlightInfo)
	newServer.reserveMap = make(map[string]*list.List)

	newServer.lockMap = make(map[string]*sync.Mutex)
	newServer.mutex = new(sync.Mutex)
	
	
	
	/***This is a master server***/
	if len(masterServerHostPort) == 0 {
		
		//add itself to nodeMap
		newServer.nodeMap[nodeID] = storagerpc.Node{paxosPort, nodeID}
		
		//if only one noede, registerServer will not be call by other slaves
		if numNodes == 1 {
			newServer.nodeList[0] =  storagerpc.Node{paxosPort, nodeID}
		}
		
		//register to receive RPC from the other servers in the system
		rpc.RegisterName("MasterServer", storagerpc.Wrap(newServer))  //for slavves
		rpc.RegisterName("StorageServer", storagerpc.Wrap(newServer)) //for libServer
		rpc.HandleHTTP()                             //can only call for one time!!!

		//listen for incoming connections
		ln, e := net.Listen("tcp", portStr)
		if e != nil {
			fmt.Println("listen error:", e)
			return nil, e
		}
		go http.Serve(ln, nil)

		//wait until all slave join the consitent hashing ring
		for len(newServer.nodeMap) < numNodes {
			time.Sleep(time.Second)
		}

	} else { /****This is a slave server****/
		//connect to server using masterServerHostPort
		client, err := rpc.DialHTTP("tcp", masterServerHostPort)
		if err != nil {
			fmt.Println("Err dialing to master node:", err)
			return nil, err
		}

		nodeInfo := new(storagerpc.Node)
		nodeInfo.HostPort = portStr
		nodeInfo.NodeID = nodeID
		args := &storagerpc.RegisterArgs{*nodeInfo}
		var reply storagerpc.RegisterReply

		//register to master of the servers group
		err = client.Call("MasterServer.RegisterServer", args, &reply)
		if err != nil {
			fmt.Println("Start slave server error:", err)
			return nil, err
		}

		

		//sleep for second and resend request until startup complete
		for reply.Status != storagerpc.OK {
			
			if reply.Status == storagerpc.AlreadySetUp {
				fmt.Println("Exceed servers number limitation. If any server is down, please try to use the recovery mode")
				return nil, errors.New("Exceed servers number limitation")
			}

			time.Sleep(100 * time.Millisecond)
			err := client.Call("MasterServer.RegisterServer", args, &reply)
			if err != nil {
				fmt.Println(" Start slave server error:", err)
				return nil, err
			}
		}

		//update nodeList in struct
		newServer.nodeList = reply.Servers

		//register to RPC
		rpc.RegisterName("StorageServer", storagerpc.Wrap(newServer)) //for libServer
		rpc.HandleHTTP()                             //can only call for one time!!!

		//listen for incoming connections
		ln, e := net.Listen("tcp", portStr)
		if e != nil {
			fmt.Println("listen error:", e)
			return nil, e
		}
		go http.Serve(ln, nil)
	}

	//Set up paxos and data
	newServer.paxos = paxos.NewPaxos(newServer.flightMap, newServer.reserveMap, newServer.nodeList, paxosPort)
	if newServer.paxos == nil {
		return nil, errors.New("error generating Paxos")
	}

	//Insert hard-coded data
	t1 := time.Date(2014, time.April, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2014, time.April, 3, 0, 0, 0, 0, time.UTC)
	newServer.flightMap["1"] = util.FlightInfo{"1", "NYC", "LA", t1, t2, 1000.0, 100}

	
	time.Sleep(1 * time.Second) //wait for other server
	fmt.Println("New Storage Server")
	newServer.serverReady = true
	return newServer, nil
}

//For slave storage servers
func (ss *storageServer) RegisterServer(args *storagerpc.RegisterArgs, reply *storagerpc.RegisterReply) error {
	//Server ring has already set up, so do noting
	if ss.serverReady == true {
		reply.Status = storagerpc.AlreadySetUp
		return nil
	}

	//Add to node map if it's a new slave
	_, exist := ss.nodeMap[args.ServerInfo.NodeID]
	if !exist {
		ss.nodeMap[args.ServerInfo.NodeID] = args.ServerInfo
		fmt.Println("Slave " + args.ServerInfo.HostPort + " registered")
	}

	//RPC reply
	if len(ss.nodeMap) < ss.numNodes {
		reply.Status = storagerpc.NotReady
	} else {
		//all the storage servers are registered in the nodeMap, now copy them to nodeList
		i := 0
		for _, node := range ss.nodeMap {
			ss.nodeList[i] = node
			i = i + 1
		}

		reply.Status = storagerpc.OK
		reply.Servers = ss.nodeList
	}

	return nil
}

// GetServers retrieves a list of all connected nodes in the ring. It
// replies with status NotReady if not all nodes in the ring have joined.
func (ss *storageServer) GetServers(args *storagerpc.GetServersArgs, reply *storagerpc.GetServersReply) error {
	if ss.serverReady == false {
		reply.Status = storagerpc.NotReady
		reply.Servers = nil
	} else {
		reply.Status = storagerpc.OK
		reply.Servers = ss.nodeList
	}

	return nil
}

// Get retrieves the specified flight from the data store and replies with
// a list of flight that satisfiy requirement. If no flight is available, 
// it will reply with an empty list

//get reply by using "source-dest", return a slice of tickets
func (ss *storageServer) Get(args *storagerpc.GetArgs, reply *storagerpc.GetReply) error {
	//fmt.Println("get it!")
	//a temporary list to store found flight 
	resultList := new(list.List)
	
	//if the source, dest and depart time matches reuqest, add the flight to return list
	for _, flightInfo := range(ss.flightMap) {
		if flightInfo.DepartTime.Before(args.DepartTime) {
			fmt.Println("yes!")
		} else {
			fmt.Println("no!")
		}
		if (flightInfo.Source == args.Source) && (flightInfo.Dest == args.Dest) && (flightInfo.DepartTime.Before(args.DepartTime)) {
			rlist, exist := ss.reserveMap[flightInfo.FlightNum]
			if !exist {
				fmt.Println("get it")
				resultList.PushBack(flightInfo)
			} else {
				if (flightInfo.TotalSeat - rlist.Len()) >= args.NumOfTicket {
					resultList.PushBack(flightInfo)	
				}
			}
			
		}	
	}
	
	//copy return list to a slice, to match return type
	if resultList.Len() == 0 {
		reply.Status = storagerpc.KeyNotExist //no matched airline is found
	} else {
		i := 0
		reply.TicketList = make([]util.FlightInfo, resultList.Len())
		reply.TicketNumList = make([]int, resultList.Len())
		for e := resultList.Front(); e != nil; e = e.Next() {
			flightInfo := (e.Value).(util.FlightInfo)
			reply.TicketList[i] = flightInfo

			rlist, exist := ss.reserveMap[flightInfo.FlightNum]
			if !exist {
				reply.TicketNumList[i] = flightInfo.TotalSeat	
			} else {
				reply.TicketNumList[i] = flightInfo.TotalSeat - rlist.Len()
			}
			//reply.TicketNumList[i] = flightInfo.TotalSeat - ss.reserveMap[flightInfo.FlightNum].Len()
		}
		reply.Status = storagerpc.OK
	}
	return nil
}

// Cancel enable the user to canel a previously booked ticket and replies
// Ok if the operation is successful. If the user haven't reserved any 
// tickets, it will reply with user not exist
func (ss *storageServer) Cancel(args *storagerpc.CancelArgs, reply *storagerpc.CancelReply) error {
	
	//see if flight exist
	_, exist := ss.flightMap[args.FlightNum]
	if !exist {
		reply.Status = storagerpc.KeyNotExist
		return nil
	}

	//make sure all the client to be cancel exist in the the reservation list
	reservedClients := ss.reserveMap[args.FlightNum]
	if reservedClients == nil {
		reply.Status = storagerpc.UserNotExist
		return nil
	}

	for i := 0; i < len(args.Usrid); i++ {
		hasClient := false
		id := args.Usrid[i]
		for e := reservedClients.Front(); e != nil; e = e.Next() {
			//val := (e.Value).(string)
			val := (e.Value).(*util.ClientInfo)
			if val.ID == id {
				hasClient = true
				break
			}
		}
		if(!hasClient) {
			reply.Status = storagerpc.UserNotExist //if any client isn't in the list, reply fail
			return nil
		}
	}
	
	//all the clients have reservation, do paxos
	var err error
	reply.Status, err = ss.paxos.Submit(args.Usrid, args.FlightNum, -1 * len(args.Usrid))
	if err != nil {
		return err;
	} 
	
	return nil
}

// Reserve enable the user to reserve a ticket and replies OK if the 
// reservation is succeed. If the flight user want to exist does not exist,
// it replies KeyNotExist. If the flight has no more tickets, it replies 
// with TicketNotExist.
func (ss *storageServer) Reserve(args *storagerpc.ReserveArgs, reply *storagerpc.ReserveReply) error {
	
	flightNum := args.FlightNum
	
	//see if flight exist
	_, exist := ss.flightMap[flightNum]
	if !exist {
		reply.Status = storagerpc.KeyNotExist
		return nil
	}

	//if no enough tickets remains, return fail
	rlist, exist := ss.reserveMap[flightNum]
	if exist {
		if (ss.flightMap[flightNum].TotalSeat - rlist.Len()) < len(args.Usrid) {
		//reply.Status = paxosrpc.Cancel
			reply.Status = storagerpc.TicketNotExist
			return nil
		}
	} else {
		rlist := list.New()
		ss.reserveMap[flightNum] = rlist
	}
	
	
	//if currently have enough tickets, do paxos
	var err error
	reply.Status, err = ss.paxos.Submit(args.Usrid, args.FlightNum, len(args.Usrid))
	if err != nil {
		return err;
	} 
	reply.Status = storagerpc.OK
	return nil
}


/****Methods for stop recovery***/
//function for generating a new server by copying data from "recoverFromAddr"
func RecoverStorageServer (numNodes int, port int, nodeID int, copyFromAddr string) (StorageServer, error) {
	//hostport string for current storage server
	portStr := fmt.Sprintf("%d", port)
	portStr = "localhost:" + portStr
	
	storageServer := new(storageServer)
	storageServer.port = port
	storageServer.nodeID = nodeID
	storageServer.numNodes = numNodes
	
	//connect to a existing server to get copy data
	//client, err := rpc.DialHTTP("tcp", copyFromAddr)			//RPC is blocked
	_, err := rpc.DialHTTP("tcp", copyFromAddr)
	if err != nil {
		fmt.Println("Err dialing to copy node:", err)
		return nil, err
	}

	//args := &storagerpc.GetCopyArgs{}							//RPC is blocked
	var reply storagerpc.GetCopyReply
	//err = client.Call("StorageServer.GetCopy", args, &reply)  //RPC is blocked

	if err != nil {
		fmt.Println("get copy error:", err)
		return nil, err
	}

	//for setup server ring
	storageServer.nodeMap = make(map[int]storagerpc.Node)
	storageServer.serverReady = false
	
	//copy from other server
	storageServer.nodeList = reply.Copy.NodeList
	storageServer.flightMap = reply.Copy.FlightMap
	storageServer.reserveMap = reply.Copy.ReserveMap
	
	
	storageServer.lockMap = make(map[string]*sync.Mutex)
	storageServer.mutex = new(sync.Mutex)
	
	//copy paxos content
	storageServer.paxos = paxos.RecoverPaxos(reply.Copy, portStr)
	
	//register RPC
	rpc.RegisterName("StorageServer", storagerpc.Wrap(storageServer)) //for libServer
	rpc.HandleHTTP()                             //can only call for one time!!!

	//listen for incoming connections
	ln, e := net.Listen("tcp", portStr)
	if e != nil {
		fmt.Println("listen error:", e)
		return nil, e
	}
	go http.Serve(ln, nil)
	
	storageServer.serverReady = true
	return storageServer, nil
}

//get the content of the server to new added server, requirement: serverNum already exist in server ring 
func (ss *storageServer) GetCopy(args *storagerpc.GetCopyArgs, reply *storagerpc.GetCopyReply) error {
	var err error
	reply.Copy, err = ss.paxos.GetCopy()
	if err!= nil {
		fmt.Println("get PAXOS copy err", err)
		return err
	}
	reply.Status = storagerpc.OK

	return nil	
}

//stop current server 
func (ss *storageServer) StopServer(args *storagerpc.StopArgs, reply *storagerpc.StopReply) error {
	//do someting to lock Paxos.Submit
	reply.Status = storagerpc.OK
	return nil
}

//resume server to process request
func (ss *storageServer) ResumeServer(args *storagerpc.ResumeArgs, reply *storagerpc.ResumeReply) error {
	//do someting to unlock Paxos.Submit
	reply.Status = storagerpc.OK
	return nil	
}

//stop current server 
func (ss *storageServer) GetContent(args *storagerpc.GetContentArgs, reply *storagerpc.GetContentReply) error {
	var err error
	//do someting to lock Paxos.Submit
	reply.Copy, err = ss.paxos.GetCopy()
	if err!= nil {
		fmt.Println("get PAXOS copy err", err)
		return err
	}
	reply.Status = storagerpc.OK
	return nil
}
