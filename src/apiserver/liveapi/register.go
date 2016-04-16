package liveapi

import (
	"apiserver/service"
	"encoding/json"
	"gotye_protocol"
	"io/ioutil"
	"net/http"

	"github.com/futurez/litego/httplib"
	"github.com/futurez/litego/logger"
	"github.com/futurez/litego/util"
)

func Register(w http.ResponseWriter, r *http.Request) {
	resp := gotye_protocol.RegisterResponse{}
	req := gotye_protocol.RegisterRequest{}

	readdata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warn("Register : ", err.Error())
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	err = json.Unmarshal(readdata, &req)
	if err != nil {
		logger.Warn("Register : reqdata not json ", string(readdata))
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	if !util.CheckPhone(req.Phone) {
		logger.Warn("Register : phone is null. ", string(readdata))
		resp.SetFormatStatus(gotye_protocol.API_PARAM_ERROR, "手机号码不存在")
		goto end
	}

	if len(req.Passwd) < 6 {
		logger.Warn("Register : pwd length less 6, =", req.Passwd)
		resp.SetFormatStatus(gotye_protocol.API_PARAM_ERROR, "密码小于6位")
		goto end
	}

	logger.Info("Register : phone=", req.Phone)
	service.UserRegister(&resp, &req)

end:
	resp.SetAccess("/live/Register")
	httplib.HttpResponseJson(w, http.StatusOK, &resp)
}
