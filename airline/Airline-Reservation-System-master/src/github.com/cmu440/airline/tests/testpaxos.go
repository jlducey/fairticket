package main

import (
	"fmt"
	"flag"
	"net/rpc"
	"strconv"
	"time"
	"container/list"
	"github.com/cmu440/airline/util"
	"github.com/cmu440/airline/paxos"
	"github.com/cmu440/airline/rpc/storagerpc"
)

type testFunc struct {
	name string
	f 	 func()
}

func createPaxos(p []paxos.Paxos, flightMap map[string]util.FlightInfo, reserveMap map[string]*list.List, nodeList[]storagerpc.Node, hostport string, me int) {
	p[me] = paxos.NewPaxos(flightMap, reserveMap, nodeList, hostport)
	rpc.HandleHTTP()
}

func handle(paxos.Paxos) {
	rpc.HandleHTTP()
}

func port(base int, port int) string {
	var s string

	port = base + port
	s = "localhost:" + strconv.Itoa(port)
	return s
}

func InitNodeList(npaxos int, baseport int) []storagerpc.Node {
	nodeList := make([]storagerpc.Node, npaxos)
	for i := 0; i < npaxos; i++ {
		nodeList[i].HostPort = port(baseport, i)
		nodeList[i].NodeID = i
	}

	return nodeList
}

func CheckLog(pxa []paxos.Paxos, slotnumber int) int {
	count := 0
	var v string
	for i := 0; i < len(pxa); i++ {
		value, exist := pxa[i].GetLog(slotnumber)
		if exist {
			if count > 0  && v != value {
				return -1
			}
			count++
			v = value
		}
	}

	return count
}

func SingleProposerTest() {
	fmt.Println("Begin SingleProposerTest ...")

	const npaxos = 3
	var pxa []paxos.Paxos = make([]paxos.Paxos, npaxos)

	nodeList := InitNodeList(npaxos, 7000)

	// make new paxos
	//rpc.HandleHTTP()
	for i := 0; i < npaxos; i++ {
		pxa[i] = paxos.CreatePaxos(nil, nil, nodeList, i)
		//go createPaxos(pxa, nil, nil, nodeList, nodeList[i].HostPort, i)
	}
	//go createPaxos(pxa ,nil, nil, nodeList, nodeList[i].HostPort)

	rpc.HandleHTTP()

	fmt.Println("Test: Single proposer ...")

	// submit one value
	pxa[0].Success("1:1:1")

	// wait for a while
	waittime := 1000 * time.Millisecond
	time.Sleep(waittime)

	//status := CheckStatus(pxa, 0, npaxos)
	count := CheckLog(pxa, 1)
	if count == 3 {
		fmt.Println("Single Proposer Passed ...")
	} else if count == -1 {
		fmt.Println("Values do not match! Test Failed!")
	} else {
		fmt.Println("not enough number! Test Failed!")
	}

	return 
}

func MultipleProposerTest() {
	fmt.Println("Begin MultipleProposerTest ...")

	const npaxos = 3
	var pxa []paxos.Paxos = make([]paxos.Paxos, npaxos)

	nodeList := InitNodeList(npaxos, 7010)

	//make new paxos
	for i := 0; i < npaxos; i++ {
		pxa[i] = paxos.CreatePaxos(nil, nil, nodeList, i)
	}

	rpc.HandleHTTP()

	fmt.Println("One slot, One value ...")
	time.Sleep(1 * time.Second)

	for i := 0; i < npaxos; i++ {
		pxa[i].Success("1:1:1")
	}

	// wait for a while
	waittime := 1000 * time.Millisecond
	time.Sleep(waittime)

	count := CheckLog(pxa, 1)

	if count == 3 {
		fmt.Println("One value Passed ...")
	} else if count == -1 {
		fmt.Println("Values do not match! Test Failed!")
		return
	} else {
		fmt.Println("not enough number! Test Failed!")
		return
	}

	fmt.Println("One slot, different values ...")
	for i := 0; i < npaxos; i++ {
		ele := strconv.Itoa(i)
		str := ele + ":" + "1" + ":" + ele
		pxa[i].Success(str)
	}

	time.Sleep(waittime)

	count = CheckLog(pxa, 2)
	if count == 3 {
		fmt.Println("different values Passed ...")
	} else if count == -1 {
		fmt.Println("Values do not match! Test Failed!")
	} else {
		fmt.Println("not enough number! Test Failed!")
	}

	// kill paxos

	return
}

