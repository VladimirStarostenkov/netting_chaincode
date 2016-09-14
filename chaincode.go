package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("chaincode")

func checkCriticalError(e error) {
	if e != nil {
		log.Error(e.Error())
		panic(e)
	}
}

// NettingChaincode implementation
type Chaincode struct {
}

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("Error starting Netting chaincode: %s", err)
	}
}

func (t *Chaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	log.Debugf("Init called with function name: %s, with arguments: %s", function, args)

	return initSmartContract(stub, []string{})
}

func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	log.Debugf("Invoke called with function name: %s, with arguments: %s", function, args)

	f, ok := invokes[function]
	if ok {
		s := smartContract{}
		return f(s, stub, args)
	}
	return nil, nil
}

func (t *Chaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	log.Debugf("Query called with function name: %s, with arguments: %s", function, args)

	f, ok := queries[function]
	if ok {
		s := smartContract{}
		return f(s, stub, args)
	}
	return nil, nil
}