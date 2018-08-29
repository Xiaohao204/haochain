package controller

import (
	"net/http"
	"github.com/haochain/bill/blockchain"
	"fmt"
	"encoding/json"
)

var cuser User

type Application struct {
	Fabric *blockchain.FabricSetup
}

func (app *Application) LoginView(w http.ResponseWriter, r *http.Request){
	fmt.Println("开始进入登陆页面")
	response(w, r, "login.html", nil)
}


func (app *Application) IssueView(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("开始进入发布票据页面")
	data := &struct {
		Flag bool
		Msg string
		Cuser User
	}{
		Flag: false,
		Msg: "",
		Cuser: cuser,
	}
	response(w, r, "issue.html", data)
}

func (app *Application) Loginout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("登出")
	cuser = User{}
	app.LoginView(w, r)
}

func (app *Application) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("登陆")
	// 接收请求参数
	userName := r.FormValue("userName")
	password := r.FormValue("password")
	// 封装响应数据
	var flag= false
	for _, user := range Users {
		// 验证请求数据的正确性
		if userName == user.UserName && password == user.Password {
			cuser = user
			flag = true
			break
		}

	}

	if flag {
		app.FindBills(w, r)
	} else {
		data := &struct {
			Name string
			Flag bool
		}{
			Name: userName,
			Flag: true,
		}
		response(w, r, "login.html", data)
	}

}

// 查询我的票据列表
func (app *Application) FindBills(w http.ResponseWriter, r *http.Request) {
	fmt.Println("FindBills")
	//holdeCmId := r.FormValue("holdeCmId")
	//holdeCmId := cuser.CmId
	// 调用业务层
	result, err := app.Fabric.QueryBills(cuser.CmId)

	if err != nil{
		fmt.Println("查询当前用户的票据列表失败: ", err.Error())
	}

	// 反序列化并封装响应数据
	var bills []blockchain.Bill
	json.Unmarshal(result, &bills)

	data := &struct {
		Bills []blockchain.Bill
		Cuser User
	}{
		Bills: bills,
		Cuser: cuser,
	}

	//响应客户端
	response(w, r, "bills.html", data)
}

// 发布票据
func (app *Application) Issue(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("Issue")
	bill := blockchain.Bill{
		BillInfoID: r.FormValue("BillInfoID"),
		BillInfoType: r.FormValue("BillInfoType"),
		BillInfoAmt: r.FormValue("BillInfoAmt"),

		DrwrAcct: r.FormValue("DrwrAcct"),
		DrwrCmID: r.FormValue("DrwrCmID"),

		AccptrAcct: r.FormValue("AccptrAcct"),
		AccptrCmID: r.FormValue("AccptrCmID"),

		PyeeAcct: r.FormValue("PyeeAcct"),
		PyeeCmID: r.FormValue("PyeeCmID"),

		HoldrAcct: r.FormValue("HoldrAcct"),
		HoldrCmID: r.FormValue("HoldrCmID"),
	}

	transactionId, err := app.Fabric.SaveBill(bill)
	var msg string
	var flag bool
	if err != nil {
		msg = "票据发布失败： " + err.Error()
		flag = false
	} else {
		msg = "票据发布成功： " + transactionId
		flag = true
	}

	data := &struct {
		Flag bool
		Msg string
		Cuser User
	}{
		Flag: flag,
		Msg: msg,
		Cuser: cuser,
	}

	//响应客户端
	response(w, r, "issue.html", data)

}

// 查询票据详情
func (app *Application) QueryBillInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("QueryBillInfo")
	// 获取提交数据
	billInfoNo := r.FormValue("billInfoNo")

	result, err := app.Fabric.FindBillByNo(billInfoNo)
	if err != nil {
		fmt.Println(err.Error())
	}

	var bill blockchain.Bill
	json.Unmarshal(result, &bill)

	data := &struct {
		Bill blockchain.Bill
		Cuser User
		Flag bool
		Msg string
	}{
		Bill: bill,
		Cuser: cuser,
		Flag: false,
		Msg: "",
	}

	flag := r.FormValue("flag")
	if flag == "t" {
		data.Flag = true
		data.Msg = r.FormValue("Msg")
	}

	//响应客户端
	response(w, r, "billInfo.html", data)

}

// 票据背书请求
func (app *Application) Endorse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endorse")
	waitEndorseAcct := r.FormValue("waitEndorseAcct")
	waitEndorseCmId := r.FormValue("waitEndorseCmId")
	billNo := r.FormValue("billNo")

	result, err := app.Fabric.Endorse(billNo, waitEndorseCmId, waitEndorseAcct)
	if err != nil {
		fmt.Println(err.Error())
	}

	r.Form.Set("billInfoNo", billNo)
	r.Form.Set("flag", "t")
	r.Form.Set("Msg", result)

	app.QueryBillInfo(w, r)

	//响应客户端
	//response(w, r, "billInfo.html", data)

}

// 查询待背书票据列表
func (app *Application) WaitEndorBills(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("WaitEndorBills")
	waitEndorseCmId := cuser.CmId
	result, err := app.Fabric.FindWaitBills(waitEndorseCmId)
	if err != nil {
		fmt.Println(err.Error())
	}

	var bills []blockchain.Bill
	json.Unmarshal(result, &bills)

	data := &struct {
		Bills []blockchain.Bill
		Cuser User
	}{
		Bills: bills,
		Cuser: cuser,
	}

	response(w, r, "waitEndorses.html", data)

}

// 查询待背书票据详情
func (app *Application) WaitEndorseInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("WaitEndorseInfo")
	billNo := r.FormValue("billNo")
	result, err := app.Fabric.FindBillByNo(billNo)
	if err != nil {
		fmt.Println(err.Error())
	}

	var bill blockchain.Bill
	json.Unmarshal(result, &bill)

	data := &struct {
		Bill blockchain.Bill
		Cuser User
		Flag bool
		Msg string
	}{
		Bill: bill,
		Cuser: cuser,
		Flag: false,
		Msg: "",
	}

	flag := r.FormValue("flag")
	if flag == "t" {
		data.Flag = true
		data.Msg = r.FormValue("Msg")
	}

	response(w, r, "waitEndorseInfo.html", data)

}

// 签收
func (app *Application) Accept(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Accept")
	billNo := r.FormValue("billNo")
	cmid := cuser.CmId
	acct := cuser.Acct

	result, err := app.Fabric.EndorseAccept(billNo, cmid, acct)
	if err != nil {
		fmt.Println(err.Error())
	}

	r.Form.Set("billNo", billNo)
	r.Form.Set("flag", "t")
	r.Form.Set("Msg", result)

	app.WaitEndorseInfo(w, r)

}

// 拒签
func (app *Application) Reject(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Reject")
	billNo := r.FormValue("billNo")
	cmid := cuser.CmId
	acct := cuser.Acct

	result, err := app.Fabric.EndorseReject(billNo, cmid, acct)
	if err != nil {
		fmt.Println(err.Error())
	}

	r.Form.Set("billNo", billNo)
	r.Form.Set("flag", "t")
	r.Form.Set("Msg", result)

	app.WaitEndorseInfo(w, r)

}