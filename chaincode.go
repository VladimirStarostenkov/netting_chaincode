package main

import (
	"errors"
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

	return initSmartContract(stub, args)
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

func callerAttribute(stub shim.ChaincodeStubInterface, attributeName string) string {
	value, err := stub.ReadCertAttribute(attributeName)
	if err != nil {
		log.Error("Failed fetching caller's attribute. Error: " + err.Error())
		return ""
	}
	log.Debugf("Caller %s is: %s", attributeName, value)
	return string(value)
}

func verifyCallerAttribute(stub shim.ChaincodeStubInterface, attributeName string, attributeValue string) error {
	if testMode {
		return nil
	}
	value := callerAttribute(stub, attributeName)
	if value != attributeValue {
		return errors.New(fmt.Sprintf("Caller attribute '%s' expected to be: %s, got: %s",
			attributeName, attributeValue, value))
	}
	return nil
}

func verifyCallerRole(stub shim.ChaincodeStubInterface, role string) error {
	return verifyCallerAttribute(stub, "role", role)
}

func verifyCallerCompany(stub shim.ChaincodeStubInterface, company string) error {
	return verifyCallerAttribute(stub, "company", company)
}
