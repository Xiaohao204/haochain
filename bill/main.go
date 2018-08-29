package main

import (
	"fmt"
	"github.com/haochain/bill/blockchain"
	"os"
	"github.com/haochain/bill/web/controller"
	"encoding/json"
	"strconv"
	"github.com/haochain/bill/web"
)

func main() {
	// Definition of the Fabric SDK properties
	fSetup := blockchain.FabricSetup{
		// Network parameters
		OrdererID: "orderer.hf.haochain.io",

		// Channel parameters
		ChannelID:     "haochain",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/haochain/bill/fixtures/artifacts/haochain.channel.tx",

		// Chaincode parameters
		ChainCodeID:     "hao-service",
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/haochain/bill/chaincode/",
		OrgAdmin:        "Admin",
		OrgName:         "org1",
		ConfigFile:      "config.yaml",

		// User parameters
		UserName: "User1",
	}

	// Initialization of the Fabric SDK from the previously set properties
	err := fSetup.Initialize()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
		return
	}

	// Install and instantiate the chaincode
	err = fSetup.InstallAndInstantiateCC()
	if err != nil {
		fmt.Printf("Unable to install and instantiate the chaincode: %v\n", err)
		return
	}

	// ==========================测试开始 ==============================
	// 发布票据
	bill := blockchain.Bill{
		BillInfoID:       "BOC10000001",
		BillInfoAmt:      "222",
		BillInfoType:     "111",
		BillInfoIsseDate: "20180826",
		BillInfoDueDate:  "20180901",
		DrwrCmID:         "111",
		DrwrAcct:         "111",
		AccptrCmID:       "111",
		AccptrAcct:       "111",
		PyeeCmID:         "111",
		PyeeAcct:         "111",
		HoldrCmID:        "BCMID",
		HoldrAcct:        "B公司 ",
	}

	// 发布票据
	resp, err := fSetup.SaveBill(bill)
	if err != nil {
		fmt.Printf("发布票据失败: %v\n", err)
	} else {
		fmt.Println("发布票据成功 =========> " + resp)
	}

	// 根据当前用户的证件号码查询票据列表
	b, err := fSetup.QueryBills("BCMID")
	// 查看返回结果
	if err != nil {
		fmt.Errorf(err.Error())
	} else {
		fmt.Println("查询用户票据成功")
		var bills = []blockchain.Bill{}
		json.Unmarshal(b, &bills)

		fmt.Println("bill length = " + strconv.Itoa(len(bills)))

		for _, temp := range bills {
			fmt.Println(temp)
		}
	}

	// 根据票据号码获取票据状态及该票据的背书历史
	//b, err = fSetup.FindBillByNo("BOC10000001")
	//
	//// 查看返回结果
	//if err != nil {
	//	fmt.Errorf(err.Error())
	//} else {
	//	fmt.Println("根据票据号码查询票据详情成功")
	//	bill = blockchain.Bill{}
	//	json.Unmarshal(b, &bill)
	//
	//	for _, hisItem := range bill.History {
	//		fmt.Println(hisItem)
	//	}
	//}

	//// 票据背书请求
	//resp, err = fSetup.Endorse("BOC10000001", "CCMID", "C公司")
	//// 查看返回结果
	//if err != nil {
	//	fmt.Errorf(err.Error())
	//} else {
	//	fmt.Println("票据背书成功")
	//	fmt.Println(resp)
	//}
	//
	//// 根据待背书人证件号码, 查询当前用户的待背书票据
	//b, err = fSetup.FindWaitBills("CCMID")
	//// 查看返回结果
	//if err != nil {
	//	fmt.Errorf(err.Error())
	//} else {
	//	fmt.Println("根据待背书人证件号码查询其对应的待背书票据成功")
	//	var bills = []service.Bill{}
	//	json.Unmarshal(b, &bills)
	//
	//	for _, temp := range bills {
	//		fmt.Println(temp)
	//	}
	//}
	//
	//// 票据背书签收
	//resp, err = fSetup.EndorseAccept("BOC10000001", "CCMID", "C公司")
	//// 查看返回结果
	//if err != nil {
	//	fmt.Errorf(err.Error())
	//} else {
	//	fmt.Println("票据背书签收成功")
	//	fmt.Println(resp)
	//}
	//
	//// 票据背书拒签
	//resp, err = fSetup.EndorseReject("BOC10000001", "CCMID", "C公司")
	//// 查看返回结果
	//if err != nil {
	//	fmt.Errorf(err.Error())
	//} else {
	//	fmt.Println("票据背书拒签成功")
	//	fmt.Println(resp)
	//}

	app := controller.Application{
		Fabric: &fSetup,
	}

	web.WebStart(&app)

}