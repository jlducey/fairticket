A distributed airline reservation system

Author:
Xiaowen Chu: xchu@cs.cmu.edu 
Da Teng: dateng@cs.cmu.edu

FIRST, Please set GOPATH to Airline-Reservation-System dir

eg: export GOPATH=/home/usr/15640hw3/Airline-Reservation-System
After setting GOPATH, you can do following:


==============Test Usage===============================

To run tests, build as follows:
go build github.com/cmu440/airline/tests/

Then you can just run the shell script to check all the tests
sh $GOPATH/src/paxostest.sh



========================Use Runners=====================

To use runners to run the system, we first need to start up servers, for example, as follows:

/****************Server runner**********************/

go run $GOPATH/src/github.com/cmu440/airline/runners/srunner/srunner.go -port=9009 -N=3

go run $GOPATH/src/github.com/cmu440/airline/runners/srunner/srunner.go -master="localhost:9009" -port=8088

go run $GOPATH/src/github.com/cmu440/airline/runners/srunner/srunner.go -master="localhost:9009" -port=8089


Then we can issue whatever command from a client, like the following:

/****************Client runner**********************/

go run $GOPATH/src/github.com/cmu440/airline/runners/crunner/crunner.go -cmd=s -server="localhost:8089" -source="NYC" -dest="LA" -departTime="2014-Jul-01"

go run $GOPATH/src/github.com/cmu440/airline/runners/crunner/crunner.go -cmd=r -server="localhost:8089" -flightNum="1" -clientID="xchu" 

go run $GOPATH/src/github.com/cmu440/airline/runners/crunner/crunner.go -cmd=s -server="localhost:8089" -source="NYC" -dest="LA" -departTime="2014-Jul-01"

go run $GOPATH/src/github.com/cmu440/airline/runners/crunner/crunner.go -cmd=c -server="localhost:8089" -flightNum="1" -clientID="xchu" 

go run $GOPATH/src/github.com/cmu440/airline/runners/crunner/crunner.go -cmd=s -server="localhost:8089" -source="NYC" -dest="LA" -departTime="2014-Jul-01"


=======================Dead Node Recovery================

After we start some servers using srunners, we might want to simulate the stop-recovery situation in which a
we start a new server and copying content from other alive servers to replace the dead node.

To acheive this, we first need to stop currently running servers, for example (supose our server on port 8088 is down), using following command:
go run $GOPATH/src/github.com/cmu440/airline/recovery/stop/stop.go localhost:8089 localhost:9009

Then we can create a new server and copy the data from exsisting servers to it:
go run $GOPATH/src/github.com/cmu440/airline/recovery/newserver/newserver.go 3 8080 localhost:8081 localhost:8082

At last, we will resume the servers we just stopped:
go run $GOPATH/src/github.com/cmu440/airline/recovery/resume/resume.go localhost:8089 localhost:9009


Some References:

Paxos made simple: http://research.microsoft.com/en-us/um/people/lamport/pubs/paxos-simple.pdf

Paxos made Live: http://www.cs.utexas.edu/users/lorenzo/corsi/cs380d/papers/paper2-1.pdf

CMU 15440 Distributed System: http://www.cs.cmu.edu/~dga/15-440/S14/

