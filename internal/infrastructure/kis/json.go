package kis

import (
	"encoding/json"
)


const (
	custtype            string = "P"
	tr_type_subscribe   string = "1"
	tr_type_unsubscribe string = "0"
)

type getApprovalKeyReqeust struct {
	GrantType string `json:"grant_type"`
	AppKey    string `json:"appkey"`
	SecretKey string `json:"secretkey"`
}

type getApprovalKeyResponse struct {
	ApprovalKey string `json:"approval_key"`
}

type HeaderJson struct {
	ApprovalKey string `json:"approval_key"` // 실시간 접속키
	Custtype    string `json:"custtype"`     // 고객 타입 (P: 개인, B: 법인)
	TrType      string `json:"tr_type"`      // 거래 타입 (1. 등록, 2. 해제)
	ContentType string `json:"content-type"` // 컨텐츠 타입 (utf-8 고정)
}

type RequestBodyJson struct {
	Input RequestInputJson `json:"input"`
}

type RequestInputJson struct {
	TrId  string `json:"tr_id"`  // 거래 ID (H0STCNT0: 실시간 주식 체결가, H0STASP0: 주식 호가, HDFSCNT0: 실시간 미국장)
	TrKey string `json:"tr_key"` // 종목코드
}

type RequestJson struct {
	Header HeaderJson      `json:"header"`
	Body   RequestBodyJson `json:"body"`
}

type ResponseJson struct {
	Header ResponseHeaderJson `json:"header"`
	Body   ResponseBodyJson   `json:"body"`
}

type ResponseHeaderJson struct {
	TrID    string `json:"tr_id"`
	TrKey   string `json:"tr_key"`
	Encrypt string `json:"encrypt"`
}

type ResponseBodyJson struct {
	RtCd   string `json:"rt_cd"`
	MsgCd  string `json:"msg_cd"`
	Msg1   string `json:"msg1"`
	Output struct {
		Iv  string `json:"iv"`
		Key string `json:"key"`
	} `json:"output"`
}



func parseToSubscriptionResponse(data []byte) (string, bool) {

	var res ResponseJson
	if err := json.Unmarshal([]byte(data), &res); err != nil {
		return "", false
	}
	if res.Header.TrID == "" || res.Header.TrKey == "" || res.Header.Encrypt == "" {
		return "", false
	}

	return res.Header.TrKey, true
}