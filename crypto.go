package wx

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strings"
)

/**
企业微信加解密
*/

const (
	ValidateSignatureError int = -40001
	ParseJsonError         int = -40002
	ComputeSignatureError  int = -40003
	IllegalAesKey          int = -40004
	ValidateCorpidError    int = -40005
	EncryptAESError        int = -40006
	DecryptAESError        int = -40007
	IllegalBuffer          int = -40008
	EncodeBase64Error      int = -40009
	DecodeBase64Error      int = -40010
	GenJsonError           int = -40011
	IllegalProtocolType    int = -40012
)

type ProtocolType int

const (
	JsonType ProtocolType = 1
)

type CryptError struct {
	ErrCode int
	ErrMsg  string
}

type BizJsonMsg4Recv struct {
	Tousername string `json:"tousername"`
	Encrypt    string `json:"encrypt"`
	Agentid    string `json:"agentid"`
}

type BizJsonMsg4Send struct {
	Encrypt   string `json:"encrypt"`
	Signature string `json:"msgsignature"`
	Timestamp string `json:"timestamp"`
	Nonce     string `json:"nonce"`
}
type BizMsgCrypt struct {
	token             string
	encodingAeskey    string
	receiverId        string
	protocolProcessor ProtocolProcessor
}

func NewCryptError(errCode int, errMsg string) *CryptError {
	return &CryptError{ErrCode: errCode, ErrMsg: errMsg}
}
func NewWXBizJsonMsg4Send(encrypt, signature, timestamp, nonce string) *BizJsonMsg4Send {
	return &BizJsonMsg4Send{Encrypt: encrypt, Signature: signature, Timestamp: timestamp, Nonce: nonce}
}

type ProtocolProcessor interface {
	parse(srcData []byte) (*BizJsonMsg4Recv, *CryptError)
	serialize(msgSend *BizJsonMsg4Send) ([]byte, *CryptError)
}

type JsonProcessor struct{}

func (jp *JsonProcessor) parse(srcData []byte) (*BizJsonMsg4Recv, *CryptError) {
	var msg4Recv BizJsonMsg4Recv
	err := json.Unmarshal(srcData, &msg4Recv)
	if nil != err {
		fmt.Println("Unmarshal fail", err)
		return nil, NewCryptError(ParseJsonError, "json to msg fail")
	}
	return &msg4Recv, nil
}

func (jp *JsonProcessor) serialize(msg4Send *BizJsonMsg4Send) ([]byte, *CryptError) {
	jsonMsg, err := json.Marshal(msg4Send)
	if nil != err {
		return nil, NewCryptError(GenJsonError, err.Error())
	}

	return jsonMsg, nil
}
func NewWXBizMsgCrypt(token, encodingAeskey, receiverId string, protocolType ProtocolType) *BizMsgCrypt {
	var protocolProcessor ProtocolProcessor
	if protocolType != JsonType {
		panic("unsupport protocal")
	} else {
		protocolProcessor = new(JsonProcessor)
	}

	return &BizMsgCrypt{token: token, encodingAeskey: encodingAeskey + "=", receiverId: receiverId, protocolProcessor: protocolProcessor}
}

