package appserver

import "github.com/cmu440/airline/rpc/apprpc"

type AppServer interface {
	GetTicket(args *apprpc.GetTicketArgs, reply *apprpc.GetTicketReply) (error)
	ReserveTicket(args *apprpc.ReserveTicketArgs, reply *apprpc.ReserveTicketReply) (error)
	CancelTicket(args *apprpc.CancelTicketArgs, reply *apprpc.CancelTicketReply) (error)
	//Close()
}