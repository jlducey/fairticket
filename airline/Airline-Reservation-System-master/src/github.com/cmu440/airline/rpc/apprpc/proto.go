package apprpc

import "github.com/cmu440/airline/util"

type Status int

const (
	OK			Status = iota + 1
	FlightNotExist
	TicketSoldOut
	Fail
)

type GetTicketArgs struct{
	TicketReq *util.TicketRequest
}

type ReserveTicketArgs struct {
	ClientId	string
	TicketReq *util.TicketRequest
}

type CancelTicketArgs struct {
	ClientId string
	TicketReq *util.TicketRequest
}

type GetTicketReply struct {
	Status Status
	Tickets []*util.TicketReply
}

type ReserveTicketReply struct {
	Status Status
}

type CancelTicketReply struct {
	Status Status
}