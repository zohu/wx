package wx

type ParamNotify struct {
	MsgSignature string `json:"msg_signature"`
	Timestamp    string `json:"timestamp"`
	Nonce        string `json:"nonce"`
	Echostr      string `json:"echostr"`
}

type Response struct {
	Errcode int64  `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}
