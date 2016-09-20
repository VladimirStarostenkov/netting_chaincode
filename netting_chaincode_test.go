package main

import (
	"encoding/json"
	"fmt"
	"github.com/VladimirStarostenkov/netting"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"testing"
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
	log.Info("\n\nInit test")
	scc := new(Chaincode)
	stub := shim.NewMockStub("netting", scc)
	// calls
	checkInit(t, stub, []string{"testMode"})
}

func TestNettingChaincode_QueryEmptyStats(t *testing.T) {
	log.Info("\n\nQuery empty stats test")
	scc := new(Chaincode)
	stub := shim.NewMockStub("netting", scc)
	referenceStats := netting.NettingTableStats{
		NumberOfCounterParties: 0,
		NumberOfClaims:         0,
		MetricL1:               -1.0,
		MetricL2:               -1.0,
		SumH:                   0.0,
	}
	referenceBytes, _ := json.Marshal(referenceStats)
	// calls
	checkInit(t, stub, []string{"testMode"})
	checkQuery(t, stub, "Stats", []string{}, string(referenceBytes))
}

func TestNettingChaincode_Query3NodesStats(t *testing.T) {
	log.Info("\n\nQuery 3 nodes stats test")
	scc := new(Chaincode)
	stub := shim.NewMockStub("netting", scc)
	referenceStats := netting.NettingTableStats{
		NumberOfCounterParties: 3,
		NumberOfClaims:         0,
		MetricL1:               0.0,
		MetricL2:               0.0,
		SumH:                   0.0,
	}
	referenceBytes, _ := json.Marshal(referenceStats)
	//calls
	checkInit(t, stub, []string{"testMode"})
	checkInvoke(t, stub, "AddCounterParty", []string{})
	checkInvoke(t, stub, "AddCounterParty", []string{})
	checkInvoke(t, stub, "AddCounterParty", []string{})
	checkQuery(t, stub, "Stats", []string{}, string(referenceBytes))
}

func TestNettingChaincode_Query3NodesWithClaim(t *testing.T) {
	log.Info("\n\nQuery 3 nodes with claim stats test")
	scc := new(Chaincode)
	stub := shim.NewMockStub("netting", scc)
	referenceString := "[{\"f\":1,\"t\":2,\"v\":3.14}]"
	//calls
	checkInit(t, stub, []string{"testMode"})
	checkInvoke(t, stub, "AddCounterParty", []string{})
	checkInvoke(t, stub, "AddCounterParty", []string{})
	checkInvoke(t, stub, "AddCounterParty", []string{})
	checkInvoke(t, stub, "AddClaim", []string{"1", "2", "3.14"})
	checkQuery(t, stub, "Claims", []string{"1"}, referenceString)
}

func TestNettingChaincode_Query3NodesWith2Claims(t *testing.T) {
	log.Info("\n\nQuery 3 nodes with claims stats test")
	scc := new(Chaincode)
	stub := shim.NewMockStub("netting", scc)
	referenceString := "[{\"f\":1,\"t\":2,\"v\":6.28}]"
	//calls
	checkInit(t, stub, []string{"testMode"})
	checkInvoke(t, stub, "AddCounterParty", []string{})
	checkInvoke(t, stub, "AddCounterParty", []string{})
	checkInvoke(t, stub, "AddCounterParty", []string{})
	checkInvoke(t, stub, "AddClaim", []string{"1", "2", "3.14"})
	checkInvoke(t, stub, "AddClaim", []string{"1", "2", "3.14"})
	checkQuery(t, stub, "Claims", []string{"1"}, referenceString)
}

