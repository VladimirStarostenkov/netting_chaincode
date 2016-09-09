package main

import (
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"errors"
	"github.com/VladimirStarostenkov/netting"
)

type smartContract struct {
	netting.NettingTable
	storeKey string
}

func (this *smartContract) init(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	log.Debugf("init called with args: %s\n", args)

	this.NettingTable.Init()
	this.storeKey = "NettingTable"

	if err := this.save(stub); err != nil {
		return nil, err
	}
	return nil, nil
}
// args: From int, To int, Value float
func (this *smartContract) invokeAddClaim(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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
	err = this.load(stub)
	checkCriticalError(err)

	this.NettingTable.AddClaim(from, to, value)

	// Save new data
	err = this.save(stub)
	checkCriticalError(err)

	return nil, nil
}
// args: -
func (this *smartContract) invokeAddNode(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	log.Debugf("invokeAddNode called with args: %s\n", args)

	// Load existing data
	err := this.load(stub)
	checkCriticalError(err)

	_ = this.NettingTable.AddCounterParty()

	// Save new data
	err = this.save(stub)
	checkCriticalError(err)

	return nil, nil
}
// args: -
func (this *smartContract) invokeRunNetting(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	log.Debugf("invokeRunNetting called with args: %s\n", args)

	// Load existing data
	err := this.load(stub)
	checkCriticalError(err)

	// Run netting algorithm
	this.NettingTable.Optimize()

	// Save new data
	err = this.save(stub)
	checkCriticalError(err)

	return nil, nil
}
// args: -
func (this *smartContract) queryStats(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	log.Debugf("queryStats called with args: %s\n", args)

	// Load existing data
	err := this.load(stub)
	checkCriticalError(err)

	return this.GetStats(), nil
}
// args: CounterPartyId int
func (this *smartContract) queryClaims(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	message := fmt.Sprintf("queryClaims called with args: %s\n", args)
	log.Debugf(message)

	if len(args) < 1 {
		log.Errorf(message)
		return nil, errors.New(message)
	}

	counterPartyId, err := strconv.Atoi(args[0])
	if err != nil {
		log.Errorf("strconv.Atoi(args[0]) error: %s", err.Error())
		return nil, err
	}

	// Load existing data
	err = this.load(stub)
	checkCriticalError(err)

	return this.NettingTable.GetClaims(counterPartyId), nil
}

func (this *smartContract) save(stub shim.ChaincodeStubInterface) (error) {
	log.Debugf("Saving...\n")

	// Data to Bytes
	bytes, err := this.NettingTable.ToBytes()
	if err != nil {
		log.Errorf("this.netting.ToBytes() error: %s", err.Error())
		return err
	}
	// Save Bytes
	err = stub.PutState(this.storeKey, bytes)
	if err != nil {
		log.Errorf("stub.PutState(this.key, bytes) error: %s", err.Error())
		return err
	}
	log.Debugf("Saved data : %s\n", bytes)

	return nil
}

func (this *smartContract) load(stub shim.ChaincodeStubInterface) (error) {
	log.Debugf("Loading...\n")

	bytes, err := stub.GetState(this.storeKey)
	if err != nil {
		log.Errorf("stub.GetState(this.storeKey) error: %s", err.Error())
		return err
	}

	this.NettingTable = netting.NettingTable{}
	err = this.NettingTable.InitFromBytes(bytes)
	if err != nil {
		log.Errorf("this.NettingTable.InitFromBytes(bytes) error: %s", err.Error())
		return err
	}

	return err
}