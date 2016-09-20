package main

import (
	"errors"
	"fmt"
	"github.com/VladimirStarostenkov/netting"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
)

const storeKey string = "NettingTable"
const adminRole string = "admin"

var testMode bool = false

var invokes map[string]func(smartContract, shim.ChaincodeStubInterface, []string) ([]byte, error) = map[string]func(smartContract, shim.ChaincodeStubInterface, []string) ([]byte, error){
	"AddClaim":        (smartContract).invoke_AddClaim,
	"AddCounterParty": (smartContract).invoke_AddCounterParty,
	"RunNetting":      (smartContract).invoke_RunNetting,
	"Clear":           (smartContract).invoke_Clear,
}

var queries map[string]func(smartContract, shim.ChaincodeStubInterface, []string) ([]byte, error) = map[string]func(smartContract, shim.ChaincodeStubInterface, []string) ([]byte, error){
	"Stats":  (smartContract).query_Stats,
	"Graph":  (smartContract).query_Graph,
	"Claims": (smartContract).query_Claims,
}

type smartContract struct {
}

func initSmartContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	log.Debugf("init called with args: %s\n", args)

	if (len(args) > 0) && (args[0] == "testMode") {
		testMode = true
	}

	if err := verifyCallerRole(stub, adminRole); err != nil {
		return nil, err
	}

	nettingTable := netting.NettingTable{}
	nettingTable.Init()

	if err := save(&nettingTable, stub); err != nil {
		return nil, err
	}
	return nil, nil
}

// args: From int, To int, Value float
func (smartContract) invoke_AddClaim(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	message := fmt.Sprintf("invokeAddClaim called with args: %s\n", args)
	log.Debugf(message)

	// Check arguments
	if len(args) < 3 {
		log.Errorf(message)
		return nil, errors.New(message)
	}
	from, err := strconv.Atoi(args[0])
	if err != nil {
		log.Errorf("strconv.Atoi(args[0]) error: %s", err.Error())
		return nil, err
	}
	to, err := strconv.Atoi(args[1])
	if err != nil {
		log.Errorf("strconv.Atoi(args[1]) error: %s", err.Error())
		return nil, err
	}
	value, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		log.Errorf("strconv.ParseFloat(args[2], 64) error: %s", err.Error())
		return nil, err
	}

	// We are not interested in "negative claims"
	if value < 0.0 {
		return nil, nil
	}

	// Load existing data
	nettingTable, err := load(stub)
	checkCriticalError(err)

	if err := verifyCallerCompany(stub, strconv.Itoa(from)); err != nil {
		return nil, err
	}

	nettingTable.AddClaim(from, to, value)

	// Save new data
	err = save(nettingTable, stub)
	checkCriticalError(err)

	return nil, nil
}

// args: -
func (smartContract) invoke_AddCounterParty(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	log.Debugf("invokeAddNode called with args: %s\n", args)

	if err := verifyCallerRole(stub, adminRole); err != nil {
		return nil, err
	}

	// Load existing data
	nettingTable, err := load(stub)
	checkCriticalError(err)

	_ = nettingTable.AddCounterParty()

	// Save new data
	err = save(nettingTable, stub)
	checkCriticalError(err)

	return nil, nil
}

// args: -
func (smartContract) invoke_RunNetting(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	log.Debugf("invokeRunNetting called with args: %s\n", args)

	if err := verifyCallerRole(stub, adminRole); err != nil {
		return nil, err
	}

	// Load existing data
	nettingTable, err := load(stub)
	checkCriticalError(err)

	// Run netting algorithm
	nettingTable.Optimize()

	// Save new data
	err = save(nettingTable, stub)
	checkCriticalError(err)

	return nil, nil
}

// args: -
func (smartContract) invoke_Clear(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if err := verifyCallerRole(stub, adminRole); err != nil {
		return nil, err
	}
	return initSmartContract(stub, args)
}

// args: -
func (smartContract) query_Stats(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	log.Debugf("queryStats called with args: %s\n", args)

	// Load existing data
	nettingTable, err := load(stub)
	checkCriticalError(err)

	return nettingTable.GetStats(), nil
}

// args: -
func (smartContract) query_Graph(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	log.Debugf("queryGraph called with args: %s\n", args)

	if err := verifyCallerRole(stub, adminRole); err != nil {
		return nil, err
	}

	// Load existing data
	nettingTable, err := load(stub)
	checkCriticalError(err)

	bts, err := nettingTable.ToBytes()
	if err != nil {
		log.Errorf("nettingTable.ToBytes() error: \n", err.Error())
		return nil, err
	}

	return bts, nil
}

// args: CounterPartyId int
func (smartContract) query_Claims(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	message := fmt.Sprintf("queryClaims called with args: %s\n", args)
	log.Debugf(message)

	if len(args) < 1 {
		log.Errorf(message)
		return nil, errors.New(message)
	}

	if err := verifyCallerCompany(stub, args[0]); err != nil {
		return nil, err
	}

	counterPartyId, err := strconv.Atoi(args[0])
	if err != nil {
		log.Errorf("strconv.Atoi(args[0]) error: %s", err.Error())
		return nil, err
	}

	// Load existing data
	nettingTable, err := load(stub)
	checkCriticalError(err)

	return nettingTable.GetClaims(counterPartyId), nil
}

func save(this *netting.NettingTable, stub shim.ChaincodeStubInterface) error {
	log.Debugf("Saving...\n")

	// Data to Bytes
	bytes, err := this.ToBytes()
	if err != nil {
		log.Errorf("this.netting.ToBytes() error: %s", err.Error())
		return err
	}
	// Save Bytes
	err = stub.PutState(storeKey, bytes)
	if err != nil {
		log.Errorf("stub.PutState(this.key, bytes) error: %s", err.Error())
		return err
	}
	log.Debugf("Saved data : %s\n", bytes)

	return nil
}

func load(stub shim.ChaincodeStubInterface) (*netting.NettingTable, error) {
	log.Debugf("Loading...\n")

	bytes, err := stub.GetState(storeKey)
	if err != nil {
		log.Errorf("stub.GetState(storeKey) error: %s", err.Error())
		return nil, err
	}

	result := netting.NettingTable{}
	err = result.InitFromBytes(bytes)
	if err != nil {
		log.Errorf("this.NettingTable.InitFromBytes(bytes) error: %s", err.Error())
		return nil, err
	}

	return &result, nil
}
