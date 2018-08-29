package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
)

// 根据指定的票据号码获取对应的票据信息
func (t *BillChainCode) getBill(stub shim.ChaincodeStubInterface, bill_No string) (Bill, bool) {
	var bill Bill
	// 1. 将定义的前缀与票据号码连接成一个key
	key := Bill_Prefix + bill_No

	// 2. 以Key为参数查询票据信息
	b, err := stub.GetState(key)
	if b == nil || err != nil {
		return bill, false
	}

	// 3. 反序列化查询到的票据
	err = json.Unmarshal(b, &bill)
	if err != nil {
		return bill, false
	}

	// 4. 返回
	return bill, true
}

// 保存指定的票据信息
func (t *BillChainCode) putBill(stub shim.ChaincodeStubInterface, bill Bill) ([]byte, bool) {
	// 1. 将参数bill序列化
	b, err := json.Marshal(bill)
	if err != nil {
		return nil, false
	}

	// 2. 以定义的前缀与票据号码连接形成一个Key将bill保存在账本中
	err = stub.PutState(Bill_Prefix + bill.BillInfoID, b)
	if err != nil {
		return nil, false
	}

	// 3. 返回
	return b, true
}

// 发布票据
// args: 0 - {bill object}
func (t *BillChainCode) issue(stub shim.ChaincodeStubInterface, args[] string) peer.Response {
	// 1. 检查参数个数是否为一个
	if len(args) != 1 {
		res := GetRetString(1, "发布票据指定的参数个数只能是一个")
		return shim.Error(res)
	}

	// 2. 反序列化args[0]参数为bill
	var bill Bill
	err := json.Unmarshal([]byte(args[0]), &bill)
	if err != nil {
		res := GetRetString(1, "发布票据时指定的票据信息反序列化失败")
		return shim.Error(res)
	}

	// 3. 查重(根据票据号码查询是否存在. 如果根据票据号码查询到信息, 证明 已存在)
	_, exist := t.getBill(stub, bill.BillInfoID)
	if exist {
		res := GetRetString(1, "发布票据失败，该票据号码已存在：" + bill.BillInfoID)
		return shim.Error(res)
	}

	// 4. 更改票据状态(票据状态设为新发布). 并保存票据:
	bill.State = BillInfo_State_NewPublish

	_, b1 := t.putBill(stub, bill)
	if !b1 {
		res := GetRetString(1, "发布票据失败，保存票据信息时发生错误！")
		return shim.Error(res)
	}

	// 5. 保存以当前持票人ID与票据号码构造的复合键, 以便持票人批量查询. value为空即可
	holderNameBillNoIndexKey, err := stub.CreateCompositeKey(IndexName,
		[]string{bill.HoldrCmID, bill.BillInfoID})
	if err != nil {
		res := GetRetString(1, "创建持票人ID与票据号码的复合键失败")
		return shim.Error(res)
	}

	stub.PutState(holderNameBillNoIndexKey, []byte{0x00})

	// 6. 返回
	res := GetRetByte(0, "票据发布成功")
	return  shim.Success(res)
}

// 查询当前用户的票据(根据当前持票人证件号码, 批量查询票据)
// args: 0 - holderCmId(当前持票人证件号码)
func (t *BillChainCode) queryMyBills(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// 1. 检查参数个数是否为一个
	if len(args) != 1 {
		res := GetRetString(1, "发布票据指定的参数个数只能是一个")
		return shim.Error(res)
	}

	// 根据指定的组合键查询分类账上的状态
	// -- 在发布票据方法中实现了将复合键(持票人ID与票据号码构造的复合键)保 存在分类账中
	// 2. 根据当前持票人证件号码从search中查询所持有的票号
	billsIterator, err := stub.GetStateByPartialCompositeKey(IndexName, []string{args[0]})
	if err != nil {
		res := GetRetString(1, "查询票据失败，查询所持有的票号时发生错误")
		return shim.Error(res)
	}
	defer billsIterator.Close()

	// 3. 迭代处理
	var bills = []Bill{}
	for billsIterator.HasNext() {
		kv, _ := billsIterator.Next()
		// 获取持票人名下的票号
		_, compositeKeyParts, err := stub.SplitCompositeKey(kv.Key)
		if err != nil {
			res := GetRetString(1, "查询票据失败，分割复合键时发生错误")
			return shim.Error(res)
		}

		// 根据获取到的票据号码查询对应的票据状态
		bill, b1 := t.getBill(stub, compositeKeyParts[1])
		if !b1 {
			res := GetRetString(1, "根据票据号码查询对应的票据状态时发生错误")
			return shim.Error(res)
		}

		bills = append(bills, bill)
	}

	// 4. 序列化票据数组
	b, err := json.Marshal(bills)
	if err != nil {
		res := GetRetString(1, "查询票据失败，序列化票据状态时发生错误")
		return shim.Error(res)
	}

	// 5. 返回查询结果
	return shim.Success(b)
}

