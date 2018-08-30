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
//Logger := shim.NewLogger("logger")
type AccountStructure struct{
	Account_Number string `json:"account_Number"`
	Account_Holder_Name string `json:"account_Holder_Name"`
	Account_Balance string `json:"account_Balance"`
	Bank_Name string `json:"bank_Name"`
}
type ContractStructure struct{
	Contract_Id string `json:"contract_Id"`
	Content_Description string`json:"content_Description"`
	Value string`json:"value"`
	Importer_Bank_Name string  `json:"importer_Bank_Name"`
	Exporter_Bank_Name string `json:"exporter_Bank_Name"`
	Custom_Authority string `json:"custom_Authority"`
	Port_Of_Loading string `json:"port_Of_Loading"`
	Port_Of_Entry string `json:"port_Of_Entry"`
}
// func (s *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response{
// 	return shim.Success(nil)
// }

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	//Logger.Info("invoke called")
	fmt.Println("exporter Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "invoke" {
		return t.invoke(stub, args)
	}else if function == "create_Account" {
		return t.create_Account(stub, args)
	}else if function == "create_Contract" {
		return t.create_Contract(stub, args)
	} else if function == "get_Contract_By" {
		return t.get_Contract_By(stub, args)
	}else if function == "get_Balance_By" {
		return t.get_Balance_By(stub, args)
	}else if function == "get_Account" {
		return t.get_Account(stub, args)
	}else if function == "accept_By_Importer" {
		return t.accept_By_Importer(stub, args)
	}else if function == "accept_By_Exporter" {
		return t.accept_By_Exporter(stub, args)
	}else if function == "accept_By_Custom" {
		return t.accept_By_Custom(stub, args)
	}else if function == "accept_By_ImporterBank" {
		return t.accept_By_ImporterBank(stub, args)
	}else if function == "accept_By_Insurance" {
		return t.accept_By_Insurance(stub, args)
	}else if function == "query" {
		// the old "Query" is now implemtned in invoke
		return t.query(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"create_Account\" \"create_Contract\" \"get_Balance_By\"\"get_Account\"\"accept_By_Importer\"\"accept_By_Exporter\"\"accept_By_Custom\"\"accept_By_ImporterBank\"\"accept_By_Insurance\"\"query\"")
}

func (s *SimpleChaincode) create_Account(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	account := AccountStructure{Account_Number: args[0], Account_Holder_Name: args[1], Account_Balance: args[2], Bank_Name: args[3], }
	accountAsBytes,error:=json.Marshal(account)
	if error!=nil{
	return shim.Error("something wrong")
	}
	stub.PutState(args[0], accountAsBytes)
	return shim.Success(accountAsBytes)
}
func (s *SimpleChaincode) create_Contract(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments. Expecting 8")
	}
	contract := ContractStructure{Contract_Id:args[0], Content_Description:args[1], Value:args[2],  Importer_Bank_Name:args[3], Exporter_Bank_Name:args[4], Custom_Authority:args[5], Port_Of_Loading:args[6], Port_Of_Entry:args[7]}
	contractAsBytes,error:=json.Marshal(contract)
	if error!=nil{
		return shim.Error("something wrong")
	}
	stub.PutState(args[0],contractAsBytes)
	return shim.Success(contractAsBytes)
}
func (s *SimpleChaincode) get_Balance_By(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	 account_Number := args[0]

	 accountAsByte, accountError := stub.GetState(account_Number)
	 if accountError != nil {
	 	return shim.Error("error.... account information in getstate method ")
	 }
	  accountSructure:= AccountStructure{}
	 errorAccount := json.Unmarshal(accountAsByte, &accountSructure)
	 if errorAccount != nil {
	 	return shim.Error("Something Wrong")
	 }
	 return shim.Success([]byte(accountSructure.Account_Balance))
}
func (s *SimpleChaincode) get_Account(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	 account_Number := args[0]

	 accountAsByte, accountError := stub.GetState(account_Number)
	 if accountError != nil {
	 	return shim.Error("something wrong in getstate method ")
	 }
	 accountSructure := AccountStructure{}
	 errorAccount := json.Unmarshal(accountAsByte, &accountSructure)
	 if errorAccount != nil {
	 	return shim.Error("something wrong")
	 }
	 return shim.Success(accountAsByte)
}
func (s *SimpleChaincode) get_Contract_By(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	contract_Id := args[0]

	contractAsByte, contractError := stub.GetState(contract_Id)
	if contractError != nil {
		return shim.Error("something wrong in getstate method ")
	}
	contractSructure := ContractStructure{}
	errorContract := json.Unmarshal(contractAsByte, &contractSructure)
	if errorContract != nil {
		return	shim.Error("something wrong")
	}
	return shim.Success(contractAsByte)
}

