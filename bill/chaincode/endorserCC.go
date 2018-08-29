package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
)

// 票据背书请求
// args: 0 - Bill_No; 1 - endorseCmId(待背书人ID); 2 - endorseAcct(待背书人名称)
func (t *BillChainCode) endorse(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// 1. 检查参数长度是否为3(票据号码, 待背书人ID, 待背书人名称)
	if len(args) != 3 {
		res := GetRetString(1, "必须指定票据号码，待背书人证件号码及待背书人名称")
		return shim.Error(res)
	}

	// 2. 根据票据号码获取票据状态
	bill, b1 := t.getBill(stub, args[0])
	if !b1 {
		res := GetRetString(1, "根据指定的票据号码查询信息时发生错误")
		return shim.Error(res)
	}

	// 3. 检查待背书人与当前持票人是否为同一人
	if bill.HoldrCmID == args[1] {
		res := GetRetString(1, "被背书人不能是当前持票人")
		return shim.Error(res)
	}

	// 4. 获取票据历史变更数据
	iterator, err := stub.GetHistoryForKey(bill.BillInfoID)
	if err != nil {
		res := GetRetString(1, "获取票据流转历史信息时发生错误")
		return shim.Error(res)
	}
	defer iterator.Close()

	// 5. 检查待背书人是否为票据历史持有人
	var hisBill Bill
	for iterator.HasNext() {
		hisData, err := iterator.Next()
		if err != nil {
			res := GetRetString(1, "获取历史数据时发生错误")
			return shim.Error(res)
		}
		json.Unmarshal(hisData.Value, &hisBill)
		if hisData.Value == nil {
			var empty Bill
			hisBill = empty
		}

		if bill.HoldrCmID == args[1] {
			res := GetRetString(1, "被背书人不能是该票据的历史持有人")
			return shim.Error(res)
		}
	}

	// 6. 更改票据信息与状态: 添加待背书人信息(证件号码与名称), 票据状态 更改为待背书, 重置已拒绝背书人
	bill.State = BillInfo_State_EndorseWaitSign
	bill.WaitEndorseCmID = args[1]
	bill.WaitEndorseAcct = args[2]
	bill.RejectEndorseAcct = ""
	bill.RejectEndorseCmID = ""

	// 7. 保存票据
	_, b1 = t.putBill(stub, bill)
	if !b1 {
		res := GetRetString(1, "票据背书请求失败，保存票据信息时发生错误")
		return shim.Error(res)
	}

	// 8. 保存以待背书人ID与票据号码构造的复合键, 以便待背书人批量查询. value为空即可
	waitEndorserCmIdBillInfoIndexKey, err := stub.CreateCompositeKey(IndexName,
		[]string{bill.HoldrCmID, bill.BillInfoID})
	if err != nil {
		res := GetRetString(1, "根据待背书人的证件号码及票据号码创建复合键失败")
		return shim.Error(res)
	}

	stub.PutState(waitEndorserCmIdBillInfoIndexKey, nil)

	// 9. 返回
	return shim.Success([]byte("发送请求成功，此票据待被背书人处理"))
}

// 票据背书签收
// args: 0 - Bill_No; 1 - endorseCmId(待背书人ID); 2 - endorseAcct(待背书人名称)
func (t *BillChainCode) accept(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// 1. 检查参数长度是否为3(票据号码, 待背书人ID, 待背书人名称)
	if len(args) != 3 {
		res := GetRetString(1, "必须指定票据号码，待背书人证件号码及待背书人名称")
		return shim.Error(res)
	}

	// 2. 根据票据号码获取票据状态
	bill, b1 := t.getBill(stub, args[0])
	if !b1 {
		res := GetRetString(1, "根据指定的票据号码查询信息时发生错误")
		return shim.Error(res)
	}

	// 3. 以前手持票人ID与票据号码构造复合键, 删除该key, 以便前手持票人无 法再查到该票据
	holderCmIdBillInfoIdIndexKey, err := stub.CreateCompositeKey(IndexName,
		[]string{bill.HoldrCmID, bill.BillInfoID})
	if err != nil {
		res := GetRetString(1, "创建复合键失败")
		return shim.Error(res)
	}
	stub.DelState(holderCmIdBillInfoIdIndexKey)

	// 4. 更改票据信息与状态: 将当前持票人更改为待背书人(证件与名称), 票 据状态更改为背书签收, 重置待背书人
	bill.State = BillInfo_State_EndorseSigned
	bill.WaitEndorseCmID = args[1]
	bill.WaitEndorseAcct = args[2]
	bill.RejectEndorseAcct = ""
	bill.RejectEndorseCmID = ""

	// 5. 保存票据
	_, b1 = t.putBill(stub, bill)
	if !b1 {
		res := GetRetString(1, "票据背书签收失败，保存票据信息时发生错误")
		return shim.Error(res)
	}

	// 6.返回
	return shim.Success([]byte("票据背书签收成功"))
}

// 票据背书拒签(拒绝背书)
// args: 0 - bill_NO; 1 - endorseCmId(待背书人ID); 2 - endorseAcct(待背书人名称)
func (t *BillChainCode) reject(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// 1. 检查参数长度是否为3(票据号码, 待背书人ID, 待背书人名称)
	if len(args) != 3 {
		res := GetRetString(1, "必须指定票据号码，待背书人证件号码及待背书人名称")
		return shim.Error(res)
	}

	// 2. 根据票据号码查询对应的票据状态
	bill, b1 := t.getBill(stub, args[0])
	if !b1 {
		res := GetRetString(1, "根据指定的票据号码查询信息时发生错误")
		return shim.Error(res)
	}

	// 3. 以待背书人ID及票据号码构造复合键, 从search中删除该key, 以便当 前被背书人无法再次查询到该票据
	holderCmIdBillInfoIdIndexKey, err := stub.CreateCompositeKey(IndexName,
		[]string{args[1], bill.BillInfoID})
	if err != nil {
		res := GetRetString(1, "创建复合键失败")
		return shim.Error(res)
	}
	stub.DelState(holderCmIdBillInfoIdIndexKey)

	// 4. 更改票据信息与状态: 将拒绝背书人更改为当前待背书人(证件号码与名 称), 票据状态更改为背书拒绝, 重置待背书人
	bill.State = BillInfo_State_EndorseReject
	bill.RejectEndorseAcct = args[1]
	bill.RejectEndorseCmID = args[2]
	bill.WaitEndorseCmID = ""
	bill.WaitEndorseAcct = ""

	// 5. 保存票据状态
	_, b1 = t.putBill(stub, bill)
	if !b1 {
		res := GetRetString(1, "票据背书拒签失败，保存票据信息时发生错误")
		return shim.Error(res)
	}

	// 6. 返回
	return shim.Success([]byte("票据背书拒签成功"))
}