// This file contains constants and arguments used to perform RPCs between
// an application server and the storage servers. DO NOT MODIFY!

package storagerpc

import (
	"github.com/cmu440/airline/util"
	"container/list"
	"time"
)

// Status represents the status of a RPC's reply
type Status int

const (
	OK          Status = iota + 1   // The RPC was a success. 
	KeyNotExist					// The specific key does not exist.
	TicketNotExist				// No left ticket.
	UserNotExist				// No such client in flight reservation
	
	//For server group start up
	NotReady					// The storage serers are still getting ready.
	AlreadySetUp					// The storage servers already set up and receive addtional register request
	
	//For stop recovery
	ServerNotFound 				
)

type Node struct {
	HostPort string					// The host:port address of the storage server node.
	NodeID int					// The ID identifying this storage server node.
}


//a copy struct for stop recovery
type PaxosCopy struct {
	Log map[int] *util.LogEntry    // written to database log
	PendingLog map[int] *util.LogEntry    // commited but not written to database log
	UncommitLog map[int] *util.LogEntry    // uncommited log

	NodeList []Node    // paxos list
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

type RegisterArgs struct {
	ServerInfo Node
}

type RegisterReply struct {
	Status Status
	Servers []Node
}

type GetServersArgs struct {
	// Intentionally left empty
}

type GetServersReply struct {
	Status Status
	Servers []Node
}

type GetArgs struct {
	//TicketReq *util.TicketRequest
	Source string
	Dest string
	DepartTime time.Time
	NumOfTicket int
}

type GetReply struct {
	Status Status
	//Result []*util.TicketReply
	TicketList []util.FlightInfo //a slice of matched flights
	TicketNumList []int //remaining ticket fot the matched flights
}

type CancelArgs struct {
//	TicketReq *util.TicketRequest
	FlightNum string
	Usrid []string
}

type CancelReply struct {
	Status Status
}

type ReserveArgs struct {
	//TicketReq *util.TicketRequest
	FlightNum string
	Usrid []string
}

type ReserveReply struct {
	Status Status
}

//Arguments for stop recovery
type StopArgs struct {
	//left empty
}

type StopReply struct {
	Status Status
}

type ResumeArgs struct {
	//left empty
}

type ResumeReply struct {
	Status Status
}

type GetCopyArgs struct {
	//left empty
}

type GetCopyReply struct {
	Status Status
	Copy PaxosCopy
}

type GetContentArgs struct{

}

type GetContentReply struct {
	Status Status
	Copy PaxosCopy
}