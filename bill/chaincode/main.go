package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	peer "github.com/hyperledger/fabric/protos/peer"
	"fmt"
)

type BillChainCode struct {

}

func (t *BillChainCode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("########### HaoServiceChaincode Init ###########")

	// Get the function and arguments from the request
	function, _ := stub.GetFunctionAndParameters()

	// Check if the request is the init function
	if function != "init" {
		return shim.Error("Unknown function call")
	}

	// Put in the ledger the key/value hello/world
	err := stub.PutState("hello", []byte("world"))
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return a successful message
	return shim.Success(nil)
}

func (t *BillChainCode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "issue" {
		return t.issue(stub, args)
	} else if function == "queryMyBills" {
		return t.queryMyBills(stub, args)
	} else if function == "queryBillByNo" {
		return t.queryBillByNo(stub, args)
	} else if function == "queryMyWaitBills" {
		return t.queryMyWaitBills(stub, args)
	} else if function == "endorse" {
		return t.endorse(stub, args)
	} else if function == "accept" {
		return t.accept(stub, args)
	} else if function == "reject" {
		return t.reject(stub, args)
	}
	return shim.Error("指定的函数名称错误")
}

// 这里是智能合约的入口
func main()  {
	err := shim.Start(new(BillChainCode))
	if err != nil {
		fmt.Println("启动链码错误: ", err)
	}
}
