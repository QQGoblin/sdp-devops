package wxwork

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"time"
)

type WXWorkClient struct {
	token
	CorpId     string
	CorpSecret string
}

type GetTokenResp struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"err_msg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func New(corpId, corpsecret string) (t *WXWorkClient) {

	wxWorkClient := WXWorkClient{
		token:      token{},
		CorpId:     corpId,
		CorpSecret: corpsecret,
	}

	wxWorkClient.token.setGetTokenFunc(wxWorkClient.getTokenFunc)
	wxWorkClient.token.tokenRefresher(context.Background())

	return &wxWorkClient

}

func (t *WXWorkClient) getToken() (tokenInfo, error) {

	tinfo := tokenInfo{}

	client := resty.New()
	var rbody GetTokenResp
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetPathParam("CORPID", t.CorpId).
		SetPathParam("CORPSECRET", t.CorpSecret).
		SetResult(&rbody).
		Get("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid={CORPID}&corpsecret={CORPSECRET}")

	if err != nil || resp.StatusCode() != 200 {
		return tinfo, err
	}
	tinfo.token = rbody.AccessToken
	tinfo.expiresIn = time.Duration(rbody.ExpiresIn) * time.Second

	return tinfo, nil

}

type TextCardMsg struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
}

type AppMsg struct {
	ToParty string      `json:"toparty"`
	MsgType string      `json:"msgtype"`
	AgentId int         `json:"agentid"`
	Content TextCardMsg `json:"textcard"`
}

func (t *WXWorkClient) SendMsg(agentid int, toparty string, content TextCardMsg) {

	reqBody := AppMsg{
		ToParty: toparty,
		MsgType: "textcard",
		AgentId: agentid,
		Content: content,
	}

	client := resty.New()
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		SetPathParam("token", t.token.getToken()).
		Post("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token={token}")
	if err != nil {
		logrus.Error("发送微信告警失败")
	}
}