func (s *SimpleChaincode) accept_By_Exporter(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	response:="success"
	return shim.Success([]byte(response))
}

func (s *SimpleChaincode) accept_By_Importer(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	account_Number := args[0]
	accountAsByte, accountError := stub.GetState(account_Number)
	if accountError != nil {
		return shim.Error("Account number Invalid ")
	}
	accountSructure := AccountStructure{}
	errorAccount := json.Unmarshal(accountAsByte, &accountSructure)
	if accountSructure.Account_Balance!="10000"{
		return shim.Error("account balance is less than 10000")
	}
	//errorAccount := json.Unmarshal(accountAsByte, &accountSructure)
	if errorAccount != nil {
		return shim.Error("Account number Invalid")
	}
	response:="success"
	return shim.Success([]byte(response))
}

func (s *SimpleChaincode) accept_By_Custom(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	contract_Id:=args[0]
	 customAsByte, customError := stub.GetState(contract_Id)
	 if customError != nil {
	 	return shim.Error("Contract_Id is Invalid ")
	 }
	 contractSruct := ContractStructure{}
	 errorContract_Custom := json.Unmarshal(customAsByte, &contractSruct)
	 if contractSruct.Port_Of_Loading!="India" && contractSruct.Port_Of_Entry!="USA"{
		return shim.Error("Port_Of_Loading or Port_Of_Entry is Invalid")
	 }
	 //errorContract_Custom := json.Unmarshal(customAsByte, &contractSruct)
	 if errorContract_Custom != nil {
	 	return shim.Error("something wrong")
	 }
	 response:="success"
	 return shim.Success([]byte(response))}

func (s *SimpleChaincode) accept_By_ImporterBank(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	account_Number := args[0]
    account_Balance:="10000"
	accountAsByte, accountError := stub.GetState(account_Number)
	if accountError != nil {
		return shim.Error("Account number is Invalid")
	}
	accountSruct := AccountStructure{}
	errorAccount := json.Unmarshal(accountAsByte, &accountSruct)
	if accountSruct.Account_Balance!=account_Balance{
		return shim.Error("account balance is less than 10000")
	}
	//errorAccount := json.Unmarshal(accountAsByte, &accountSruct)
	if errorAccount != nil {
	return	shim.Error("something wrong")
	}
	response:="success"
	return shim.Success([]byte(response))}

func (s *SimpleChaincode) accept_By_Insurance(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	contract_Id:=args[0]
	 customAsByte, customError := stub.GetState(contract_Id)
	 if customError != nil {
	 	return shim.Error("something wrong in getstate method ")
	 }
	 contractSruct := ContractStructure{}
	 errorContract_Custom := json.Unmarshal(customAsByte, &contractSruct)
	 if contractSruct.Port_Of_Loading!="India" && contractSruct.Port_Of_Entry!="USA"{
		return shim.Error("Port_Of_Loading or Port_Of_Entry is Invalid")
	 }
	 //errorContract_Custom := json.Unmarshal(customAsByte, &contractSruct)
	 if errorContract_Custom != nil {
	 	return shim.Error("something wrong")
	 }
	 response:="success"
	 return shim.Success([]byte(response))}

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
