package blockchain


type Bill struct {
	BillInfoID string `json:"BillInfoID"` //票据号码
	BillInfoAmt string `json:"BillInfoAmt"` //票据金额
	BillInfoType string `json:"BillInfoType"` //票据类型
	BillInfoIsseDate string `json:"BillInfoIsseDate"` // 票据出票日期
	BillInfoDueDate string `json:"billInfoDueDate"` // 票据到期日期
	DrwrCmID string `json:"DrwrCmID"` // 出票人证件号码
	DrwrAcct string `json:"DrwrAcct"` // 出票人名称
	AccptrCmID string `json:"AccptrCmID"` // 承兑人证件号码
	AccptrAcct string `json:"AccptrAcct"` // 承兑人名称
	PyeeCmID string `json:"PyeeCmID"` // 收款人证件号码
	PyeeAcct string `json:"PyeeAcct"` // 收款人名称
	HoldrCmID string `json:"HodrCmID"` // 当前持票人证件号码
	HoldrAcct string `json:"HodrAcct` // 当前持票人名称
	WaitEndorseCmID string `json:"WaitEndorseCmID"` // 待背书人证件号码
	WaitEndorseAcct string `json:"WaitEndorseAcct` // 待背书人名称
	RejectEndorseCmID string `json:"WaitEndorseCmID"` // 拒绝背书人证件号码
	RejectEndorseAcct string `json:"WaitEndorseAcct"` // 拒绝背书人名称
	State string `json:"State"` // 票据状态
	History []HistoryItem  `json:"History"` // 票据背书历史
}

// 票据历史信息
type HistoryItem struct {
	TxId        string      `json:"TxId"`
	Bill        Bill        `json:"bill"`
}