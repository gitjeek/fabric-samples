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
	State   string `json:"state"`
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
	} else if function == "modifyStdInfo" {
		return s.modifyStdInfo(APIstub, args)
	} else if function == "deleteStd" {
		return s.deleteStd(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name : "+function)
}

//在区块链节点调用该函数时，需要传入指定学生的ID，函数通过学生ID查询学生信息。查询成功规则返回学生信息，查询出现错误或者学生ID参数传入错误，函数则会报错。
func (s *SmartContract) queryStd(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	stdAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(stdAsBytes)
}

//在区块链网络第一次初始化时，系统会调用初始化账本函数（initLedger）。在初始化账本是，会写入系统自带的两位学生的信息，以方便在试验阶段检测系统功能完善性。在系统进入真实使用环境时，应该把初始化的两位学生信息删除掉。
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	stds := []Std{
		Std{Number: "2014141093031", Name: "FanMing", Major: "ISM", School: "SCU", State: "valid"},
		Std{Number: "2014141093032", Name: "ZhaoLei", Major: "ISM", School: "SCU", State: "valid"},
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

//系统通过调用创建学生（createStd）函数来创建学生。在调用该函数时，需要传入学生的学号、姓名、专业、学校四个参数。创建成果会返回成功消息，如果参数传入错误、或者学生创建失败，函数则会报错。
func (s *SmartContract) createStd(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var std = Std{Number: args[1], Name: args[2], Major: args[3], School: args[4], State: "valid"}

	stdAsBytes, _ := json.Marshal(std)
	APIstub.PutState(args[0], stdAsBytes)

	return shim.Success(nil)
}

//显示所有的学生信息功能通过调用查询所有学生信息（queryAllStds）函数实现。调用此函数不需要传入任何参数，并且默认返回前1000位学生的信息（鉴于系统处于试验阶段），以json的格式返回。
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

//系统通过调用修改学生信息（modifyStdInfo）函数来对学生的信息进行修改。调用此函数需要传入学生ID以及学生的所有信息（包括需要修改的和未修改的）。如果某一个字段需要修改，那么传入新的值；如果某一个值不需要修改，则传入“NoChange”字段。函数调用成功后将修改学生的信息，并返回成功消息。如果函数传入参数不规范或者修改失败，程序将报错。
func (s *SmartContract) modifyStdInfo(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	stdAsBytes, _ := APIstub.GetState(args[0])

	std := Std{}

	json.Unmarshal(stdAsBytes, &std)

	if args[1] != "NoChange"{
		std.Major = args[1]
	}
	if args[2] != "NoChange"{
		std.Name = args[2]
	}
	if args[3] != "NoChange"{
		std.Number = args[3]
	}
	if args[4] != "NoChange"{
		std.School = args[4]
	}

	stdAsBytes, _ = json.Marshal(std)
	APIstub.PutState(args[0], stdAsBytes)

	return shim.Success(nil)
}

//系统通过调用删除学生信息（deleteStd）函数来对学生信息进行“删除”。这里需要注意的是，在Hyperledger Fabric（即区块链网络）中，并没有真正的删除功能。每一条数据操作和数据更新记录都写在区块中，无法修改和抹除。所以这里的删除功能是将学生的“数据有效性”字段（state）从valid修改为invalid。调用此函数需要传入学生的ID。如果函数调用成功，则成功删除学生的信息，并返回成功消息。如果参数传入不规范或者删除失败，程序将会报错。
func (s *SmartContract) deleteStd(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	stdAsBytes, _ := APIstub.GetState(args[0])

	std := Std{}

	json.Unmarshal(stdAsBytes, &std)

	std.State = "invalid"

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
