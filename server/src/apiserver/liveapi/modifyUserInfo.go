package liveapi

import (
	"apiserver/service"
	"encoding/json"
	"gotye_protocol"
	"io/ioutil"
	"net/http"

	"github.com/futurez/litego/httplib"
	"github.com/futurez/litego/logger"
)

func ModifyUserInfo(w http.ResponseWriter, r *http.Request) {
	resp := gotye_protocol.ModifyUserInfoResponse{}
	req := gotye_protocol.ModifyUserInfoRequest{}

	readdata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warn("ModifyUserInfo : ", err.Error())
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	err = json.Unmarshal(readdata, &req)
	if err != nil {
		logger.Warn("ModifyUserInfo : reqdata not json ", string(readdata))
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	if len(req.SessionID) == 0 {
		logger.Warn("ModifyUserInfo : sessionID is nul ", string(readdata))
		resp.SetStatus(gotye_protocol.API_EXPIRED_SESSION_ERROR)
		goto end
	}

	if len(req.NickName) == 0 && len(req.Address) == 0 &&
		(req.Sex != 1 || req.Sex != 2) {
		logger.Warn("ModifyUserInfo : param error. ", string(readdata))
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	logger.Info("ModifyUserInfo : req=", string(readdata))
	service.UserInfoModify(&resp, &req)

end:
	resp.SetAccess("/live/ModifyUserInfo")
	httplib.HttpResponseJson(w, http.StatusOK, &resp)
}
