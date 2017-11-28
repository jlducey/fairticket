package appserver

import (
	"github.com/cmu440/airline/rpc/apprpc"
	"errors"
)

type appServer struct {

	//Close()
}

func (as *appServer) GetTicket(args *apprpc.GetTicketArgs, reply *apprpc.GetTicketReply) (error) {
	return errors.New("not implemented!")	
}

func (as *appServer) ReserveTicket(args *apprpc.ReserveTicketArgs, reply *apprpc.ReserveTicketReply) (error) {
	return errors.New("not implemented!")
}

func (as *appServer) CancelTicket(args *apprpc.CancelTicketArgs, reply *apprpc.CancelTicketReply) (error) {
	return errors.New("not implemented!")
}