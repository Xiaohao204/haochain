package controller

type User struct {
	UserName    string  `json:"UserName"`
	Name        string  `json:"Name"`
	Password    string  `json:"Password"`
	CmId        string  `json:"CmId"`
	Acct        string  `json:"Acct"`
}
var Users []User

func init() {
	admin := User{UserName: "admin", Name: "管理员", Password:
		"123456", CmId: "HODR01", Acct: "管理员"}

	alice := User{UserName: "alice", Name: "A公司", Password:
		"123456", CmId: "ACMID", Acct: "A公司"}

	bob := User{UserName: "bob", Name: "B公司", Password: "123456",
		CmId: "BCMID", Acct: "B公司"}

	jack := User{UserName: "jack", Name: "C公司", Password:
		"123456", CmId: "CCMID", Acct: "C公司"}

	Users = append(Users, admin)
	Users = append(Users, alice)
	Users = append(Users, bob)
	Users = append(Users, jack)
}
