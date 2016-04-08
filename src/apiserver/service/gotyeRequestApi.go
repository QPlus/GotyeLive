package service

import (
	"gotye_protocol/gotye_sdk"
	"strconv"
	"sync"
	"time"

	"github.com/futurez/litego/httplib"
	"github.com/futurez/litego/logger"
	"github.com/futurez/litego/util"
)

type TokenCache struct {
	sync.RWMutex
	accessToken string
	expires     time.Time
	bValid      bool
}

func (t *TokenCache) GetAccessToken() (string, bool) {
	t.RLock()
	defer t.RUnlock()

	if !t.bValid {
		return "", false
	}

	if t.expires.After(time.Now()) {
		logger.Infof("AccessAppToken : quick accessToken=%s", t.accessToken)
		return t.accessToken, true
	}
	return "", false
}

func (t *TokenCache) SetAccessToken(token string, expire int64) {
	t.Lock()
	defer t.Unlock()

	t.accessToken = token
	//before 60s accessAppToken
	t.expires = time.Now().Add(time.Second * time.Duration(expire-60))
	t.bValid = true
}

func (t *TokenCache) ResetAccessToken() {
	t.Lock()
	defer t.Unlock()

	t.bValid = false
	t.accessToken = ""
}

var SP_tokenCache = &TokenCache{}

func GotyeClearAccessToken() {
	SP_tokenCache.ResetAccessToken()
}

func GotyeAccessAppToken() (string, error) {
	if token, ok := SP_tokenCache.GetAccessToken(); ok {
		return token, nil
	}

	//send request.
	req := gotye_sdk.AccessTokenAppRequest{}
	req.Scope = "app"
	req.UserName = SP_appInfo.username
	req.Password = SP_appInfo.password
	resp := gotye_sdk.AccessTokenResponse{}

	logger.Info("GotyeAccessAppToken : req=", req)
	err := httplib.HttpRequestJson(gotye_sdk.HttpUrlAccessToken, &req, &resp)
	if err != nil {
		logger.Error("GotyeAccessAppToken : ", err.Error())
		return "", err
	}
	logger.Info("GotyeAccessAppToken : resp=", resp)
	SP_tokenCache.SetAccessToken(resp.AccessToken, resp.ExpiresIn)
	return resp.AccessToken, nil
}

func GotyeAccessRoomToken(liveRoomId int64, password, nickname string) (string, error) {
	//send request.
	req := gotye_sdk.AccessTokenRoomRequest{}
	req.Scope = "room"
	req.RoomId = liveRoomId
	req.Password = password
	req.NickName = nickname
	req.SecretKey = util.Md5Hash(strconv.FormatInt(liveRoomId, 10) + password + SP_appInfo.accessSecret)

	resp := gotye_sdk.AccessTokenResponse{}

	logger.Info("GotyeAccessRoomToken : ", req)
	err := httplib.HttpRequestJson(gotye_sdk.HttpUrlAccessToken, &req, &resp)
	if err != nil {
		logger.Error("GotyeAccessRoomToken : ", err.Error())
		return "", err
	}
	logger.Info("GotyeAccessRoomToken : ", resp)
	return resp.AccessToken, nil
}

func GotyeCreateRoom(appToken, roomName, anchorPwd, assistPwd, userPwd, anchorDesc,
	contentDesc string) (*gotye_sdk.CreateRoomResponse, error) {

	headers := map[string]string{
		"Authorization": appToken,
	}

	resp := &gotye_sdk.CreateRoomResponse{}
	req := gotye_sdk.CreateRoomRequest{}
	req.RoomName = roomName
	req.AnchorPwd = anchorPwd
	req.AssistPwd = assistPwd
	req.UserPwd = userPwd
	req.AnchorDesc = anchorDesc
	req.ContentDesc = contentDesc

	logger.Info("GotyeCreateRoom : req=", req)
	err := httplib.HttpRequestJsonToken(gotye_sdk.HttpUrlCreateRoom, headers, &req, resp)
	if err != nil {
		logger.Error("GotyeCreateRoom : ", err.Error())
		return nil, err
	}
	logger.Info("GotyeCreateRoom : resp=", resp)
	return resp, nil
}

func GotyeModifyRoom(appToken string, roomId int64, roomName, anchorPwd, assistPwd,
	userPwd, anchorDesc, contentDesc string) int {

	headers := map[string]string{
		"Authorization": appToken,
	}

	resp := gotye_sdk.ModifyRoomResponse{}
	//req := gotye_sdk.ModifyRoomRequest{}
	req := make(map[string]interface{})
	req["roomId"] = roomId
	req["enableRecordFlag"] = 1
	req["permanentPlayFlag"] = 1
	req["startPlayTime"] = time.Now().Second()
	if len(roomName) > 0 {
		req["roomName"] = roomName
	}
	if len(anchorPwd) > 0 {
		req["anchorPwd"] = anchorPwd
	}
	if len(assistPwd) > 0 {
		req["assistPwd"] = assistPwd
	}
	if len(userPwd) > 0 {
		req["userPwd"] = userPwd
	}
	if len(anchorDesc) > 0 {
		req["anchorDesc"] = anchorDesc
	}
	if len(contentDesc) > 0 {
		req["contentDesc"] = contentDesc
	}

	logger.Info("GotyeModifyRoom : req=", req)
	err := httplib.HttpRequestJsonToken(gotye_sdk.HttpUrlModifyRoom, headers, req, &resp)
	if err != nil {
		logger.Error("GotyeCreateRoom : ", err.Error())
		return -1
	}
	logger.Info("GotyeModifyRoom : resp=", resp)
	return resp.Status
}

func GotyeGetLiveContext(appToken string, roomId int64) (int, int, error) {
	headers := map[string]string{
		"Authorization": appToken,
	}
	req := gotye_sdk.GetLiveContextRequest{RoomId: roomId}
	resp := gotye_sdk.GetLiveContextResponse{}

	logger.Info("GotyeGetLiveContext : req=", req)
	err := httplib.HttpRequestJsonToken(gotye_sdk.HttpGetLiveContext, headers, req, &resp)
	if err != nil {
		logger.Error("GotyeGetLiveContext : ", err.Error())
		return 0, 0, err
	}
	logger.Info("GotyeGetLiveContext : resp=", resp)
	return resp.Entity.PlayUserCount, resp.Status, nil
}
