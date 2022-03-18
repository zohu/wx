package wxcpt

type JsonBizMsg4Send struct {
	Encrypt   string `json:"encrypt"`
	Signature string `json:"msgsignature"`
	Timestamp string `json:"timestamp"`
	Nonce     string `json:"nonce"`
}

// NewJsonBizMsg4Send
// @Description: 待发送消息json构造器
// @param encrypt
// @param signature
// @param timestamp
// @param nonce
// @return *JsonBizMsg4Send
func NewJsonBizMsg4Send(encrypt, signature, timestamp, nonce string) *JsonBizMsg4Send {
	return &JsonBizMsg4Send{Encrypt: encrypt, Signature: signature, Timestamp: timestamp, Nonce: nonce}
}
