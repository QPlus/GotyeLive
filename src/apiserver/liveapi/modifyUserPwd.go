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

func ModifyUserPwd(w http.ResponseWriter, r *http.Request) {
	resp := gotye_protocol.ModifyUserPwdResponse{}
	req := gotye_protocol.ModifyUserPwdRequest{}

	readdata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warn("ModifyUserPwd : ", err.Error())
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	err = json.Unmarshal(readdata, &req)
	if err != nil {
		logger.Warn("ModifyUserPwd : reqdata not json ", string(readdata))
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	logger.Info("ModifyUserPwd : req=", string(readdata))
	service.UserPwdModify(&resp, &req)

end:
	resp.SetAccess("/live/ModifyUserPwd")
	httplib.HttpResponseJson(w, http.StatusOK, &resp)
}