func TestNettingChaincode_testReferenceTable(t *testing.T) {
	log.Info("\n\nReference table stats test")
	scc := new(Chaincode)
	stub := shim.NewMockStub("netting", scc)
	referenceStats := netting.NettingTableStats{
		NumberOfCounterParties: 10,
		NumberOfClaims:         44,
		MetricL1:               53.44444444444444,
		MetricL2:               64.00086804966875,
		SumH:                   0.0,
	}
	referenceBytes, _ := json.Marshal(referenceStats)
	//calls
	checkInit(t, stub, []string{"testMode"})
	// adds 10
	for i := 0; i < 10; i++ {
		checkInvoke(t, stub, "AddCounterParty", []string{})
	}

	checkInvoke(t, stub, "AddClaim", []string{"0", "5", "55.0"})
	checkInvoke(t, stub, "AddClaim", []string{"0", "6", "20.0"})
	checkInvoke(t, stub, "AddClaim", []string{"0", "2", "115.0"})
	checkInvoke(t, stub, "AddClaim", []string{"0", "3", "30.0"})
	checkInvoke(t, stub, "AddClaim", []string{"0", "4", "65.0"})
	checkInvoke(t, stub, "AddClaim", []string{"1", "3", "70.0"})
	checkInvoke(t, stub, "AddClaim", []string{"1", "4", "85.0"})
	checkInvoke(t, stub, "AddClaim", []string{"1", "5", "65.0"})
	checkInvoke(t, stub, "AddClaim", []string{"1", "8", "40.0"})
	checkInvoke(t, stub, "AddClaim", []string{"1", "0", "70.0"})
	checkInvoke(t, stub, "AddClaim", []string{"2", "7", "80.0"})
	checkInvoke(t, stub, "AddClaim", []string{"2", "9", "100.0"})
	checkInvoke(t, stub, "AddClaim", []string{"2", "1", "60.0"})
	checkInvoke(t, stub, "AddClaim", []string{"2", "8", "20.0"})
	checkInvoke(t, stub, "AddClaim", []string{"2", "3", "50.0"})
	checkInvoke(t, stub, "AddClaim", []string{"2", "4", "110.0"})
	checkInvoke(t, stub, "AddClaim", []string{"2", "6", "35.0"})
	checkInvoke(t, stub, "AddClaim", []string{"3", "5", "5.0"})
	checkInvoke(t, stub, "AddClaim", []string{"3", "6", "30.0"})
	checkInvoke(t, stub, "AddClaim", []string{"3", "9", "130.0"})
	checkInvoke(t, stub, "AddClaim", []string{"4", "3", "155.0"})
	checkInvoke(t, stub, "AddClaim", []string{"4", "6", "30.0"})
	checkInvoke(t, stub, "AddClaim", []string{"4", "8", "30.0"})
	checkInvoke(t, stub, "AddClaim", []string{"5", "4", "45.0"})
	checkInvoke(t, stub, "AddClaim", []string{"5", "9", "30.0"})
	checkInvoke(t, stub, "AddClaim", []string{"5", "2", "80.0"})
	checkInvoke(t, stub, "AddClaim", []string{"5", "8", "70.0"})
	checkInvoke(t, stub, "AddClaim", []string{"6", "1", "55.0"})
	checkInvoke(t, stub, "AddClaim", []string{"6", "5", "15.0"})
	checkInvoke(t, stub, "AddClaim", []string{"7", "0", "5.0"})
	checkInvoke(t, stub, "AddClaim", []string{"7", "3", "95.0"})
	checkInvoke(t, stub, "AddClaim", []string{"7", "4", "65.0"})
	checkInvoke(t, stub, "AddClaim", []string{"7", "5", "20.0"})
	checkInvoke(t, stub, "AddClaim", []string{"7", "6", "25.0"})
	checkInvoke(t, stub, "AddClaim", []string{"7", "9", "40.0"})
	checkInvoke(t, stub, "AddClaim", []string{"8", "6", "35.0"})
	checkInvoke(t, stub, "AddClaim", []string{"8", "7", "45.0"})
	checkInvoke(t, stub, "AddClaim", []string{"8", "0", "15.0"})
	checkInvoke(t, stub, "AddClaim", []string{"8", "3", "50.0"})
	checkInvoke(t, stub, "AddClaim", []string{"8", "9", "65.0"})
	checkInvoke(t, stub, "AddClaim", []string{"9", "1", "10.0"})
	checkInvoke(t, stub, "AddClaim", []string{"9", "4", "30.0"})
	checkInvoke(t, stub, "AddClaim", []string{"9", "6", "115.0"})
	checkInvoke(t, stub, "AddClaim", []string{"9", "0", "45.0"})

	checkQuery(t, stub, "Stats", []string{}, string(referenceBytes))
}

