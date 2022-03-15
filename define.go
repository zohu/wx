package wx

import "encoding/xml"

type ParamNotify struct {
	MsgSignature string `json:"msg_signature"`
	Timestamp    string `json:"timestamp"`
	Nonce        string `json:"nonce"`
	Echostr      string `json:"echostr"`
}
type NotifyEncrypt struct {
	XMLName    xml.Name `xml:"xml"`
	ToUserName string
	Encrypt    string
	AgentID    string
}
