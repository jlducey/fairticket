package paxos

import (
	"github.com/cmu440/airline/rpc/paxosrpc"
	"github.com/cmu440/airline/rpc/storagerpc"
	"github.com/cmu440/airline/util"

	"strings"
	"errors"
	"strconv"
	"sync"
	"fmt"
	"net"
	//"math/rand"
	"container/list"
	"net/rpc"
	"net/http"
)

type paxos struct {
	log map[int] *util.LogEntry    // written to database log
	pendingLog map[int] *util.LogEntry    // commited but not written to database log
	uncommitLog map[int] *util.LogEntry    // uncommited log

	nodeList []storagerpc.Node    // paxos list
	rpcClientMap map[string]*rpc.Client // hostport <-> *rpc.Client, store rpc client for nodelist
	
	nextSlot int    // next slot number
	majorityNum int

	lock *sync.Mutex

	flightMap map[string] util.FlightInfo   // flight number to remaining ticket
	reserveMap map[string] *list.List      //flight number <-> clientList (type: util/client)
	
	ln net.Listener
	dead bool
	deaf bool
  	unreliable bool

  	me int
  	test bool
}



/*
//A copy struct for stop recovery
type PaxosCopy struct {
	Log map[int] *util.LogEntry    // written to database log
	PendingLog map[int] *util.LogEntry    // commited but not written to database log
	UncommitLog map[int] *util.LogEntry    // uncommited log

	NodeList []storagerpc.Node    // paxos list
	//RpcClientMap map[string]*rpc.Client // hostport <-> *rpc.Client, store rpc client for nodelist
	
	NextSlot int    // next slot number
	MajorityNum int

	//lock *sync.Mutex

	FlightMap map[string] util.FlightInfo   // flight number to remaining ticket
	ReserveMap map[string] *list.List      //flight number <-> clientList (type: util/client)
	
//	Ln net.Listener
	Dead bool
	Deaf bool
  	Unreliable bool
}
*/

func CreatePaxos(flightMap map[string]util.FlightInfo, reserveMap map[string]*list.List, nodeList[]storagerpc.Node, index int) Paxos {
	//fmt.Println("begin Make")
	paxos := new(paxos)

	log := make(map[int]*util.LogEntry)
	pendingLog := make(map[int]*util.LogEntry)
	uncommitLog := make(map[int]*util.LogEntry)

	paxos.log = log
	paxos.pendingLog = pendingLog
	paxos.uncommitLog = uncommitLog
	
	paxos.nodeList = nodeList
	paxos.rpcClientMap = make(map[string]*rpc.Client)
	
	paxos.nextSlot = 1
	paxos.majorityNum = len(paxos.nodeList) / 2 + 1

	if flightMap == nil {
		paxos.flightMap = make(map[string]util.FlightInfo)
	} else {
		paxos.flightMap = flightMap
	}

	if reserveMap == nil {
		paxos.reserveMap = make(map[string]*list.List)
	} else {
		paxos.reserveMap = reserveMap
	}
	
	paxos.lock = new(sync.Mutex)
	//register to receive RPC from the other servers in the system
	paxosName := getPaxosName(index)
	rpc.RegisterName(paxosName, paxosrpc.Wrap(paxos))  //for slavves
	//rpc.HandleHTTP() 

	//listen for incoming connections
	hostport := nodeList[index].HostPort
	ln, e := net.Listen("tcp", hostport)
	if e != nil {
		fmt.Println("listen error:", e)
	}
	//paxos.ln, _ = net.Listen("tcp", hostport)
	go http.Serve(ln, nil)
	
	//variable for testing
	paxos.dead = false
	paxos.deaf = false
	paxos.unreliable = false
	paxos.test = true
	paxos.me = index

	paxos.ln = ln

	return paxos
}

//function for stop recovery
func RecoverPaxos(copy storagerpc.PaxosCopy, hostport string) Paxos {
	paxos := new(paxos)
	
	paxos.log = copy.Log
	paxos.pendingLog = copy.PendingLog
	paxos.uncommitLog = copy.UncommitLog
	
	paxos.nodeList = copy.NodeList
	paxos.rpcClientMap = make(map[string]*rpc.Client)
	
	paxos.nextSlot = copy.NextSlot
	paxos.majorityNum = copy.MajorityNum
	
	paxos.lock = new(sync.Mutex)
	
	paxos.flightMap = copy.FlightMap
	paxos.reserveMap = copy.ReserveMap
	
	//register to receive RPC from the other servers in the system
	rpc.RegisterName("Paxos", paxosrpc.Wrap(paxos))  //for slavves
	//rpc.HandleHTTP() 

	//FOLLOWING CONFLICT!!!
	//listen for incoming connections
	/*ln, e := net.Listen("tcp", hostport)
	if e != nil {
		fmt.Println("listen error:", e)
	}
	go http.Serve(ln, nil)
	paxos.ln = ln
	*/
	paxos.dead = false
	paxos.deaf = false
	paxos.unreliable = false
	
	return paxos
}

