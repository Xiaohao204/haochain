package blockchain

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"encoding/json"
	"fmt"
)

// 发布票据
func (setup *FabricSetup) SaveBill(bill Bill) (string, error) {
	var args []string
	args = append(args, "issue")
	b, _ := json.Marshal(bill)

	response, err := setup.client.Query(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{b}})
	if err != nil {
		return "", fmt.Errorf("failed to query: %v", err)
	}
	//// 设置交易请求参数
	//req := chclient.Request{ChaincodeID: setup.Fabric.ChaincodeID,
	//	Fcn: args[0], Args: [][]byte{b}}
	//// 执行交易
	//response, err := setup.Fabric.Client.Execute(req) if err != nil {
	//	returnn "", fmt.Errorf("保存票据时发生错误: %v\n", err) }
	return string(response.Payload), nil
	//return response.TransactionID.ID, nil
}

// 根据当前持票人证件号码, 批量查询票据
func (setup *FabricSetup) QueryBills(holderCmId string) ([]byte, error) {
	//var args []string
	// 1. 将所需参数添加至args中
	var args []string
	args = append(args, "queryMyBills")
	args = append(args, holderCmId)

	// 2. 设置查询的请求参数,执行查询
	response, err := setup.client.Query(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(holderCmId)}})
	if err != nil {
		return nil, fmt.Errorf("根据持票人的证件号码查询票据失败: %v", err)
	}
	// 3. 返回
	b := response.Payload
	return b[:], nil
}

// 根据票据号码获取票据状态及该票据的背书历史
func (setup *FabricSetup) FindBillByNo(billNo string) ([]byte, error) {
	//var args []string
	// 1. 将所需参数添加至args中
	var args []string
	args = append(args, "queryBillByNo")
	args = append(args, billNo)

	// 2. 设置查询的请求参数,执行查询
	response, err := setup.client.Query(channel.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1])}})
	if err != nil {
		return nil, fmt.Errorf("根据billNo查询票据失败: %v", err)
	}

	// 3. 返回
	b := response.Payload
	return b[:], nil
}

