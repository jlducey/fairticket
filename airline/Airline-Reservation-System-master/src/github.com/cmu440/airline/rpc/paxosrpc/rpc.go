// This file provides a type-safe wrapper that should be used to register
// the libstore to receive RPCs from the storage server. DO NOT MODIFY!

package paxosrpc

//import "github.com/cmu440/airline/rpc/paxosrpc"

// STAFF USE ONLY! Students should not use this interface in their code.
type RemotePaxos interface {
	//Prepare(args *paxosrpc.PrepareArgs, reply *paxosrpc.PrepareReply) (error)
	Prepare(args *PrepareArgs, reply *PrepareReply) (error)
	Accept(args *AcceptArgs, reply *AcceptReply) (error)
	Commit(args *CommitArgs, reply *CommitReply) (error)
}

type Paxos struct {
	// Embed all methods into the struct. See the Effective Go section about
	// embedding for more details: golang.org/doc/effective_go.html#embedding
	RemotePaxos
}

// Wrap wraps l in a type-safe wrapper struct to ensure that only the desired
// LeaseCallbacks methods are exported to receive RPCs.
func Wrap(p RemotePaxos) RemotePaxos {
	return &Paxos{p}
}
