package paxosrpc

import "github.com/cmu440/airline/util"

type Status int 

const (
	OK 		Status = iota + 1
	Cancel
)

type SubmitArgs struct {
	ClientId string
	Number	int //postive represents reserve ticket, negative represents cancel ticket
	FlightInfo util.FlightInfo
	//Operation 	//type int equals Reserve or Cancel
}

type SubmitReply struct {
	Status Status
}

type PrepareArgs struct {
	SlotNumber int
	ProposalNumber int
}

type PrepareReply struct {
	SlotNumber int
	Status Status
	Va	string //previously accepted, uncommitted value
	Na	int//Sequence number of leader who proposed Va
}

type AcceptArgs struct {
	SlotNumber int
	V	string //Value being proposed
	N	int //Sequence number of the proposal of v
}

type AcceptReply struct {
	SlotNumber int
	Status Status 
}

type CommitArgs struct {
	SlotNumber int
	V	string
}

type CommitReply struct {
	//nothing
}