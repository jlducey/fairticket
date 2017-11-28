package paxos

import (
	"github.com/cmu440/airline/rpc/storagerpc"
)


type Paxos interface {
	//operation is the number of tickets we want to modify, positive for reserve, negative for cancel
	Submit(clientIdList []string, flightNum string, operation int) (storagerpc.Status, error)
	GetCopy() (storagerpc.PaxosCopy, error)
	Success(value string) (bool, error)
	GetLog(slotnumber int) (string, bool)
	Kill()
	SetDeaf()
	UnsetDeaf()
}


