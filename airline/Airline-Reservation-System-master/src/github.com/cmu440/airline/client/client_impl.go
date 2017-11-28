package client

import (
	"github.com/cmu440/airline/rpc/apprpc"
	"github.com/cmu440/airline/rpc/storagerpc"
	"github.com/cmu440/airline/util"
	//"fmt"
	"strconv"
	"net"
	"net/rpc"
)

type client struct{
	client *rpc.Client
	id string
}

func NewClient(serverHost string, serverPort int, clientid string) (Client, error) {
	cli, err := rpc.DialHTTP("tcp", net.JoinHostPort(serverHost, strconv.Itoa(serverPort)))
	if err != nil {
		return nil, err
	}

	client := new(client)
	client.client = cli
	client.id = clientid

	return client, nil
}

func (c *client) GetTicket(ticketReq *util.TicketRequest) (storagerpc.GetReply, error) {
	args := new(storagerpc.GetArgs)
	var reply storagerpc.GetReply

	args.Source = ticketReq.Source
	args.Dest = ticketReq.Dest
	args.DepartTime = ticketReq.DepartTime
	args.NumOfTicket = ticketReq.NumOfTicket

	err := c.client.Call("StorageServer.Get", args, &reply)
	if err != nil {
		return reply, err
	}

	return reply, nil	
}

//api changed by dateng
func (c *client) ReserveTicket(flightNum string) (apprpc.Status, error) {
	args := new(storagerpc.ReserveArgs)
	var reply storagerpc.ReserveReply

	args.FlightNum = flightNum
	userlist := make([]string, 1)
	userlist[0] = c.id
	args.Usrid = userlist

	err := c.client.Call("StorageServer.Reserve", args, &reply)
	if err != nil {
		return -1, err
	}
	
	// need to fix here
	if reply.Status == storagerpc.OK {
		return apprpc.OK, nil
	} else if reply.Status == storagerpc.KeyNotExist {
		return apprpc.FlightNotExist, nil
	} else if reply.Status == storagerpc.TicketNotExist {
		return apprpc.TicketSoldOut, nil
	} else {
		return apprpc.Fail, nil	
	}
}

//api changed by dateng
func (c *client) CancelTicket(flightNum string) (apprpc.Status, error) {
	args := new(storagerpc.CancelArgs)
	var reply storagerpc.CancelReply

	args.FlightNum = flightNum
	userlist := make([]string, 1)
	userlist[0] = c.id
	args.Usrid = userlist

	err := c.client.Call("StorageServer.Cancel", args, &reply)
	if err != nil {
		return -1, err
	}

	//need to fix here
	if reply.Status == storagerpc.OK {
		return apprpc.OK, nil	
	} else if reply.Status == storagerpc.KeyNotExist {
		return apprpc.FlightNotExist, nil
	} else if reply.Status == storagerpc.TicketNotExist {
		return apprpc.TicketSoldOut, nil
	} else {
		return apprpc.Fail, nil
	}
	
}

func (c *client) Close() error {
	err := c.client.Close()
	if err != nil {
		return err
	}

	return nil
}