// 根据票据号码获取票据状态及该票据的背书历史
// args: 0 - bill_No
func (t *BillChainCode) queryBillByNo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// 1. 检查参数个数是否为一个
	if len(args) != 1 {
		res := GetRetString(1, "发布票据指定的参数个数只能是一个")
		return shim.Error(res)
	}

	// 2. 根据票据号码获取对应的票据状态
	bill, b1 := t.getBill(stub, args[0])
	if !b1 {
		res := GetRetString(1, "获取票据状态及背书历史失败，根据给定的票据号码查询对应的票据状态时发生错误")
		return shim.Error(res)
	}

	// 3. 获取票据背书变更历史
	billIterator, err := stub.GetHistoryForKey(Bill_Prefix + args[0])
	if err != nil {
		res := GetRetString(1, "获取票据状态及背书历史失败，查询背书历史变更时发生错误")
		return shim.Error(res)
	}

	// 4. 迭代处理
	var bills []HistoryItem
	var hisBill Bill

	for billIterator.HasNext() {
		hisData, err := billIterator.Next()
		if err != nil {
			res := GetRetString(1, "获取历史流转信息时发生错误")
			return shim.Error(res)
		}

		var historyItem HistoryItem
		historyItem.TxId = hisData.TxId
		json.Unmarshal(hisData.Value, &hisBill)
		if hisData.Value == nil {
			var empty Bill
			historyItem.Bill = empty
		} else {
			historyItem.Bill = bill
		}
		bills = append(bills, historyItem)
	}

	// 5. 将背书历史做为票据的属性返回
	bill.History = bills

	// 6. 序列化bill
	b, err := json.Marshal(bill)
	if err != nil {
		res := GetRetString(1, "序列化票据时发生错误")
		return shim.Error(res)
	}

	// 7. 返回结果
	return shim.Success(b)
}

// 查询当前用户的待背书票据(根据待背书人证件号码)
// args: 0 - endorserCmId
func (t *BillChainCode) queryMyWaitBills(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// 1. 检查参数个数是否为一个
	if len(args) != 1 {
		res := GetRetString(1, "发布票据指定的参数个数只能是一个")
		return shim.Error(res)
	}

	// 2. 根据待背书人证件号码从search中查询待背书的票据号码
	// -- 票据背书请求方法中实现了将待背书人ID与票据号码构造的复合键信息的 保存
	billIterator, err := stub.GetStateByPartialCompositeKey(IndexName, []string{args[0]})
	if err != nil {
		res := GetRetString(1, "查询待背书票据失败，查询待背书人的票据号码时发生错误")
		return shim.Error(res)
	}
	defer billIterator.Close()

	// 3. 迭代处理
	var bills = []Bill{}
	for billIterator.HasNext() {

		kv, _ := billIterator.Next()
		// 获取持票人名下的票号
		_, compositeKeyParts, err := stub.SplitCompositeKey(kv.Key)

		if err != nil {
			res := GetRetString(1, "分割复合key时发生错误")
			return shim.Error(res)
		}

		// 根据票据号码查询相应的票据状态
		bill, b1 := t.getBill(stub, compositeKeyParts[1])
		if !b1 {
			res := GetRetString(1, "查询待背书票据失败，根据待背书票据号码查询对应票据状态时发生错误")
			return shim.Error(res)
		}

		if bill.State == BillInfo_State_EndorseWaitSign && bill.WaitEndorseCmID == args[0] {
			bills = append(bills, bill)
		}

	}

	// 4. 序列化待背书票据数组
	b, err := json.Marshal(bills)
	if err != nil {
		res := GetRetString(1, "查询待背书票据失败，序列化待背书票据数组时发生错误")
		return shim.Error(res)
	}

	// 5. 返回结果
	return shim.Success(b)
}