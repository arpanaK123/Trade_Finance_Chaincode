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
*/

package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	// "bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
// SimpleAsset implements a simple chaincode to manage an asset
type SimpleChaincode struct {
}

type AccountStructure struct{
	account_Number string `json:"account_Number"`
	account_Holder_Name string `json:"account_Holder_Name"`
	account_Balance string `json:"account_Balance"`
}
type ContractStructure struct{
	contract_Id string `json:"contract_Id"`
	importer_Name string  `json:"importer_Name"`
	exporter_Name string `json:"exporter_Name"`
	port_Authority string `json:"port_Authority"`
	custom_Authority string `json:"custom_Authority"`
	importer_Bank_Name string `json:"importer_Bank_Name"`
	insurance_Name string `json:"insurance_Name"`
}
// func (s *SimpleChaincode) Init(APIstub shim.ChaincodeStubInterface) pb.Response{
// 	return shim.Success(nil)
// }

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("exporter Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "invoke" {
		// Make payment of X units from A to B
		return t.invoke(stub, args)
	}else if function == "create_Account" {
		// Make payment of X units from A to B
		return t.create_Account(stub, args)
	} else if function == "create_Contract" {
		// Create contract an entity from its state
		return t.create_Contract(stub, args)
	} else if function == "get_Balance_By" {
		// Get balance by acc_no an entity from its state
		return t.get_Balance_By(stub, args)
	}else if function == "get_Account" {
		// get account an entity from its state
		return t.get_Account(stub, args)
	}else if function == "query" {
		// the old "Query" is now implemtned in invoke
		return t.query(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"create_Account\" \"create_Contract\" \"get_Balance_By\"\"get_Account\"\"query\"")
}

func (s *SimpleChaincode) create_Account(APIstub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	var account=AccountStructure{account_Number:args[1], account_Holder_Name: args[2], account_Balance: args[3] }
	accountAsBytes,_:=json.Marshal(account)
	APIstub.PutState(args[0], accountAsBytes)
	return shim.Success(nil)
}
func (s *SimpleChaincode) create_Contract(APIstub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments. Expecting 8")
	}
	var contract=ContractStructure{contract_Id:args[1], importer_Name:args[2], exporter_Name:args[3],  port_Authority:args[4], custom_Authority:args[5], importer_Bank_Name:args[6], insurance_Name:args[7]}
	contractAsBytes,_:=json.Marshal(contract)
	APIstub.PutState(args[0],contractAsBytes)
	return shim.Success(nil)
}
func (s *SimpleChaincode) get_Balance_By(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	return shim.Success(nil)
}
func (s *SimpleChaincode) get_Account(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	return shim.Success(nil)
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("exporter Init")
	_, args := stub.GetFunctionAndParameters()
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var err error

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4 ")
	}

	// Initialize the chaincode
	A = args[0]
	Aval, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	B = args[2]
	Bval, err = strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
// 	fmt.Println("exporter Invoke")
// 	function, args := stub.GetFunctionAndParameters()
// 	if function == "invoke" {
// 		// Make payment of X units from A to B
// 		return t.invoke(stub, args)
// 	} else if function == "delete" {
// 		// Deletes an entity from its state
// 		return t.delete(stub, args)
// 	} else if function == "query" {
// 		// the old "Query" is now implemtned in invoke
// 		return t.query(stub, args)
// 	}

// 	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
// }

//Transaction makes payment of X units from A to B
func (t *SimpleChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var X int          // Transaction value
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	A = args[0]
	B = args[1]

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Avalbytes == nil {
		return shim.Error("Entity not found")
	}
	Aval, _ = strconv.Atoi(string(Avalbytes))

	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Bvalbytes == nil {
		return shim.Error("Entity not found")
	}
	Bval, _ = strconv.Atoi(string(Bvalbytes))

	// Perform the execution
	X, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	Aval = Aval - X
	Bval = Bval + X
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state back to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

// query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
