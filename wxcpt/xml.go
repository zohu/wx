package wxcpt

import (
	"encoding/xml"
)

type CDATA struct {
	Value string `xml:",cdata"`
}

type XmlBizMsg4Send struct {
	XMLName   xml.Name `xml:"xml"`
	Encrypt   CDATA    `xml:"Encrypt"`
	Signature CDATA    `xml:"MsgSignature"`
	Timestamp string   `xml:"TimeStamp"`
	Nonce     CDATA    `xml:"Nonce"`
}

// NewXmlBizMsg4Send
// @Description: 待发送消息xml构造器
// @param encrypt
// @param signature
// @param timestamp
// @param nonce
// @return *BizXmlMsg4Send
func NewXmlBizMsg4Send(encrypt, signature, timestamp, nonce string) *XmlBizMsg4Send {
	return &XmlBizMsg4Send{Encrypt: CDATA{Value: encrypt}, Signature: CDATA{Value: signature}, Timestamp: timestamp, Nonce: CDATA{Value: nonce}}
}