func TestNettingChaincode_testNetting(t *testing.T) {
	log.Info("\n\nReference table + Netting stats test")
	scc := new(Chaincode)
	stub := shim.NewMockStub("netting", scc)
	//calls
	checkInit(t, stub, []string{"testMode"})
	// adds 10
	for i := 0; i < 10; i++ {
		checkInvoke(t, stub, "AddCounterParty", []string{})
	}

	checkInvoke(t, stub, "AddClaim", []string{"0", "5", "55.0"})
	checkInvoke(t, stub, "AddClaim", []string{"0", "6", "20.0"})
	checkInvoke(t, stub, "AddClaim", []string{"0", "2", "115.0"})
	checkInvoke(t, stub, "AddClaim", []string{"0", "3", "30.0"})
	checkInvoke(t, stub, "AddClaim", []string{"0", "4", "65.0"})
	checkInvoke(t, stub, "AddClaim", []string{"1", "3", "70.0"})
	checkInvoke(t, stub, "AddClaim", []string{"1", "4", "85.0"})
	checkInvoke(t, stub, "AddClaim", []string{"1", "5", "65.0"})
	checkInvoke(t, stub, "AddClaim", []string{"1", "8", "40.0"})
	checkInvoke(t, stub, "AddClaim", []string{"1", "0", "70.0"})
	checkInvoke(t, stub, "AddClaim", []string{"2", "7", "80.0"})
	checkInvoke(t, stub, "AddClaim", []string{"2", "9", "100.0"})
	checkInvoke(t, stub, "AddClaim", []string{"2", "1", "60.0"})
	checkInvoke(t, stub, "AddClaim", []string{"2", "8", "20.0"})
	checkInvoke(t, stub, "AddClaim", []string{"2", "3", "50.0"})
	checkInvoke(t, stub, "AddClaim", []string{"2", "4", "110.0"})
	checkInvoke(t, stub, "AddClaim", []string{"2", "6", "35.0"})
	checkInvoke(t, stub, "AddClaim", []string{"3", "5", "5.0"})
	checkInvoke(t, stub, "AddClaim", []string{"3", "6", "30.0"})
	checkInvoke(t, stub, "AddClaim", []string{"3", "9", "130.0"})
	checkInvoke(t, stub, "AddClaim", []string{"4", "3", "155.0"})
	checkInvoke(t, stub, "AddClaim", []string{"4", "6", "30.0"})
	checkInvoke(t, stub, "AddClaim", []string{"4", "8", "30.0"})
	checkInvoke(t, stub, "AddClaim", []string{"5", "4", "45.0"})
	checkInvoke(t, stub, "AddClaim", []string{"5", "9", "30.0"})
	checkInvoke(t, stub, "AddClaim", []string{"5", "2", "80.0"})
	checkInvoke(t, stub, "AddClaim", []string{"5", "8", "70.0"})
	checkInvoke(t, stub, "AddClaim", []string{"6", "1", "55.0"})
	checkInvoke(t, stub, "AddClaim", []string{"6", "5", "15.0"})
	checkInvoke(t, stub, "AddClaim", []string{"7", "0", "5.0"})
	checkInvoke(t, stub, "AddClaim", []string{"7", "3", "95.0"})
	checkInvoke(t, stub, "AddClaim", []string{"7", "4", "65.0"})
	checkInvoke(t, stub, "AddClaim", []string{"7", "5", "20.0"})
	checkInvoke(t, stub, "AddClaim", []string{"7", "6", "25.0"})
	checkInvoke(t, stub, "AddClaim", []string{"7", "9", "40.0"})
	checkInvoke(t, stub, "AddClaim", []string{"8", "6", "35.0"})
	checkInvoke(t, stub, "AddClaim", []string{"8", "7", "45.0"})
	checkInvoke(t, stub, "AddClaim", []string{"8", "0", "15.0"})
	checkInvoke(t, stub, "AddClaim", []string{"8", "3", "50.0"})
	checkInvoke(t, stub, "AddClaim", []string{"8", "9", "65.0"})
	checkInvoke(t, stub, "AddClaim", []string{"9", "1", "10.0"})
	checkInvoke(t, stub, "AddClaim", []string{"9", "4", "30.0"})
	checkInvoke(t, stub, "AddClaim", []string{"9", "6", "115.0"})
	checkInvoke(t, stub, "AddClaim", []string{"9", "0", "45.0"})
	checkInvoke(t, stub, "RunNetting", []string{})

	initial := netting.NettingTableStats{
		NumberOfCounterParties: 10,
		NumberOfClaims:         44,
		MetricL1:               53.44444444444444,
		MetricL2:               64.00086804966875,
		SumH:                   0.0,
	}

	var stats netting.NettingTableStats
	bts, _ := stub.MockQuery("Stats", []string{})
	_ = json.Unmarshal(bts, &stats)
	if stats.SumH != 0.0 ||
		stats.NumberOfCounterParties != initial.NumberOfCounterParties ||
		stats.NumberOfClaims >= initial.NumberOfClaims ||
		stats.MetricL1 >= initial.MetricL1 ||
		stats.MetricL2 >= initial.MetricL2 {
		t.FailNow()
	}
}
