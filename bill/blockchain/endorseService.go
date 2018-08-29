package blockchain

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"fmt"
)

// 票据背书请求
// args: 0 - Bill_No; 1 - endorserCmId(待背书人ID); 2 - endorserAcct(待背书人名称)
func (setup *FabricSetup) Endorse(billNo string, endorseCmId string, endorseAcct string) (string, error) {
	//var args []string
	// 1. 将所需参数添加至args中
	var args []string
	args = append(args, "endorse")
	args = append(args, billNo)
	args = append(args, endorseCmId)
	args = append(args, endorseAcct)


	// 2. 设置查询的请求参数,执行查询
	response, err := setup.client.Query(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3])}})
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	// 3. 返回
	b := response.Payload
	return string(b[:]), nil
}

// 根据待背书人证件号码, 查询当前用户的待背书票据
func (setup *FabricSetup) FindWaitBills(endorserCmId string) ([]byte, error) {
	//var args []string
	// 1. 将所需参数添加至args中
	var args []string
	args = append(args, "queryMyWaitBills")
	args = append(args, endorserCmId)


	// 2. 设置查询的请求参数,执行查询
	response, err := setup.client.Query(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1])}})
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	// 3. 返回
	b := response.Payload
	return b[:], nil
}

// 票据背书签收
// args: 0 - Bill_No; 1 - endorserCmId(待背书人ID); 2 - endorserAcct(待背书人名称)
func (setup *FabricSetup) EndorseAccept(billNo string, endorseCmId string, endorseAcct string) (string, error) {
	//var args []string
	// 1. 将所需参数添加至args中
	var args []string
	args = append(args, "accept")
	args = append(args, billNo)
	args = append(args, endorseCmId)
	args = append(args, endorseAcct)


	// 2. 设置查询的请求参数,执行查询
	response, err := setup.client.Query(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3])}})
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	// 3. 返回
	b := response.Payload
	return string(b[:]), nil
}

// 票据背书拒签(拒绝背书)
// args: 0 - bill_NO; 1 - endorserCmId(待背书人ID); 2 - endorserAcct(待背书人名称)
func (setup *FabricSetup) EndorseReject(billNo string, endorseCmId string, endorseAcct string) (string, error) {
	//var args []string
	// 1. 将所需参数添加至args中
	var args []string
	args = append(args, "reject")
	args = append(args, billNo)
	args = append(args, endorseCmId)
	args = append(args, endorseAcct)


	// 2. 设置查询的请求参数,执行查询
	response, err := setup.client.Query(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3])}})
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	// 3. 返回
	b := response.Payload
	return string(b[:]), nil
}