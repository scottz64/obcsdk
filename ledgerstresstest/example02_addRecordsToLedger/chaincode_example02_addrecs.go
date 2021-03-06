
/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Modified by Ratnakar Asara, and Scott Zwierzynski.

*/

package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"errors"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"github.com/hyperledger/fabric/core/chaincode"

	//"hyperledger/ccs"
	//main "github.com/scottz64/obcsdk/blob/master/ledgerstresstest/example02_addRecordsToLedger"
	//main "github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02"
)

type Chaincode_example02_addrecs struct {
}

func (t *Chaincode_example02_addrecs) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var A, B string    // Entities for the Deploy/Init
	var Bval int
	var err error

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	// Initialize the chaincode
	A = args[0]      // "a"
	Aval := args[1]  // DATA : FixedString or RandomString
	B = args[2]      // "counter"
	Bval, err = strconv.Atoi(args[3])    // cntr value integer
	if err != nil {
		return nil, errors.New("Expecting integer value for counter index")
	}
	fmt.Printf("Aval (INIT DATA STRING) = %s, Bval (INIT counter value) = %d\n", Aval, Bval)

	// Write the state to the ledger
	err = stub.PutState(A, []byte(Aval))    // "a" , DATA
	if err != nil {
		return nil, err
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))  // "counter" , counter value (typically zero)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Transaction makes payment of X units from A to B
func (t *Chaincode_example02_addrecs) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	}

	var A, B, Aval string    // Entities
	var Bval int // Asset holdings
	var err error

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	A = args[0]       // aN ("a0" or "a1" or whatever, with the number being the cntr value, the top of the ledger stack)
	Aval = args[1]    // DATA
	B = args[2]       // "counter"


	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	if Bvalbytes == nil {
		return nil, errors.New("Entity not found")
	}
	Bval, _ = strconv.Atoi(string(Bvalbytes))

	// Perform the execution
	Bval = Bval + 1
	fmt.Printf("Aval (DATA) = %s, Bval (counter++ value) = %d\n", Aval, Bval)

	// Write the state back to the ledger
	err = stub.PutState(A, []byte(Aval))     // aN, DATA
	if err != nil {
		return nil, err
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))   // "counter", incremented_counter_value
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Deletes an entity from state
func (t *Chaincode_example02_addrecs) delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *Chaincode_example02_addrecs) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function != "query" {
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}
	var A string // Entities
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return nil, errors.New(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return nil, errors.New(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return Avalbytes, nil
}

func main() {
	self := &Chaincode_example02_addrecs{}
	err := shim.Start(self) // Our one instance implements both Transactions and Queries interfaces
	if err != nil {
		fmt.Printf("Error starting chaincode_example02_addrecs chaincode: %s", err)
	}
}

