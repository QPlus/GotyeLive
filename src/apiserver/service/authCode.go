package service

import (
	"fmt"
	"gotye_protocol"
	"net/url"
	"sync"

	"github.com/futurez/litego/httplib"
	"github.com/futurez/litego/logger"
	"github.com/futurez/litego/util"
)

func SendAuthCode(phone string) (string, error) {
	v := url.Values{}
	v.Set("uid", SP_appInfo.smsUid)
	v.Set("auth", SP_appInfo.smsAuth)
	v.Set("mobile", phone)
	v.Set("expid", "0")
	v.Set("encode", "utf-8")

	authCode := util.AuthCode()

	v.Set("msg", fmt.Sprintf("验证码%s,请您进行校验,请勿泄漏.", authCode))
	req := v.Encode()

	resp, err := httplib.HttpRequest("http://sms.10690221.com:9011/hy/", httplib.METHOD_GET, nil, []byte(req))
	if err != nil {
		logger.Error("SendAuthCode : ", err.Error())
		return "", err
	}
	logger.Infof("SendAuthCode : authCode=%s, resp=%s", authCode, string(resp))
	return authCode, nil
}

type PhoneAuthcode struct {
	sync.RWMutex
	phoneMap map[string]string
}

func (p *PhoneAuthcode) Set(phone, code string) {
	p.Lock()
	defer p.Unlock()

	p.phoneMap[phone] = code
}

func (p *PhoneAuthcode) Delete(phone string) {
	p.RLock()
	defer p.RUnlock()

	delete(p.phoneMap, phone)
}

func (p *PhoneAuthcode) Check(phone, code string) bool {
	p.RLock()
	defer p.RUnlock()

	v, ok := p.phoneMap[phone]
	if !ok {
		return false
	}
	return (v == code)
}

var SP_phoneCode = &PhoneAuthcode{phoneMap: make(map[string]string, 0)}

func RequestAuthCode(resp *gotye_protocol.AuthCodeResponse, req *gotye_protocol.AuthCodeRequest) {
	if DBIsPhoneExists(req.Phone) {
		resp.SetStatus(gotye_protocol.API_PHONE_EXISTS_ERROR)
		logger.Info("RequestAuthCode : phone = ", req.Phone, " is exists.")
		return
	}

	code, err := SendAuthCode(req.Phone)
	if err != nil {
		resp.SetStatus(gotye_protocol.API_SERVER_ERROR)
		logger.Info("RequestAuthCode : SendAuthCode failed.")
		return
	}

	SP_phoneCode.Set(req.Phone, code)
	logger.Infof("RequestAuthCode: phone=%s, code=%s", req.Phone, code)
	resp.SetStatus(gotye_protocol.API_SUCCESS)
}
