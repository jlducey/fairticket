package client

import (
	"github.com/cmu440/airline/util"
	"github.com/cmu440/airline/rpc/storagerpc"
	"github.com/cmu440/airline/rpc/apprpc"
)

type Client interface {
	// Get Ticket information
	GetTicket(ticketReq *util.TicketRequest) (storagerpc.GetReply, error)

	// Reserve ticket
	ReserveTicket(flightNum string) (apprpc.Status, error)

	// Cancel the booked ticket
	CancelTicket(flightNum string) (apprpc.Status, error)

	// Close connection
	Close() error
}