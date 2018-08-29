package main

import (
	"encoding/json"
	"fmt"
)

// chaincode Response结构
type chaincodeRet struct {
	Code int // 0代表成功, 否则为1
	Dec string // 描述
}

// 根据返回码和描述信息返回序列化后的字节数组
func GetRetByte(code int, dec string) ([]byte) {
	b, err := getRet(code, dec)
	if err != nil {
		return []byte{}
	}
	return b
}

// 根据返回码和描述信息返回序列化后的字符串
func GetRetString(code int, dec string) (string) {
	b, err := getRet(code, dec)
	if err != nil {
		return "" }
	return string(b[:])
}

// 根据返回码和描述信息进行序列化
func getRet(code int, dec string) ([]byte, error) {
	var c chaincodeRet
	c.Code = code
	c.Dec = dec
	b, err := json.Marshal(c)
	if err != nil {
		return []byte{0x00}, fmt.Errorf("序列化失败")
	}
	return b, err
}
