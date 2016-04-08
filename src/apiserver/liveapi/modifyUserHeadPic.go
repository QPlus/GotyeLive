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

func ModifyUserHeadPic(w http.ResponseWriter, r *http.Request) {
	resp := gotye_protocol.ModifyUserHeadPicResponse{}
	req := gotye_protocol.ModifyUserHeadPicRequest{}

	readdata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warn("ModifyUserHeadPic : ", err.Error())
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	err = json.Unmarshal(readdata, &req)
	if err != nil {
		logger.Warn("ModifyUserHeadPic : reqdata not json ", string(readdata))
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end

	}

	if len(req.SessionID) == 0 {
		logger.Warn("ModifyUserHeadPic : sessionID is nul ")
		resp.SetStatus(gotye_protocol.API_EXPIRED_SESSION_ERROR)
		goto end
	}

	logger.Infof("ModifyUserHeadPic : req sessionId=%s, piclen=%d", req.SessionID, len(req.HeadPic))
	service.UserHeadPicModify(&resp, &req)

end:
	resp.SetAccess("/live/ModifyUserHeadPic")
	httplib.HttpResponseJson(w, http.StatusOK, &resp)
}
