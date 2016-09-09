package main

import (
	"fmt"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("netting-chaincode")

func checkCriticalError(e error) {
	if e != nil {
		log.Error(e.Error())
		panic(e)
	}
}

// NettingChaincode implementation
type NettingChaincode struct {
	invokes map[string]func(*smartContract, shim.ChaincodeStubInterface, []string) ([]byte, error)
	queries map[string]func(*smartContract, shim.ChaincodeStubInterface, []string) ([]byte, error)
}

func (self *NettingChaincode) setMethods() {
	self.invokes = map[string]func(*smartContract, shim.ChaincodeStubInterface, []string) ([]byte, error) {
		"AddClaim":(*smartContract).invokeAddClaim,
		"AddNode":(*smartContract).invokeAddNode,
		"RunNetting":(*smartContract).invokeRunNetting,
	}
	self.queries = map[string]func(*smartContract, shim.ChaincodeStubInterface, []string) ([]byte, error) {
		"Stats":(*smartContract).queryStats,
		"Claims":(*smartContract).queryClaims,
	}
	log.Debug("setMethods: ")
	log.Debugf("        %+v\n", self.invokes)
	log.Debugf("        %+v\n", self.queries)
}

func main() {
	err := shim.Start(new(NettingChaincode))
	if err != nil {
		fmt.Printf("Error starting Netting chaincode: %s", err)
	}
}

func (t *NettingChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	log.Debugf("Init called with function name: %s, with arguments: %s", function, args)

	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}

	t.setMethods()
	s := &smartContract{}

	return s.init(stub, []string{})
}

func (t *NettingChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	log.Debugf("Invoke called with function name: %s, with arguments: %s", function, args)

	f, ok := t.invokes[function]
	if ok {
		s := smartContract{}
		return f(&s, stub, args)
	}
	return nil, nil
}

func (t *NettingChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	log.Debugf("Query called with function name: %s, with arguments: %s", function, args)

	f, ok := t.queries[function]
	if ok {
		s := smartContract{}
		return f(&s, stub, args)
	}
	return nil, nil
}