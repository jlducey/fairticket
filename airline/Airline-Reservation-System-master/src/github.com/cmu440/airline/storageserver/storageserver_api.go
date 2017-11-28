package storageserver

import "github.com/cmu440/airline/rpc/storagerpc"
//import "github.com/cmu440/airline/util"

type StorageServer interface {

	// RegisterServer adds a storage server to the ring. It replies with
	// status NotReady if not all nodes in the ring have joined. Once
	// all nodes have joined, it should reply with status OK and a list
	// of all connected nodes in the ring.
	RegisterServer(*storagerpc.RegisterArgs, *storagerpc.RegisterReply) error

	// GetServers retrieves a list of all connected nodes in the ring. It
	// replies with status NotReady if not all nodes in the ring have joined.
	GetServers(*storagerpc.GetServersArgs, *storagerpc.GetServersReply) error

	// Get retrieves the specified flight from the data store and replies with
	// a list of flight that satisfiy requirement. If no flight is available, 
	// it will reply with an empty list
	Get(*storagerpc.GetArgs, *storagerpc.GetReply) error

	// Cancel enable the user to canel a previously booked ticket and replies
	// Ok if the operation is successful. If the user haven't reserved any 
	// tickets, it will reply with user not exist
	Cancel(*storagerpc.CancelArgs, *storagerpc.CancelReply) error

	// Reserve enable the user to reserve a ticket and replies OK if the 
	// reservation is succeed. If the flight user want to exist does not exist,
	// it replies KeyNotExist. If the flight has no more tickets, it replies 
	// with TicketNotExist.
	Reserve(*storagerpc.ReserveArgs, *storagerpc.ReserveReply) error
	
	/****Methods for stop recovery***/
	//stop current server 
	StopServer(*storagerpc.StopArgs, *storagerpc.StopReply) error
	
	//resume server to process request
	ResumeServer(*storagerpc.ResumeArgs, *storagerpc.ResumeReply) error
	
	//get the content of the server to new added server, requirement: serverNum already exist in server ring 
	GetCopy(*storagerpc.GetCopyArgs, *storagerpc.GetCopyReply) error

	GetContent(*storagerpc.GetContentArgs, *storagerpc.GetContentReply) error
}