func NewPaxos(flightMap map[string]util.FlightInfo, reserveMap map[string]*list.List, nodeList[]storagerpc.Node, hostport string) Paxos {
			
	paxos := new(paxos)

	log := make(map[int]*util.LogEntry)
	pendingLog := make(map[int]*util.LogEntry)
	uncommitLog := make(map[int]*util.LogEntry)

	paxos.log = log
	paxos.pendingLog = pendingLog
	paxos.uncommitLog = uncommitLog
	paxos.test = false
	
	paxos.nodeList = nodeList
	paxos.rpcClientMap = make(map[string]*rpc.Client)
	
	paxos.nextSlot = 1
	paxos.majorityNum = len(paxos.nodeList) / 2 + 1

	if flightMap == nil {
		paxos.flightMap = make(map[string]util.FlightInfo)
	} else {
		paxos.flightMap = flightMap
	}

	if reserveMap == nil {
		paxos.reserveMap = make(map[string]*list.List)
	} else {
		paxos.reserveMap = reserveMap
	}
	
	paxos.lock = new(sync.Mutex)
	//register to receive RPC from the other servers in the system
	rpc.RegisterName("Paxos", paxosrpc.Wrap(paxos))  //for slavves
	//rpc.HandleHTTP() 

	//listen for incoming connections
	ln, e := net.Listen("tcp", hostport)
	if e != nil {
		fmt.Println("listen error:", e)
		return nil
	}
	go http.Serve(ln, nil)
	paxos.ln = ln
	
	//variable for testing
	paxos.dead = false
	paxos.deaf = false
	paxos.unreliable = false

	return paxos
}

func (px *paxos) GetCopy() (storagerpc.PaxosCopy, error) {
	copy := new(storagerpc.PaxosCopy)
	
	copy.Log = px.log
	copy.PendingLog = px.pendingLog
	copy.UncommitLog = px.uncommitLog

	//copy.NodeList = px.nodeList
	//RpcClientMap map[string]*rpc.Client // hostport <-> *rpc.Client, store rpc client for nodelist
	
	copy.NextSlot = px.nextSlot
	copy.MajorityNum = px.majorityNum

	//lock *sync.Mutex

	copy.FlightMap = px.flightMap
	copy.ReserveMap = px.reserveMap
	
//	Ln net.Listener
	copy.Dead = px.dead
	copy.Deaf = px.deaf
  	copy.Unreliable = px.unreliable
	
	return *copy, nil
}

//Before calling this function:
//if command is "Cancel" storage server make sure all clients exist in the resveMap
func (px *paxos) Submit(clientInfoList []string, flightNum string, operation int) (storagerpc.Status, error) {
	
	number := operation //postive for reserve, negative for cancel
	//HERE IS THE FORMAT for PAXOS COMMAND (value) -- "flight:operation:clintID1:clientID2..."
	value := flightNum + ":" + strconv.Itoa(operation)
	for i := 0; i < len(clientInfoList); i++ {
		value = value + ":" + clientInfoList[i]
	}
	
	//Lock when begin paxos
	px.lock.Lock()
	
	//check whehter there is enough tickets
	if number > 0 && (px.flightMap[flightNum].TotalSeat - px.reserveMap[flightNum].Len()) < number {
		//reply.Status = paxosrpc.Cancel
		return storagerpc.TicketNotExist, nil
	}
	
	//do paxos until success
	for ok, err := px.Success(value); !ok; {
		if err != nil {
			return storagerpc.OK, err
		}
		
		if number > 0 && (px.flightMap[flightNum].TotalSeat - px.reserveMap[flightNum].Len()) < number {
			//reply.Status = paxosrpc.Cancel
			return storagerpc.TicketNotExist, nil
		}
	}

	px.lock.Unlock()
	return storagerpc.OK, nil
}

