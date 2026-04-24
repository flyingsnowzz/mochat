package wechat

import (
	"encoding/xml"
	"fmt"

	"github.com/silenceper/wechat/v2/util"
)

type CallbackMsg struct {
	ToUserName string `xml:"ToUserName"`
	AgentID    string `xml:"AgentID"`
	Encrypt    string `xml:"Encrypt"`
}

type CallbackData struct {
	XMLName        xml.Name `xml:"xml"`
	ToUserName     string   `xml:"ToUserName"`
	FromUserName   string   `xml:"FromUserName"`
	CreateTime     int64    `xml:"CreateTime"`
	MsgType        string   `xml:"MsgType"`
	Event          string   `xml:"Event"`
	ChangeType     string   `xml:"ChangeType"`
	Content        string   `xml:"Content"`
	UserID         string   `xml:"UserID"`
	ExternalUserID string   `xml:"ExternalUserID"`
	ChatID         string   `xml:"ChatID"`
	TagID          string   `xml:"TagId"`
	WelcomeCode    string   `xml:"WelcomeCode"`
}

func VerifyURL(token, timestamp, nonce, echoStr, signature string) bool {
	hash := util.Signature(token, timestamp, nonce)
	return hash == signature
}

func DecryptCallbackMsg(encodingAesKey, corpID, encrypt string) ([]byte, error) {
	_, rawMsgXMLBytes, err := util.DecryptMsg(corpID, encrypt, encodingAesKey)
	if err != nil {
		return nil, fmt.Errorf("decrypt failed: %w", err)
	}

	return rawMsgXMLBytes, nil
}
