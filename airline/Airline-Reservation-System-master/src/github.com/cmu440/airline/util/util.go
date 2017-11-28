package util

import (
	"time"
)

type FlightInfo struct {
	FlightNum string
	Source string
	Dest string
	DepartTime time.Time
	ArrivalTime time.Time
	Price float64
	TotalSeat int
}

type ClientInfo struct {
	ID	string	//the key for reservation
	Name string
	Mobile string
}


type TicketRequest struct {
	Source string
	Dest string
	DepartTime time.Time
	NumOfTicket int
}

type TicketReply struct {
	TicketList []*FlightInfo
	TicketNumList []int
}


type LogEntry struct {
	//ReqSeq string	// a failed commit routing to another server?
	Va string //"+1 clientid flightid"
	Na int
	Seqnum int
}