func (px *paxos) Success(value string) (bool, error) {
	ok, err := px.SendPrepare()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	ok, err = px.SendAccept(value)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	commitValue, err := px.SendCommit()
	if err != nil {
		return false, err
	}

	// string compare???
	if commitValue != value {
		return false, nil
	}

	fmt.Println("Success one paxos slot!")
	return true, nil
}

func (px *paxos) SendPrepare() (bool, error) {
	args := new(paxosrpc.PrepareArgs)
	var reply paxosrpc.PrepareReply

	prepareOKNum := 0 
	maxNa := 0
	var maxVa string
	
	//generate rpc args
	args.SlotNumber = px.nextSlot

	_, exist := px.uncommitLog[args.SlotNumber]
	if !exist {
		templog := new(util.LogEntry)
		templog.Seqnum = 0
		px.uncommitLog[args.SlotNumber] = templog
	}
	px.uncommitLog[args.SlotNumber].Seqnum = px.uncommitLog[args.SlotNumber].Seqnum + 1
	args.ProposalNumber = px.uncommitLog[args.SlotNumber].Seqnum

	// call rpc 
	//fmt.Println("nodeList: !!!!!", px.nodeList)
	for i :=0; i < len(px.nodeList); i++ {
		//get rpc client for the node
		hostport := px.nodeList[i].HostPort
		//fmt.Println(hostport)
		client := px.rpcClientMap[hostport]
		if client == nil {
			var err error
			//fmt.Println("dial hostport is:", hostport)
			tempstr := strings.Split(hostport, ":")
			if len(tempstr) == 1 {
				client, err = rpc.DialHTTP("tcp", "localhost"+hostport)
			} else {
				client, err = rpc.DialHTTP("tcp", hostport)
			}
			
			//client, err = rpc.DialHTTP("tcp", "localhost"+hostport)
			if err != nil {
				fmt.Println("Err dialing to master node:", err)
				continue
			}
			px.rpcClientMap[hostport] = client
		}

		//call rpc

		//fmt.Println(client)

		if px.test == true {
			paxosName := getPaxosName(i)
			err := client.Call(paxosName+".Prepare", args, &reply)
			if err != nil {
				fmt.Println("prepare paxos error: ",err)
			}
		} else {

			err := client.Call("Paxos.Prepare", args, &reply)
			if err != nil {
				fmt.Println("prepare paxos error: ",err)
			}
		}


		//fmt.Println(reply)
		// count number of OK reply
		if reply.Status == paxosrpc.OK {
			prepareOKNum++

			// get maxvalue from replys
			if reply.Na > maxNa {
				maxNa = reply.Na
				maxVa = reply.Va
			}
		}
	}

	// majority agrees
	fmt.Println("--------- prepareOKNum is :", prepareOKNum)
	if prepareOKNum >= px.majorityNum {
		
		 //should we do this until accpet phase?
		templog, exist := px.uncommitLog[args.SlotNumber]
		if !exist {
			templog = new(util.LogEntry)
		}
		templog.Na = maxNa
		templog.Va = maxVa 
		
		px.uncommitLog[args.SlotNumber] = templog
		//should we do this until accpet phase?
		fmt.Println("finished SendPrepare with Success!")
		return true, nil
	} else {
		// TODO: if we need maxNa to update???
		fmt.Println("Quit SendPrepare!")
		return false, nil
	}
}

// Acceptor function
func (px *paxos) Prepare(args *paxosrpc.PrepareArgs, reply *paxosrpc.PrepareReply) error {
	//fmt.Println("prepare me is :", px.me)

	if px.deaf == true {
		reply.Status = paxosrpc.Cancel
		return nil
	}

	slotNumber := args.SlotNumber
	
	//get util.LogEntry for the requested slot
	var templog *util.LogEntry
	if slotNumber < px.nextSlot {
		templog = px.log[slotNumber]
		//fmt.Println(templog)
	} else if slotNumber >= px.nextSlot {
		//fmt.Println("slotNumber: ", slotNumber)
		templog, exist := px.uncommitLog[slotNumber]
		if !exist {
			templog = new(util.LogEntry)

			templog.Va = ""
			templog.Na = 0
			templog.Seqnum = 0

			px.uncommitLog[slotNumber] = templog
		}
	}
	
	templog = px.uncommitLog[slotNumber]
	if templog == nil {
		templog = px.log[slotNumber]
	}
	//decide how to respond according to ProposalNumber and current Seqnum
	//fmt.Println(args.ProposalNumber)
	if args.ProposalNumber >= templog.Seqnum {
		templog.Seqnum = args.ProposalNumber //add by dateng, update seqnum
		
		reply.SlotNumber = slotNumber
		reply.Status = paxosrpc.OK
		reply.Va = templog.Va 
		reply.Na = templog.Na 
	} else {
		reply.SlotNumber = slotNumber
		reply.Status = paxosrpc.Cancel
	}

	return nil
}

