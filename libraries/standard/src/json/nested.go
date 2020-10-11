package main

import (
	"encoding/json"
	"fmt"
)

type Feedback struct {
	Msg string `json:"msg"`
}

type FeedbackMsgModel struct {
	MsgType    string      `json:"xxlx"`
	MsgContent interface{} `json:"xxnr"`
}

type FeedbackReportPerson struct {
	SerialNumber string           `json:"bbywlsh"` // 报备业务流水号
	Result       []FeedbackPerson `json:"bbjg"`    // 报备结果
}

type FeedbackCertInfo struct {
	CertType   string `json:"czlx"` // 证件类型
	CertNumber string `json:"zjhm"` // 证件号码
}

type FeedbackPerson struct {
	PersonIdent string             `json:"rybs"` // 人员标识
	State       string             `json:"bbzt"` // 报备状态
	CertList    []FeedbackCertInfo `json:"zjlb"` // 证件列表
}

func main() {
	msg := FeedbackMsgModel{
		MsgType: "type",
		MsgContent: FeedbackReportPerson{
			SerialNumber: "1",
			Result: []FeedbackPerson{{
				PersonIdent: "1",
				State:       "2",
				CertList: []FeedbackCertInfo{{
					CertType:   "1",
					CertNumber: "2",
				}}},
			},
		},
	}
	buf, _ := json.Marshal(msg)
	//fmt.Printf("%s", buf)

	feedback := Feedback{
		Msg: string(buf),// 因为是 interface{} 无法反序列化回来
	}
	var msgBuf FeedbackMsgModel
	_ = json.Unmarshal([]byte(feedback.Msg), &msgBuf)
	fmt.Printf("%+v",msgBuf)
}
