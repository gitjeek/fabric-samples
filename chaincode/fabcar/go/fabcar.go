/**
 * @Author  rjagge@163.com
 * @date	2018.3.7
 * @description	this is a fabric sample demo modified by me, in order to finish my graduation project.
 */
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the std structure, with 4 properties.  Structure tags are used by encoding/json library
type Std struct {
	Number   string `json:"number"`
	Name  string `json:"name"`
	Major string `json:"major"`
	School  string `json:"school"`
}

/*
 * The Init method is called when the Smart Contract "fabstd" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabstd"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryStd" {
		return s.queryStd(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createStd" {
		return s.createStd(APIstub, args)
	} else if function == "queryAllStds" {
		return s.queryAllStds(APIstub)
	} else if function == "changeStdSchool" {
		return s.changeStdSchool(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryStd(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	stdAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(stdAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	stds := []Std{
		Std{Number: "2014141093031", Name: "Fanwei", Major: "ISM", School: "SCU"},
		Std{Number: "2014141093032", Name: "Zhaoying", Major: "ISM", School: "SCU"},
		Std{Number: "2014141093033", Name: "Dongkaining", Major: "ISM", School: "SCU"},
		Std{Number: "2014141093034", Name: "Yuanli", Major: "ISM", School: "SCU"},
		Std{Number: "2014141093035", Name: "Luojun", Major: "ISM", School: "SCU"},
		Std{Number: "2014141093036", Name: "Liuguoxiang", Major: "ISM", School: "SCU"},
		Std{Number: "2014141093037", Name: "Pengguangming", Major: "ISM", School: "SCU"},
		Std{Number: "2014141093038", Name: "Jixuyu", Major: "ISM", School: "SCU"},
		Std{Number: "2014141093039", Name: "Xieheping", Major: "ISM", School: "SCU"},
		Std{Number: "2014141093030", Name: "Liyanrong", Major: "ISM", School: "SCU"},
	}

	i := 0
	for i < len(stds) {
		fmt.Println("i is ", i)
		stdAsBytes, _ := json.Marshal(stds[i])
		APIstub.PutState("STD"+strconv.Itoa(i), stdAsBytes)
		fmt.Println("Added", stds[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createStd(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var std = Std{Number: args[1], Name: args[2], Major: args[3], School: args[4]}

	stdAsBytes, _ := json.Marshal(std)
	APIstub.PutState(args[0], stdAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllStds(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "STD0"
	endKey := "STD999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllStds:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) changeStdSchool(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	stdAsBytes, _ := APIstub.GetState(args[0])
	std := Std{}

	json.Unmarshal(stdAsBytes, &std)
	std.School = args[1]

	stdAsBytes, _ = json.Marshal(std)
	APIstub.PutState(args[0], stdAsBytes)

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