func ManySlotTest() {
	fmt.Println("Begin ManySlotTest ...")

	const npaxos = 3
	var pxa []paxos.Paxos = make([]paxos.Paxos, npaxos)

	nodeList := InitNodeList(npaxos, 7020)

	//make new paxos
	for i := 0; i < npaxos; i++ {
		pxa[i] = paxos.CreatePaxos(nil, nil, nodeList, i)
	}

	rpc.HandleHTTP()

	pxa[0].Success("5:1:5")
	pxa[0].Success("4:1:4")
	pxa[1].Success("3:1:3")
	pxa[1].Success("2:1:2")
	pxa[2].Success("1:1:1")

	waittime := 1000 * time.Millisecond
	time.Sleep(waittime)

	for i := 0; i < 5; i++ {
		count := CheckLog(pxa, i+1)
		if count < 3 {
			fmt.Println("not enough agreement! ManySlotTest Failed!")
			return
		}
	}

	fmt.Println("ManySlotTest Passed ...")

	//delete(pxa)
	return
}

func MessageLostTest() {
	fmt.Println("Begin MessageLostTest ...")

	const npaxos = 3
	var pxa []paxos.Paxos = make([]paxos.Paxos, npaxos)

	nodeList := InitNodeList(npaxos, 7030)

	//make new paxos
	for i := 0; i < npaxos; i++ {
		pxa[i] = paxos.CreatePaxos(nil, nil, nodeList, i)
	}

	rpc.HandleHTTP()

	// set one paxos deaf
	fmt.Println("one message lost")
	pxa[2].SetDeaf()

	// begin propose
	pxa[0].Success("1:1:1")
	waittime := 1000 * time.Millisecond
	time.Sleep(waittime)

	count := CheckLog(pxa, 1)
	if count == 2 {
		fmt.Println("One message lost Passed ...")
	} else if count == 3 {
		fmt.Println("two many agreement! Test Failed")
		return
	} else if count == -1 {
		fmt.Println("Values do not match! Test Failed!")
		return
	} else {
		fmt.Println("not enough number! Test Failed!")
		return
	}

	fmt.Println("two messages lost")
	pxa[1].SetDeaf()

	// begin propose
	pxa[0].Success("1:1:1")
	waittime = 1000 * time.Millisecond
	time.Sleep(waittime)

	count = CheckLog(pxa, 2)
	if count == 0 {
		fmt.Println("Two messages lost Passed ...")
	} else if count > 1 {
		fmt.Println("two many agreement! Test Failed")
		return
	} else if count == -1 {
		fmt.Println("Values do not match! Test Failed!")
		return
	} 

	fmt.Println("MessageLostTest Passed ...")
}

func SingleFailureTest() {
	fmt.Println("Begin SingleFailureTest ...")

	const npaxos = 3
	var pxa []paxos.Paxos = make([]paxos.Paxos, npaxos)

	nodeList := InitNodeList(npaxos, 7040)

	//make new paxos
	for i := 0; i < npaxos; i++ {
		pxa[i] = paxos.CreatePaxos(nil, nil, nodeList, i)
	}

	rpc.HandleHTTP()

	fmt.Println("one node failure")
	pxa[2].Kill()
	pxa[0].Success("0:1:0")

	// set one paxos deaf
	time.Sleep(1 * time.Second)

	// begin propose
	pxa[0].Success("1:1:1")
	waittime := 1000 * time.Millisecond
	time.Sleep(waittime)

	count := CheckLog(pxa, 1)
	if count == 2 {
		fmt.Println("One failure Passed ...")
	} else if count == 3 {
		fmt.Println("two many agreement! Test Failed")
		return
	} else if count == -1 {
		fmt.Println("Values do not match! Test Failed!")
		return
	} else {
		fmt.Println("not enough number! Test Failed!")
		return
	}

	

	fmt.Println("SingleFailureTest Passed ...")

}