func (jp *BizMsgCrypt) randString(n int) string {
	b := make([]byte, n)
	letterBytes := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func (jp *BizMsgCrypt) pKCS7Padding(plaintext string, blockSize int) []byte {
	padding := blockSize - (len(plaintext) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	var buffer bytes.Buffer
	buffer.WriteString(plaintext)
	buffer.Write(padText)
	return buffer.Bytes()
}

func (jp *BizMsgCrypt) pKCS7Unpadding(plaintext []byte, blockSize int) ([]byte, *CryptError) {
	plaintextLen := len(plaintext)
	if nil == plaintext || plaintextLen == 0 {
		return nil, NewCryptError(DecryptAESError, "pKCS7Unpadding error nil or zero")
	}
	if plaintextLen%blockSize != 0 {
		return nil, NewCryptError(DecryptAESError, "pKCS7Unpadding text not a multiple of the block size")
	}
	paddingLen := int(plaintext[plaintextLen-1])
	return plaintext[:plaintextLen-paddingLen], nil
}

func (jp *BizMsgCrypt) cbcEncrypter(plaintext string) ([]byte, *CryptError) {
	aeskey, err := base64.StdEncoding.DecodeString(jp.encodingAeskey)
	if nil != err {
		return nil, NewCryptError(DecodeBase64Error, err.Error())
	}
	const blockSize = 32
	padMsg := jp.pKCS7Padding(plaintext, blockSize)

	block, err := aes.NewCipher(aeskey)
	if err != nil {
		return nil, NewCryptError(EncryptAESError, err.Error())
	}

	ciphertext := make([]byte, len(padMsg))
	iv := aeskey[:aes.BlockSize]

	mode := cipher.NewCBCEncrypter(block, iv)

	mode.CryptBlocks(ciphertext, padMsg)
	base64Msg := make([]byte, base64.StdEncoding.EncodedLen(len(ciphertext)))
	base64.StdEncoding.Encode(base64Msg, ciphertext)

	return base64Msg, nil
}

func (jp *BizMsgCrypt) cbcDecrypter(base64EncryptMsg string) ([]byte, *CryptError) {
	aeskey, err := base64.StdEncoding.DecodeString(jp.encodingAeskey)
	if nil != err {
		return nil, NewCryptError(DecodeBase64Error, err.Error())
	}

	encryptMsg, err := base64.StdEncoding.DecodeString(base64EncryptMsg)
	if nil != err {
		return nil, NewCryptError(DecodeBase64Error, err.Error())
	}

	block, err := aes.NewCipher(aeskey)
	if err != nil {
		return nil, NewCryptError(DecryptAESError, err.Error())
	}

	if len(encryptMsg) < aes.BlockSize {
		return nil, NewCryptError(DecryptAESError, "encrypt_msg size is not valid")
	}

	iv := aeskey[:aes.BlockSize]

	if len(encryptMsg)%aes.BlockSize != 0 {
		return nil, NewCryptError(DecryptAESError, "encrypt_msg not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(encryptMsg, encryptMsg)

	return encryptMsg, nil
}

func (jp *BizMsgCrypt) calSignature(timestamp, nonce, data string) string {
	sortArr := []string{jp.token, timestamp, nonce, data}
	sort.Strings(sortArr)
	var buffer bytes.Buffer
	for _, value := range sortArr {
		buffer.WriteString(value)
	}

	sha := sha1.New()
	sha.Write(buffer.Bytes())
	signature := fmt.Sprintf("%x", sha.Sum(nil))
	return signature
}

func (jp *BizMsgCrypt) ParsePlainText(plaintext []byte) ([]byte, uint32, []byte, []byte, *CryptError) {
	const blockSize = 32
	plaintext, err := jp.pKCS7Unpadding(plaintext, blockSize)
	if nil != err {
		return nil, 0, nil, nil, err
	}

	textLen := uint32(len(plaintext))
	if textLen < 20 {
		return nil, 0, nil, nil, NewCryptError(IllegalBuffer, "plain is to small 1")
	}
	random := plaintext[:16]
	msgLen := binary.BigEndian.Uint32(plaintext[16:20])
	if textLen < (20 + msgLen) {
		return nil, 0, nil, nil, NewCryptError(IllegalBuffer, "plain is to small 2")
	}

	msg := plaintext[20 : 20+msgLen]
	receiverId := plaintext[20+msgLen:]

	return random, msgLen, msg, receiverId, nil
}

func (jp *BizMsgCrypt) VerifyURL(msgSignature, timestamp, nonce, echostr string) ([]byte, *CryptError) {
	signature := jp.calSignature(timestamp, nonce, echostr)
	fmt.Println(signature, msgSignature)
	if strings.Compare(signature, msgSignature) != 0 {
		return nil, NewCryptError(ValidateSignatureError, "signature not equal")
	}

	plaintext, err := jp.cbcDecrypter(echostr)
	if nil != err {
		return nil, err
	}

	_, _, msg, receiverId, err := jp.ParsePlainText(plaintext)
	if nil != err {
		return nil, err
	}

	if len(jp.receiverId) > 0 && strings.Compare(string(receiverId), jp.receiverId) != 0 {
		fmt.Println(string(receiverId), jp.receiverId, len(receiverId), len(jp.receiverId))
		return nil, NewCryptError(ValidateCorpidError, "receiver_id is not equal")
	}

	return msg, nil
}

func (jp *BizMsgCrypt) EncryptMsg(replyMsg, timestamp, nonce string) ([]byte, *CryptError) {
	randStr := jp.randString(16)
	var buffer bytes.Buffer
	buffer.WriteString(randStr)

	msgLenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(msgLenBuf, uint32(len(replyMsg)))
	buffer.Write(msgLenBuf)
	buffer.WriteString(replyMsg)
	buffer.WriteString(jp.receiverId)

	tmpCiphertext, err := jp.cbcEncrypter(buffer.String())
	if nil != err {
		return nil, err
	}
	ciphertext := string(tmpCiphertext)

	signature := jp.calSignature(timestamp, nonce, ciphertext)

	msg4Send := NewWXBizJsonMsg4Send(ciphertext, signature, timestamp, nonce)
	return jp.protocolProcessor.serialize(msg4Send)
}

func (jp *BizMsgCrypt) DecryptMsg(msgSignature, timestamp, nonce string, msg4Recv *BizJsonMsg4Recv) ([]byte, *CryptError) {
	signature := jp.calSignature(timestamp, nonce, msg4Recv.Encrypt)

	if strings.Compare(signature, msgSignature) != 0 {
		return nil, NewCryptError(ValidateSignatureError, "signature not equal")
	}

	plaintext, cryptErr := jp.cbcDecrypter(msg4Recv.Encrypt)
	if nil != cryptErr {
		return nil, cryptErr
	}

	_, _, msg, receiverId, cryptErr := jp.ParsePlainText(plaintext)
	if nil != cryptErr {
		return nil, cryptErr
	}

	if len(jp.receiverId) > 0 && strings.Compare(string(receiverId), jp.receiverId) != 0 {
		return nil, NewCryptError(ValidateCorpidError, "receiver_id is not equal")
	}

	return msg, nil
}

/**
微信加解密
*/

func MpDecryptMsg(appID, encryptedMsg, aesKey string) (random, rawMsgXMLBytes []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("panic error: err=%v", e)
			return
		}
	}()
	var encryptedMsgBytes, key, getAppIDBytes []byte
	encryptedMsgBytes, err = base64.StdEncoding.DecodeString(encryptedMsg)
	if err != nil {
		return
	}
	key, err = mpAesKeyDecode(aesKey)
	if err != nil {
		panic(err)
	}
	random, rawMsgXMLBytes, getAppIDBytes, err = mpAESDecryptMsg(encryptedMsgBytes, key)
	if err != nil {
		err = fmt.Errorf("消息解密失败,%v", err)
		return
	}
	if appID != string(getAppIDBytes) {
		err = fmt.Errorf("消息解密校验APPID失败")
		return
	}
	return
}

func mpAesKeyDecode(encodedAESKey string) (key []byte, err error) {
	if len(encodedAESKey) != 43 {
		err = fmt.Errorf("the length of encodedAESKey must be equal to 43")
		return
	}
	key, err = base64.StdEncoding.DecodeString(encodedAESKey + "=")
	if err != nil {
		return
	}
	if len(key) != 32 {
		err = fmt.Errorf("encodingAESKey invalid")
		return
	}
	return
}
func mpAESDecryptMsg(ciphertext []byte, aesKey []byte) (random, rawXMLMsg, appID []byte, err error) {
	const (
		BlockSize = 32            // PKCS#7
		BlockMask = BlockSize - 1 // BLOCK_SIZE 为 2^n 时, 可以用 mask 获取针对 BLOCK_SIZE 的余数
	)

	if len(ciphertext) < BlockSize {
		err = fmt.Errorf("the length of ciphertext too short: %d", len(ciphertext))
		return
	}
	if len(ciphertext)&BlockMask != 0 {
		err = fmt.Errorf("ciphertext is not a multiple of the block size, the length is %d", len(ciphertext))
		return
	}

	plaintext := make([]byte, len(ciphertext)) // len(plaintext) >= BLOCK_SIZE

	// 解密
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		panic(err)
	}
	mode := cipher.NewCBCDecrypter(block, aesKey[:16])
	mode.CryptBlocks(plaintext, ciphertext)

	// PKCS#7 去除补位
	amountToPad := int(plaintext[len(plaintext)-1])
	if amountToPad < 1 || amountToPad > BlockSize {
		err = fmt.Errorf("the amount to pad is incorrect: %d", amountToPad)
		return
	}
	plaintext = plaintext[:len(plaintext)-amountToPad]

	// 反拼接
	// len(plaintext) == 16+4+len(rawXMLMsg)+len(appId)
	if len(plaintext) <= 20 {
		err = fmt.Errorf("plaintext too short, the length is %d", len(plaintext))
		return
	}
	rawXMLMsgLen := int(mpDecodeNetworkByteOrder(plaintext[16:20]))
	if rawXMLMsgLen < 0 {
		err = fmt.Errorf("incorrect msg length: %d", rawXMLMsgLen)
		return
	}
	appIDOffset := 20 + rawXMLMsgLen
	if len(plaintext) <= appIDOffset {
		err = fmt.Errorf("msg length too large: %d", rawXMLMsgLen)
		return
	}

	random = plaintext[:16:20]
	rawXMLMsg = plaintext[20:appIDOffset:appIDOffset]
	appID = plaintext[appIDOffset:]
	return
}

func mpDecodeNetworkByteOrder(orderBytes []byte) (n uint32) {
	return uint32(orderBytes[0])<<24 |
		uint32(orderBytes[1])<<16 |
		uint32(orderBytes[2])<<8 |
		uint32(orderBytes[3])
}
