#!/bin/bash

PAXOS_TEST=$GOPATH/src/tests
echo "begin SingleProposerTest"
sleep 2
${PAXOS_TEST} "SingleProposerTest"
PAXOS_TEST_PID=$!

sleep 5
kill -9 ${PAXOS_TEST_PID} 2> /dev/null
wait ${STORAGE_SERVER_PID} 2> /dev/null

echo "begin MultipleProposerTest"
sleep 2
${PAXOS_TEST} "MultipleProposerTest"
PAXOS_TEST_PID=$!

sleep 5
kill -9 ${PAXOS_TEST_PID} 2> /dev/null
wait ${STORAGE_SERVER_PID} 2> /dev/null

echo "begin ManySlotTest"
sleep 2
${PAXOS_TEST} "ManySlotTest"
PAXOS_TEST_PID=$!

sleep 5
kill -9 ${PAXOS_TEST_PID} 2> /dev/null
wait ${STORAGE_SERVER_PID} 2> /dev/null

echo "begin MessageLostTest"
sleep 2
${PAXOS_TEST} "MessageLostTest"
PAXOS_TEST_PID=$!

sleep 5
kill -9 ${PAXOS_TEST_PID} 2> /dev/null
wait ${STORAGE_SERVER_PID} 2> /dev/null

echo "begin SingleFailureTest"
sleep 2
${PAXOS_TEST} "SingleFailureTest"
PAXOS_TEST_PID=$!

sleep 5
kill -9 ${PAXOS_TEST_PID} 2> /dev/null
wait ${STORAGE_SERVER_PID} 2> /dev/null

echo "begin MultipleFailureTest"
sleep 2
${PAXOS_TEST} "MultipleFailureTest"
PAXOS_TEST_PID=$!

sleep 5
kill -9 ${PAXOS_TEST_PID} 2> /dev/null
wait ${STORAGE_SERVER_PID} 2> /dev/null

echo "begin FallbehindTest"
sleep 2
${PAXOS_TEST} "FallbehindTest"
PAXOS_TEST_PID=$!

sleep 5
kill -9 ${PAXOS_TEST_PID} 2> /dev/null
wait ${STORAGE_SERVER_PID} 2> /dev/null