// proposer function
func (px *paxos) SendAccept(value string) (bool, error) {
	args := new(paxosrpc.AcceptArgs)
	var reply paxosrpc.AcceptReply

	templog, exist := px.uncommitLog[px.nextSlot]
	if !exist {
		//fmt.Println("Pending log not exist in SendAccept!")
		return false, errors.New("Pending log not exist in SendAccept!")
	}

	// update the value to be the last maximum
	args.SlotNumber = px.nextSlot
	args.V = templog.Va
	//args.N = templog.Na //change by dateng, should be following??
	args.N = templog.Seqnum

	// first time
	if templog.Na == 0 {
		args.V = value
		args.N = px.uncommitLog[args.SlotNumber].Seqnum
	}

	acceptOKNum := 0
	// call rpc
	for i :=0; i < len(px.nodeList); i++ {
		//get rpc client for the node
		hostport := px.nodeList[i].HostPort
		//fmt.Println(hostport)
		client := px.rpcClientMap[hostport]
		if client == nil {
			var err error
			client, err = rpc.DialHTTP("tcp", hostport)
			if err != nil {
				fmt.Println("Err dialing to master node:", err)
				continue
			}
		}


		if px.test == true {
			paxosName := getPaxosName(i)
			err := client.Call(paxosName+".Accept", args, &reply)
			if err != nil {
				fmt.Println("paxos accept error: ",err)
			}
		} else {
			err := client.Call("Paxos.Accept", args, &reply)
			if err != nil {
				fmt.Println("paxos accept error: ",err)
			}

		}

		if reply.SlotNumber != px.nextSlot {
			//fmt.Println("SendAccept get wrong SlotNumber!")
			return false, errors.New("wrong SlotNumber!")
		}

		if reply.Status == paxosrpc.OK {
			acceptOKNum++
		}
	}

	fmt.Println("---------- acceptOKNum is :", acceptOKNum)
	if acceptOKNum >= px.majorityNum {
		//TODO: check if accept value equals to proposed value
		return true, nil
	} else {
		return false, nil
	}

}

// Acceptor function
func (px *paxos) Accept(args *paxosrpc.AcceptArgs, reply *paxosrpc.AcceptReply) error {
	SlotNumber := args.SlotNumber

	if px.deaf == true {
		reply.Status = paxosrpc.Cancel
		return nil
	}

	var templog *util.LogEntry
	// commited slot
	if SlotNumber < px.nextSlot {
		templog = px.log[SlotNumber]
	} else {
		// current or future slot
		templog, exist := px.uncommitLog[SlotNumber]
		if !exist {
			// did not receive prepare request
			templog = new(util.LogEntry)
			templog.Va = ""
			templog.Na = 0
			templog.Seqnum = 0

			px.uncommitLog[SlotNumber] = templog
		}
	}

	templog = px.uncommitLog[SlotNumber]
	if templog == nil {
		templog = px.log[SlotNumber]
	}

	// agree to accept
	if args.N >= templog.Seqnum {
		templog.Va = args.V
		
		templog.Na = args.N //added by dateng, we need to set Na here??????
		
		reply.SlotNumber = SlotNumber
		reply.Status = paxosrpc.OK
	} else {
		reply.SlotNumber = SlotNumber
		reply.Status = paxosrpc.Cancel
	}

	return nil
}

func (px *paxos) SendCommit() (string, error) {
	args := new(paxosrpc.CommitArgs)
	var reply paxosrpc.CommitReply

	args.SlotNumber = px.nextSlot

	templog, exist := px.uncommitLog[args.SlotNumber]
	if !exist {
		//fmt.Println("Send Commit error! No such log!")
		return "", errors.New("Send Commit error!")
	}

	args.V = templog.Va

	// call rpc
	for i :=0; i < len(px.nodeList); i++ {
		//get rpc client for the node
		hostport := px.nodeList[i].HostPort
		client := px.rpcClientMap[hostport]
		if client == nil {
			var err error
			client, err = rpc.DialHTTP("tcp", hostport)
			if err != nil {
				fmt.Println("Err dialing to master node:", err)
				continue
			}
		}
		
		//call rpc
		//fmt.Println(client)

		if px.test == true {
			paxosName := getPaxosName(i)
			err := client.Call(paxosName+".Commit", args, &reply)

			if err != nil {
				fmt.Println("paxos Commit error: ",err)			
			}
		} else {
			err := client.Call("Paxos.Commit", args, &reply)

			if err != nil {
				fmt.Println("paxos Commit error: ",err)			
			}
		}	
	}

	return args.V, nil
}

