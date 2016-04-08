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

func Login(w http.ResponseWriter, r *http.Request) {
	resp := gotye_protocol.LoginResponse{}
	req := gotye_protocol.LoginRequest{}

	readdata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warn("Login : ", err.Error())
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	err = json.Unmarshal(readdata, &req)
	if err != nil {
		logger.Warn("Login : reqdata not json ", string(readdata))
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}
	logger.Info("Login : req = ", string(readdata))
	service.UserLogin(&resp, &req)
end:
	resp.SetAccess("/live/Login")
	httplib.HttpResponseJson(w, http.StatusOK, &resp)
}