func MultipleFailureTest() {
	fmt.Println("Begin MultipleFailureTest ...")

	const npaxos = 3
	var pxa []paxos.Paxos = make([]paxos.Paxos, npaxos)

	nodeList := InitNodeList(npaxos, 7050)

	//make new paxos
	for i := 0; i < npaxos; i++ {
		pxa[i] = paxos.CreatePaxos(nil, nil, nodeList, i)
	}

	rpc.HandleHTTP()

	fmt.Println("two nodes failure")
	pxa[1].Kill()
	pxa[2].Kill()

	// begin propose
	pxa[0].Success("2:1:1")

	waittime := 1000 * time.Millisecond
	time.Sleep(waittime)

	count := CheckLog(pxa, 1)
	if count == 0 {
		fmt.Println("Two failures Passed ...")
	} else if count > 1 {
		fmt.Println("two many agreement! Test Failed")
		return
	} else if count == -1 {
		fmt.Println("Values do not match! Test Failed!")
		return
	} else {
		fmt.Println("not enough number! Test Failed!")
		return
	}

	fmt.Println("MultipleFailureTest Passed ...")
}

func FallbehindTest() {
	fmt.Println("Begin FallbehindTest ...")

	const npaxos = 3
	var pxa []paxos.Paxos = make([]paxos.Paxos, npaxos)

	nodeList := InitNodeList(npaxos, 7060)

	//make new paxos
	for i := 0; i < npaxos; i++ {
		pxa[i] = paxos.CreatePaxos(nil, nil, nodeList, i)
	}

	rpc.HandleHTTP()

	// set one paxos deaf
	fmt.Println("one node fall behind")
	pxa[2].SetDeaf()

	// begin propose
	pxa[0].Success("1:1:1")
	waittime := 1000 * time.Millisecond
	time.Sleep(waittime)

	count := CheckLog(pxa, 1)
	if count == 2 {
		fmt.Println("continue ... ")
	} else if count == 3 {
		fmt.Println("two many agreement! Test Failed")
		return
	} else if count == -1 {
		fmt.Println("Values do not match! Test Failed!")
		return
	} else {
		fmt.Println("not enough number! Test Failed!")
		return
	}

	// first one propose
	pxa[2].UnsetDeaf()
	pxa[0].Success("2:1:1")

	waittime = 1000 * time.Millisecond
	time.Sleep(waittime)

	count = CheckLog(pxa, 2)
	if count == 3 {
		fmt.Println("continue ... ")
	} else if count == -1 {
		fmt.Println("Values do not match! Test Failed!")
		return
	} else {
		//fmt.Println(count)
		fmt.Println("not enough number! Test Failed!")
		return
	}

	pxa[2].Success("3:1:3")

	waittime = 1000 * time.Millisecond
	time.Sleep(waittime)

	count = CheckLog(pxa, 2)
	if count == 3 {
		fmt.Println("FailureTest Passed ... ")
	} else if count == -1 {
		fmt.Println("Values do not match! Test Failed!")
		return
	} else {
		fmt.Println("not enough number! Test Failed!")
		return
	}
}

func main() {

	ptests := []testFunc{
		{"SingleProposerTest", SingleProposerTest},
		{"MultipleProposerTest", MultipleProposerTest},
		{"ManySlotTest", ManySlotTest},
		{"MessageLostTest", MessageLostTest},
		{"SingleFailureTest", SingleFailureTest},
		{"MultipleFailureTest", MultipleFailureTest},
		{"FallbehindTest", FallbehindTest},
	}

	flag.Parse()

	for _, t := range ptests {
		if t.name == flag.Arg(0) {
			t.f()
			break
		}
	}

}