// Acceptor function
func (px *paxos) Commit(args *paxosrpc.CommitArgs, reply *paxosrpc.CommitReply) error {
	if px.deaf == true {
		return nil
	}

	slotNumber := args.SlotNumber
	V := args.V
	//fmt.Println("me is :",px.me)

	if slotNumber < px.nextSlot {
		//fmt.Println("commit slotNumber: ", slotNumber)
		//fmt.Println("commit nextSlot: ", px.nextSlot)
		//fmt.Println("wrong number!!!!")

		px.nextSlot ++
		return nil
	}

	if slotNumber == px.nextSlot {
		templog := px.uncommitLog[slotNumber]
		templog.Va = V
		px.commitDatabase(V)
		px.log[slotNumber] = templog
		delete(px.uncommitLog, slotNumber)
		fmt.Println("commited to log!!!!")

		px.nextSlot ++

		templog = px.pendingLog[px.nextSlot]

		for templog != nil {
			px.commitDatabase(V)
			px.log[px.nextSlot] = templog
			delete(px.pendingLog, px.nextSlot)
			px.nextSlot = px.nextSlot + 1
			templog = px.pendingLog[px.nextSlot]
		}
	} else if slotNumber > px.nextSlot {
		fmt.Println("commit to pendingLog!!!!")
		templog := px.uncommitLog[slotNumber]
		templog.Va = V
		px.pendingLog[slotNumber] = templog
		delete(px.uncommitLog, slotNumber)
	}

	return nil
}

func (px *paxos) commitDatabase(value string) {
	tempStrings := strings.Split(value, ":")
	flightNum := tempStrings[0] 					
	operation, err := strconv.Atoi(tempStrings[1])	
	//clientid := tempStrings[i] //remaining are all clientId
	if err!= nil {
		fmt.Println("string to int err")
	}
	
	//get client list of the flight
	clientList, exist := px.reserveMap[flightNum]
	if !exist {
		clientList = new(list.List)
	}

	// modify client list
	if operation > 0 {
		for  i := 2; (i - 2) < operation; i++ {
			//fmt.Println(tempStrings)
			clientInfo := strings.Split(tempStrings[i], ",") //clientID,clientName,clientMobile
			
			client := new(util.ClientInfo)
			if(len(clientInfo) == 3) {
				client.ID = clientInfo[0]
				client.Name =  clientInfo[1]
				client.Mobile =  clientInfo[2]
			} else {
				client.ID = clientInfo[0]
				client.Name =  "name"
				client.Mobile =  "mobile"
			}
			
			clientList.PushBack(client)	

			flightinfo := px.flightMap[flightNum]
			flightinfo.TotalSeat--
			px.flightMap[flightNum] = flightinfo
		}
	} else {
		number := -1 * operation;
		for  i := 2; (i - 2) < number; i++ {
			
			clientInfo := strings.Split(tempStrings[i], ",") //clientID,clientName,clientMobile
			clientid := clientInfo[0] //no matter there is 3 attribute or only 1 attribute, the first one is always clientid
			
			for e := clientList.Front(); e != nil; e = e.Next() {
				val := (e.Value).(*util.ClientInfo)
				if val.ID == clientid {
					clientList.Remove(e)
					break
				}
			}

			flightinfo := px.flightMap[flightNum]
			flightinfo.TotalSeat++
			px.flightMap[flightNum] = flightinfo
		}
	}

	px.reserveMap[flightNum] = clientList
}

func getPaxosName(me int) string {
	return "Paxos"+strconv.Itoa(me)
}

func (px *paxos) GetLog(slotnumber int) (string, bool) {
	logEntry, exist := px.log[slotnumber]
	if exist {
		return logEntry.Va, true
	} else {
		logEntry, exist = px.pendingLog[slotnumber]
		if exist {
			return logEntry.Va, true
		}
	}

	return "", false
}

// Shut down paxos
func (px *paxos) Kill() {
  if px.ln != nil {
    px.ln.Close()
  }
}

func (px *paxos) SetDeaf() {
	px.deaf = true
}

func (px *paxos) UnsetDeaf() {
	px.deaf = false
}