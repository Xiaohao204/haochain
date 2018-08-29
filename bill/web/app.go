package web

import (
	"github.com/haochain/bill/web/controller"
	"net/http"
	"fmt"
)

func WebStart(app *controller.Application) error {
	// 指定文件服务器,将当前目录作为项目的根目录
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	//http.HandleFunc("/", app.LoginView)
	http.HandleFunc("/login.html", app.LoginView)
	http.HandleFunc("/issue.html", app.IssueView)
	http.HandleFunc("/loginout", app.Loginout)
	http.HandleFunc("/userLogin", app.Login)
	http.HandleFunc("/myBills", app.FindBills)
	http.HandleFunc("/billInfo", app.QueryBillInfo)
	http.HandleFunc("/endorse", app.Endorse)
	http.HandleFunc("/issue", app.Issue)
	http.HandleFunc("/waitEndorBills", app.WaitEndorBills)
	http.HandleFunc("/waitEndorseInfo", app.WaitEndorseInfo)
	http.HandleFunc("/accept", app.Accept)
	http.HandleFunc("/reject", app.Reject)


	fmt.Println("启动Web服务, 监听端口号: 9000")

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println("启动Web服务错误")
	}

	return nil

}
