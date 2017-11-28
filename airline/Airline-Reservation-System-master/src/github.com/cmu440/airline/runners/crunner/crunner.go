package main

import (
	"flag"
	"fmt"
	"os"
	"time"
	"strings"
	"strconv"
	"log"
	//"net/rpc"
	//"container/list"

	//"github.com/cmu440/airline/test"
	//"github.com/cmu440/airline/paxos"
	"github.com/cmu440/airline/client"
	"github.com/cmu440/airline/rpc/apprpc"
	"github.com/cmu440/airline/rpc/storagerpc"
	"github.com/cmu440/airline/util"
)

const defaultClientId = "1"

var (
	//must have this argument
	cmd = flag.String("cmd", "", "operation type(s/r/c)")
	serverHostPort = flag.String("server", "", "storage server host port to send the request to")
	
	//arguments for Search
	source 			= flag.String("source", "", "depart location")
	dest 			= flag.String("dest", "", "destination")
	departTime 		= flag.String("departTime", "", "departTime")

	//arguments for Reserve and Cancel
	flightNum      = flag.String("flightNum", "", "flight Num that the request is made to")
	clientID     = flag.String("clientID", "", "ID of the client who reserve the ticket")
	//clientName   = flag.String("clientName", "", "name of the client who reserve the ticket")
	//clientMobile = flag.String("clientMobile", "", "mobile of the client who reserve the ticket")
)

var cmdList = map[string]int{
	//"l": 0,
	"i": 1,
	"s": 3,
	"r": 3,
	"c": 3,
}

type cmdInfo struct {
	cmdline string
	nargs   int
}

func init() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "The crunner program is a testing tool that simulate a client of the system")
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, "Search  Ticket:		-cmd=s -server=\"hostport\" -source=\"depart from\" -dest=\"destination\" -departTime=\"2014-Jan-02\"")
		fmt.Fprintln(os.Stderr, "Reserve Ticket:		-cmd=r -server=\"hostport\" -flightNum=\"which airline\" -clientID=\"id\" ")
		fmt.Fprintln(os.Stderr, "Cancel  Ticket:		-cmd=c -server=\"hostport\" -flightNum=\"which airline\" -clientID=\"id\" ")
	}
}

func IndicateFlight(flights []util.FlightInfo) {
	if len(flights) == 0 {
		fmt.Println("No matched flight is found")
	}
	for i:=0; i < len(flights); i++ {
		v := flights[i]
		fmt.Println("=====================================")
		fmt.Println("Flight number 	  : ", v.FlightNum)
		fmt.Println("From         	  : ", v.Source)
		fmt.Println("To           	  : ", v.Dest)
		fmt.Println("Depart Time  	  : ", v.DepartTime)
		fmt.Println("Arrival Time 	  : ", v.ArrivalTime)
		fmt.Println("Price        	  : ", v.Price)
		fmt.Println("Remaining Seats  : ", v.TotalSeat)
		fmt.Println("=====================================")
	}
}

func IndicateGetReply(reply storagerpc.GetReply) {
	TicketList := &reply.TicketList
	fmt.Println(TicketList)
}

func main() {
	
	flag.Parse()
	cmdmap := make(map[string]cmdInfo)
	for k, v := range cmdList {
		cmdmap[k] = cmdInfo{cmdline: k, nargs: v}
	}

	//Generate client
	if *serverHostPort == "" {
		flag.Usage()
		return	
	}

	hostport := strings.Split(*serverHostPort, ":")
	if len(hostport) != 2 {
		fmt.Println("invalid hostport format")
		flag.Usage()
		return
	}
	hostAddr := hostport[0]
	port, err := strconv.Atoi(hostport[1])
	if err != nil {
		fmt.Println(err)
		flag.Usage()
		return
	}
	
	switch *cmd {
	case "s":
		if *source == "" || *dest == "" || *departTime == "" {
			flag.Usage()
			return
		}

		//Generate client
		client, err := client.NewClient(hostAddr, port, defaultClientId)	
		if err != nil {
			return
		}

		ticketReq := new(util.TicketRequest)
		ticketReq.Source = *source
		ticketReq.Dest = *dest
		
		//Process time format
		const shortForm = "2006-Jan-02"
		var date time.Time
 		date, err = time.Parse(shortForm, *departTime)
 		if err != nil {
 			fmt.Println("Invalid time format: ", err)
 		}
		ticketReq.DepartTime = date

		//Sent get tiket request
		reply, err := client.GetTicket(ticketReq)
		if err != nil {
			fmt.Println("get ticket error!")
		}

		//print all the flight Info
		//IndicateGetReply(reply)
		IndicateFlight(reply.TicketList)

	case "r":
		if *clientID=="" {
			flag.Usage()
			return
		}

		//Generate client
		client, err := client.NewClient(hostAddr, port, *clientID)	
		if err != nil {
			return
		}

		status, err := client.ReserveTicket(*flightNum)
		if err != nil {
			fmt.Println("reserve error!")
		}

		if status == apprpc.OK {
			fmt.Println("reserve success!")
		} else if status == apprpc.FlightNotExist {
			fmt.Println("Flight not exist!")
		} else if status == apprpc.TicketSoldOut {
			fmt.Println("No enough tickets left!")
		} else {
			fmt.Println("Reserve fail")
		}

	case "c":
		if *clientID=="" {
			flag.Usage()
			return
		}

		//Generate client
		client, err := client.NewClient(hostAddr, port, *clientID)	
		if err != nil {
			return
		}

		status, err := client.CancelTicket(*flightNum)
		if err != nil {
			fmt.Println("Cancel error!")
		}

		if status == apprpc.OK {
			fmt.Println("cancel success!")
		} else if status == apprpc.FlightNotExist {
			fmt.Println("Flight not Exist!")
		} else if status == apprpc.TicketSoldOut {
			fmt.Println("Ticket has been Sold Out!")
		} else {
			fmt.Println("Cancel fail!")
		}

	default:
		flag.Usage()
	}

}
