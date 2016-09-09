package main

import (
	"fmt"
	"testing"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkInit(t *testing.T, stub *shim.MockStub, args []string) {
	_, err := stub.MockInit("1", "init", args)
	if err != nil {
		fmt.Println("Init failed", err)
		t.FailNow()
	}
}

func checkState(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("State value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func checkQuery(t *testing.T, stub *shim.MockStub, function string, args []string, value string) {
	bytes, err := stub.MockQuery(function, args)
	if err != nil {
		fmt.Println("Query", function, "failed", err)
		t.FailNow()
	}
	if bytes == nil {
		fmt.Println("Query", function, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("Query value", string(bytes), "was not", value, "as expected")
		t.FailNow()
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, function string, args []string) {
	_, err := stub.MockInvoke("1", function, args)
	if err != nil {
		fmt.Println("Invoke", function, args, "failed", err)
		t.FailNow()
	}
}

func TestNettingChaincode_Init(t *testing.T) {
	log.Info("\n\nNetting Init test start")
	scc := new(NettingChaincode)
	stub := shim.NewMockStub("netting", scc)

	checkInit(t, stub, []string{})
}

func TestNettingChaincode_QueryEmptyStats(t *testing.T) {
	log.Info("\n\nNetting Init test start")
	scc := new(NettingChaincode)
	stub := shim.NewMockStub("netting", scc)

	checkInit(t, stub, []string{})
	checkQuery(t, stub, "Stats", []string{}, "")
}
/*
func TestNettingChaincode_Create3Nodes(t *testing.T) {
	log.Info("\n\nNetting Init3 test start")
	scc := new(NettingChaincode)
	stub := shim.NewMockStub("netting", scc)

	checkInit(t, stub, []string{})
	checkInvoke(t, stub, "AddNode", []string{})
	checkInvoke(t, stub, "AddNode", []string{})
	checkInvoke(t, stub, "AddNode", []string{})
	nodes := []int{0,1,2}
	b, _ := json.Marshal(nodes)
	checkQuery(t, stub, "readNodes", []string{}, string(b))
}

func TestNettingChaincode_addPayment(t *testing.T) {
	log.Info("\n\nNetting Add Payment test start")
	scc := new(NettingChaincode)
	stub := shim.NewMockStub("netting", scc)

	checkInit(t, stub, []string{})
	checkInvoke(t, stub, "initGraph", []string{})
	checkInvoke(t, stub, "addNode", []string{})
	checkInvoke(t, stub, "addNode", []string{})
	checkInvoke(t, stub, "addNode", []string{})
	checkInvoke(t, stub, "addPayment", []string{"1", "2", "3.14"})
}

func TestNettingChaincode_getTrades(t *testing.T) {
	log.Info("\n\nNetting Get Trades test start")
	scc := new(NettingChaincode)
	stub := shim.NewMockStub("netting", scc)

	checkInit(t, stub, []string{})
	checkInvoke(t, stub, "initGraph", []string{})
	checkInvoke(t, stub, "addNode", []string{})
	checkInvoke(t, stub, "addNode", []string{})
	checkInvoke(t, stub, "addNode", []string{})
	checkInvoke(t, stub, "addPayment", []string{"1", "2", "3.14"})
	log.Info("Got: ", string(b))
	checkQuery(t, stub, "getTrades", []string{"1"}, string(b))
}
*/