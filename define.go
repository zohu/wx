package wx

import "encoding/xml"

type Response struct {
	Errcode int64  `json:"errcode,omitempty"`
	Errmsg  string `json:"errmsg,omitempty"`
}
type ParamNotify struct {
	MsgSignature string `json:"msg_signature"`
	Timestamp    string `json:"timestamp"`
	Nonce        string `json:"nonce"`
	Echostr      string `json:"echostr,omitempty"`
	EncryptType  string `json:"encrypt_type,omitempty"`
}

func (n *ParamNotify) IsSafeMode() bool {
	return n.EncryptType == "aes"
}

// CDATA  使用该类型,在序列化为 xml 文本时文本会被解析器忽略
type CDATA string

func (c CDATA) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(struct {
		string `xml:",cdata"`
	}{string(c)}, start)
}

type Err struct {
	Errcode int64  `json:"errcode,omitempty"`
	Errmsg  string `json:"errmsg,omitempty"`
	Appid   string `json:"appid"`
	Err     string `json:"err"`
	Desc    string `json:"desc"`
}
