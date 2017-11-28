// This file contains constants and arguments used to perform RPCs between
// an application server and the storage servers. DO NOT MODIFY!

package storagerpc

// STAFF USE ONLY! Students should not use this interface in their code.
type RemoteStorageServer interface {
	RegisterServer(*RegisterArgs, *RegisterReply) error
	GetServers(*GetServersArgs, *GetServersReply) error
	Get(*GetArgs, *GetReply) error
	Cancel(*CancelArgs, *CancelReply) error
	Reserve(*ReserveArgs, *ReserveReply) error
	
	/****Methods for stop recovery***/

	StopServer(*StopArgs, *StopReply) error
	ResumeServer(*ResumeArgs, *ResumeReply) error
	GetCopy(*GetCopyArgs, *GetCopyReply) error
	
	GetContent(*GetContentArgs, *GetContentReply) error
}

type StorageServer struct {
	// Embed all methods into the struct. See the Effective Go section about
	// embedding for more details: golang.org/doc/effective_go.html#embedding
	RemoteStorageServer
}

// Wrap wraps s in a type-safe wrapper struct to ensure that only the desired
// StorageServer methods are exported to receive RPCs.
func Wrap(s RemoteStorageServer) RemoteStorageServer {
	return &StorageServer{s}
}
